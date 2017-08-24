package main

import "fmt"

type PlayScenes struct {

}

func (ps PlayScenes) New() interface{} {
	p := new(PlayScenes)
	return p
}

func (ps PlayScenes) Run(stream *Stream) {
	fmt.Println("PlayScenes")
	if src := GetStream(stream); src != nil {
		fmt.Println(src)
	}
}
