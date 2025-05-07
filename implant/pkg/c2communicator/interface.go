package c2communicator

import (
	"asritha.dev/implant/pkg/task"
)

type C2Communicator interface {
	Beacon(implantID uint, results []task.Result) (*task.Queue, error)
}

// NewC2Communicator picks the right implementation based on protocol.
//func New(protocol string) (C2Communicator, error) {
//	switch protocol {
//	case "HTTPS":
//		return NewHTTPSC2Communicator(), nil
//	case "DNS":
//		return NewDNSC2Communicator(), nil
//	default:
//		return nil, fmt.Errorf("unsupported C2 protocol %q", protocol)
//	}
//}
