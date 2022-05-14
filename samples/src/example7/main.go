package main

// #noverifyiface
type Ok struct {
}

func (o Ok) Do() error {
	return nil
}

func main() {
}

type Iface interface {
	Do() error
}
