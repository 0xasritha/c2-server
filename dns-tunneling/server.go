package main

import "github.com/miekg/dns"

func Server() {
	server := &dns.Server{
		Addr: ":53",
		Net:  "udp",
	}
	go server.ListenAndServe()
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg) // response message

		/*
			type Msg struct {
				    MsgHdr
				    Compress bool `json:"-"`
				    Question []Question
				    Answer   []RR
				    Ns       []RR
				    Extra    []RR
			}

			DNS Message Sections:
			- 1) Header: includes fields specifying which of the remaining sections are present, specifies whether message is query/response, standard query, etc.
			- 2) Question: contains fields (QTYPE, QCLASS, QNAME) that describe question to name server
			- 3) Answer: RRs answering the question
			- 4) Authority: RRs pointing towards an authoritative name server
			- 5) Additional: RRs w/ additional info related to the query

		*/

		randomIP := "23.53.44.43" // TODO: make random IP?

		m.SetReply(r) // SetReply creates a reply message from the response

		w.WriteMsg(m)

		// use Tsig for security?
	})
}
