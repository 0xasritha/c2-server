package main

/*
KrakenListener will:
- serve tasks to an implant
- store records of tasks sent + their results
*/
type Listener struct {
	TasksHandler
	ResultManager
	HistoryManager
}
