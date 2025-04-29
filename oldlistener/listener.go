package main

import (
	"asritha.dev/listener/pkg/results"
	"asritha.dev/listener/pkg/taskrecords"
	"asritha.dev/listener/pkg/tasks"
)

type Listener struct {
	tasksHandler       *tasks.Handler // SHOULD BE pointers?
	resultsHandler     *results.Handler
	taskRecordsHandler *taskrecords.Handler
}
