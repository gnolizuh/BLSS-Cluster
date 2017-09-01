package main

import (
	"fmt"
	"net/http"
)

type PlayScenes struct {

}

func (ps PlayScenes) New() interface{} {
	p := new(PlayScenes)
	return p
}

func (ps PlayScenes) Name() string {
	return "PlayScenes"
}

func (ps PlayScenes) Run(stream *Stream, w http.ResponseWriter, r *http.Request) {
	if dst_stream := GetStream(stream); dst_stream != nil {
		if !(stream.localAddr == dst_stream.localAddr &&
			stream.localPort == dst_stream.localPort) {

			dst_url := fmt.Sprintf("rtmp://%s:%d/%s/%s/%s",
				dst_stream.localAddr, dst_stream.localPort, stream.host,
				stream.app, stream.name)

			// set dst url
			w.Header().Set("Location", dst_url)

			// do w.WriteHeader() after set header
			w.WriteHeader(http.StatusFound)
		}
	}
}
