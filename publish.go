package main

import "fmt"

type PublishScenes struct {

}

func (ps PublishScenes) New() interface{} {
	p := new(PublishScenes)
	return p
}

func (ps PublishScenes) Run(stream *Stream) {
	fmt.Println("PublishScenes")
	AddStream(stream)
}