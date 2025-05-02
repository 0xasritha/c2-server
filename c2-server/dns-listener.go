package main

import "github.com/miekg/dns"

type DNSListener struct {
}

func (listener *DNSListener) listen() {
	server := &dns.Server{
		Addr: ":53",
		Net:  "udp",
	}
	go server.ListenAndServe()
	// TODO: figure out how to ship tasks as resource Records

}
