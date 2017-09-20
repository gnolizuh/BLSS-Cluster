package main

import (
	"net/http"
)

type StreamMapType map[string]*Stream
type AppMapType map[string]StreamMapType
type ServiceMapType map[string]AppMapType

type Scenes interface {
	New(m *Manager) Scenes
	Name() string
	Run(stream *Stream)
	Done(w http.ResponseWriter)
}

type BaseScenes struct {
	manager *Manager
	status   int
	headers  map[string]string
}

func (scenes *BaseScenes) Done(w http.ResponseWriter) {
	for k, v := range scenes.headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(scenes.status)
}