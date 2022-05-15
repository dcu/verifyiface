package a

type Iface interface {
	Do()
}

func DoAny(a interface{}) {
	iface, assertionWorked := a.(Iface)

	if assertionWorked {
		iface.Do()
	}
}
