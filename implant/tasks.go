package main

import (
	"log"
	"os/exec"
)

type Result struct {
	Content string `json:"content"`
	Success bool   `json:"success"`
	TaskID  uint   `json:"task_id"`
}

type Tasker interface {
	Execute() Result
	GetID() uint
}

type TaskQueue struct {
	tasks []Tasker
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{tasks: make([]Tasker, 0)}
}

func (q *TaskQueue) Enqueue(t Tasker) {
	q.tasks = append(q.tasks, t)
}

// Dequeue returns the front Tasker and removes it from the queue
func (q *TaskQueue) Dequeue() Tasker {
	if len(q.tasks) == 0 {
		return nil
	}
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task
}
func (q *TaskQueue) IsEmpty() bool {
	return len(q.tasks) == 0
}

type TaskMetadata struct {
	ID uint `json:"id"`
}

type Ping struct {
	TaskMetadata
}

func NewPingTask(ID uint) *Ping {
	return &Ping{
		TaskMetadata{
			ID: ID,
		},
	}
}

func (p *Ping) Execute() Result {
	return Result{Content: "pong", Success: true}
}

func (p *Ping) GetID() uint {
	return p.ID
}

type ShellCommand struct {
	TaskMetadata
	commandWithArgs []string
}

func NewShellCommandTask(ID uint, commandWithArgs []string) *ShellCommand {
	return &ShellCommand{
		TaskMetadata: TaskMetadata{
			ID: ID,
		},
		commandWithArgs: commandWithArgs,
	}
}
func (s *ShellCommand) GetID() uint {
	return s.ID
}

func (s *ShellCommand) Execute() Result {

	// TODO: fix this logic, have it send cmd + args separately
	cmd := exec.Command(s.commandWithArgs[0], s.commandWithArgs[1:]...)

	output, err := cmd.CombinedOutput() //  CombinedOutput runs the command and returns its combined standard output and standard error.

	var result Result

	// better way to do this?
	if err != nil {
		log.Println("Error executing command:", err)
		result.Success = false
		result.Content = err.Error()
	} else {
		result.Success = true
		result.Content = string(output)
	}

	return result
}
