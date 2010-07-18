package main

type KeyMan interface {
	Atomic(func())
}


type keyman struct {
	in, out chan bool
}

func (k *keyman) manageGate() {
	for {
		k.out <- <-k.in
	}
}

func (k *keyman) Atomic(foo func()) {
	k.in <- true
	foo()
	<-k.out
}

func TestKeyMan() KeyMan {
	k := new(keyman)
	k.in = make(chan bool)
	k.out = make(chan bool)

	go k.manageGate()

	return k
}
