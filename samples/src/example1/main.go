package main

type Ok struct { // want `struct Ok doesn't verify interface compliance for Iface`
}

func (o Ok) Do() error {
	return nil
}

// var _ Iface = (*Ok)(nil)

// var _ Iface = &Ok{}

// var _ Iface = Ok{}

func main() {
}

type Iface interface {
	Do() error
}
