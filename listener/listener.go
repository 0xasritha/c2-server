package main

type Listener interface {
	listen()
}

type HTTPSListener struct {
}

type HTTPListener struct {
}

func (listener *HTTPListener) listen() {}

type DNSListener struct {
}

func (listener *DNSListener) listen() {}
