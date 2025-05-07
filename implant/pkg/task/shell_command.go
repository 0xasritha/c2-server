package task

import (
	"log"
	"os/exec"
)

type ShellCommand struct {
	TaskMetadata
	commandWithArgs []string
}

func NewShellCommand(ID uint, commandWithArgs []string) *ShellCommand {
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
