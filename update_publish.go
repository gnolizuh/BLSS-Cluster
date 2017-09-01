package main

import "net/http"

type UpdatePublishScenes struct {

}

func (ps UpdatePublishScenes) New() interface{} {
	p := new(UpdatePublishScenes)
	return p
}

func (ps UpdatePublishScenes) Name() string {
	return "UpdatePublishScenes"
}

func (ps UpdatePublishScenes) Run(stream *Stream, w http.ResponseWriter, r *http.Request) {
	AddStream(stream)
}