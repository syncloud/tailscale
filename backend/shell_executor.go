package main

import (
	"fmt"
	"os/exec"
)

type ShellExecutor struct {
}

func (e *ShellExecutor) Run(name string, args ...string) ([]byte, error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("%s %v: %w: %s", name, args, err, output)
	}
	return output, nil
}
