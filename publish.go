package main

import "net/http"

type PublishScenes struct {

}

func (ps PublishScenes) New() interface{} {
	p := new(PublishScenes)
	return p
}

func (ps PublishScenes) Name() string {
	return "PublishScenes"
}

func (ps PublishScenes) Run(stream *Stream, w http.ResponseWriter, r *http.Request) {
	AddStream(stream)
}