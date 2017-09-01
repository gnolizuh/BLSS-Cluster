package main

import (
	"strconv"
	"sync"
	"net/http"
)

var StreamMutex = &sync.Mutex{}

type Stream struct {
	localAddr string
	localPort uint64
	service   string
	app       string
	host      string
	name      string
	streams   ServiceMapType
}

func NewStream(r *http.Request, f map[string][]string) *Stream {

	s := new(Stream)
	s.localAddr = r.Header.Get("Local-Addr")
	if i, err := strconv.ParseUint(r.Header.Get("Local-Port"), 10, 64); err == nil {
		s.localPort = i
	}
	s.service = f["service"][0]
	s.host = f["host"][0]
	s.app = f["app"][0]
	s.name = f["name"][0]

	return s
}

func AddStream(stream *Stream) {
	defer StreamMutex.Unlock()
	StreamMutex.Lock()
	stream_list := stream.streams

	apps := stream_list[stream.service]
	if apps == nil {
		apps = make(AppMapType)
		stream_list[stream.service] = apps
	}

	streams := apps[stream.app]
	if streams == nil {
		streams = make(StreamMapType)
		apps[stream.app] = streams
	}

	_, ok := streams[stream.name]
	if !ok {
		streams[stream.name] = stream
	}
}

func DelStream(stream *Stream) {
	defer StreamMutex.Unlock()
	StreamMutex.Lock()
	stream_list := stream.streams

	apps := stream_list[stream.service]
	if apps == nil {
		return
	}

	streams := apps[stream.app]
	if streams == nil {
		return
	}

	_, ok := streams[stream.name]
	if !ok {
		return
	}

	delete(streams, stream.name)
	if len(streams) == 0 {
		delete(apps, stream.app)
		if len(apps) == 0 {
			delete(stream_list, stream.service)
		}
	}
}

func GetStream(stream *Stream) *Stream {
	defer StreamMutex.Unlock()
	StreamMutex.Lock()
	stream_list := stream.streams

	apps := stream_list[stream.service]
	if apps == nil {
		return nil
	}

	streams := apps[stream.app]
	if streams == nil {
		return nil
	}

	s := streams[stream.name]
	if s == nil {
		return nil
	}

	// copy memory of s
	cs := *s

	return &cs
}