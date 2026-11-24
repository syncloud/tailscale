package main

type Executor interface {
	Run(name string, args ...string) ([]byte, error)
}
