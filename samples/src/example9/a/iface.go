package a

type Iface interface {
	Do()
}

func DoAny(a interface{}) {
	iface, ok := a.(Iface)

	if ok {
		iface.Do()
	}
}
