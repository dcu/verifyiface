package main

type Iface1 interface {
	Do1() error
}

type Iface2 interface {
	Do2() error
}

type Iface3 interface {
	Do3() error
}

type Ok struct {
}

func (o Ok) Do1() error {
	return nil
}

func (o Ok) Do2() error {
	return nil
}

func (o Ok) Do3() error {
	return nil
}

var _ Iface1 = (*Ok)(nil)
var _ Iface2 = &Ok{}
var _ Iface3 = Ok{}

func main() {
}
