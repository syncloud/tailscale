package main

import (
	"fmt"
	"strings"
)

type FakeExecutor struct {
	outputs map[string]string
	failing map[string]bool
	calls   []string
}

func NewFakeExecutor() *FakeExecutor {
	return &FakeExecutor{outputs: map[string]string{}, failing: map[string]bool{}}
}

func (e *FakeExecutor) Run(name string, args ...string) ([]byte, error) {
	call := name + " " + strings.Join(args, " ")
	e.calls = append(e.calls, call)
	for prefix, output := range e.outputs {
		if strings.HasPrefix(call, prefix) {
			if e.failing[prefix] {
				return []byte(output), fmt.Errorf("failed: %s", call)
			}
			return []byte(output), nil
		}
	}
	return nil, fmt.Errorf("unexpected call: %s", call)
}
