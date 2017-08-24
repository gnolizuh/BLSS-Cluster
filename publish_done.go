package main

import "fmt"

type PublishDoneScenes struct {

}

func (ps PublishDoneScenes) New() interface{} {
	p := new(PublishDoneScenes)
	return p
}

func (ps PublishDoneScenes) Run(stream *Stream) {
	fmt.Println("PublishDoneScenes")
	DelStream(stream)
}