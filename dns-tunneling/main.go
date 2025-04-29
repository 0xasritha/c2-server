package main

import "github.com/miekg/dns"

func main() {

	// wire format: exact binary representation of DNS message (how its serialized into raw bytes)

	// to create new resource record:
	r := new(dns.MX)
	r.Hdr = dns.RR_Header{
		Name:   "miek.n1",
		Rrtype: dns.TypeMX,
		Class:  dns.ClassINET,
		Ttl:    3600,
	}
	r.Preference = 10
	r.Mx = "mx.miek.n1"

	// DNS message is QUERY (Header + Question) or RESPONSE (Header + Question + Answer [+ Authority + Additional])

	m := new(dns.Msg)
	m.SetQuestion("miek.n1", dns.TypeMX) // question section set to ask Mx records for the miek.n1. zone

	// after creating message, query DNS server at 127.0.0.1 + port 53 with it
	c := new(dns.Client)
	in, rtt, err := c.Exchange(m, "8.8.8.8:53")
	print(in, rtt, err)

	/*
		`in` is the DNS message sent back
		4 sections:
		- 1) in.Question: question section
		- 2) in.Answer: answer section
		- 3) in.Ns: authority section
		- 4) in.Extra: additional section

		each section contains a []RR
	*/

	// to access rdata of a TXT RR as first RR in Answer section
	if t, ok := in.Answer[0].(*dns.TXT); ok {
		print(t)
	}

	// TSIG (Transaction SIGnature): security mechanism to authenticate DNS messages between client + server. Authenticates that the DNS message came from who it claims, ensures message wasn't tampered w/ during transmission

	/*

		Resource Record: single piece of info in DNS database (maps domain name to some kind of data (ex. IP address, another domain, etc.)
		- ex. `NS`: points to nameserver for that domain
		- structure: <Name> <TTL> <Class> <Type> <Data>
			- ex. `example.com. 3600 IN A 93.53.23.41`
				- IN -> Class (almost always "IN" for Internet)

	*/
}
