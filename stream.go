package main

import (
	"strconv"
	"sync"
	"net/http"
	"fmt"
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

func (s *Stream) getUniqueKey() string {
	return fmt.Sprintf("%s:%s:%s", s.service, s.app, s.name)
}

func (s *Stream) getUniqueVal() string {
	return fmt.Sprintf("%s:%d", s.localAddr, s.localPort)
}
