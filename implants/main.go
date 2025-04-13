package main

func main() {
	// set listener API endpoint #TODO: better way to do this?
}

type Implant struct {
	// thread for servicing tasks
	beacon       func()
	setMeanDwell func(meanDwell float64)
	setRunning   func(isRunning bool)
	serviceTasks func()

	host, port, uri string // listening post endpoint args

}

func NewImplant() Implant {
	return Implant{}
}
