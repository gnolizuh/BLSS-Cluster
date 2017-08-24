package main

import "fmt"

type Foo struct {
	i int
}

func main() {
	p := new(Foo)
	p.i = 1

	p1 := *p

	fmt.Printf("p:%p, &(*p):%p, &p1:%p", p, &(*p), &p1)
}