package main

type Ok struct {
}

func (o Ok) Do() error {
	return nil
}

// var _ Iface = (*Ok)(nil)

// var _ Iface = &Ok{}

// var _ Iface = Ok{}

func main() {
	// http.Get("")
}

type Iface interface {
	Do() error
}
