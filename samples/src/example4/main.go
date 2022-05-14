package main

type Iface interface {
	Do() error
}

type Ok struct {
}

func (o Ok) Do() error {
	return nil
}

var _ Iface = Ok{}

func main() {
}
