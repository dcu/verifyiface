package main

import (
	"example9/a"
	"fmt"
)

type Ok struct { // want `struct Ok doesn't verify interface compliance for a.Iface`
}

func (o *Ok) Do() {
	fmt.Println("ok.Do()")
}

func main() {
	o := &Ok{}

	a.DoAny(o)
}
