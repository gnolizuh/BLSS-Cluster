package main

import "net/http"

type PublishDoneScenes struct {

}

func (ps PublishDoneScenes) New() interface{} {
	p := new(PublishDoneScenes)
	return p
}

func (ps PublishDoneScenes) Name() string {
	return "PublishDoneScenes"
}

func (ps PublishDoneScenes) Run(stream *Stream, w http.ResponseWriter, r *http.Request) {
	DelStream(stream)
}
