package main

type DNSC2Communicator struct{}

func (D *DNSC2Communicator) Beacon(implantID uint, results []Result) (*TaskQueue, error) {

	return NewTaskQueue(), nil
}

func NewDNSC2Communicator() *DNSC2Communicator {
	return &DNSC2Communicator{}
}
