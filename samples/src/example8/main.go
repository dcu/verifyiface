package main

var _ Iface = (*Ok)(nil)

type Ok struct {
}

type Iface interface {
	Do() error
}

func (o *Ok) Do() error {
	return nil
}

func main() {
}
