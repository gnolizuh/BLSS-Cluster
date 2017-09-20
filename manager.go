package main

import (
	"net/http"
	"io/ioutil"
	"strings"
	"os"
	"log"
	"github.com/go-redis/redis"
)

type MethodType func(*http.Request) map[string][]string

type Manager struct {
	redisClient *redis.Client
	config      *Config
	methods      map[string]MethodType
	scenes       map[string]Scenes
	streams      ServiceMapType
	baseId       uint64
}

func NewManager(config *Config) (*Manager) {
	m := new(Manager)
	m.config = config
	m.methods = make(map[string]MethodType)
	m.methods["GET"] = func(r *http.Request) map[string][]string {
		return r.URL.Query()
	}

	m.methods["POST"] = func(r *http.Request) map[string][]string {
		// read body
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		// build kv maps
		kv := make(map[string][]string)
		query := strings.Split(string(body[:]), "&")
		for _, pairs := range query {
			pair := strings.Split(pairs, "=")
			if len(pair[0]) > 0 {
				if kv[pair[0]] != nil {
					kv[pair[0]] = append(kv[pair[0]], pair[1])
				} else {
					kv[pair[0]] = []string{pair[1]}
				}
			}
		}

		return kv
	}

	m.scenes = make(map[string]Scenes)
	m.scenes["publish"] = new(PublishScenes)
	m.scenes["publish_done"] = new(PublishDoneScenes)
	m.scenes["play"] = new(PlayScenes)
	m.scenes["update_publish"] = new(UpdatePublishScenes)

	m.streams = make(ServiceMapType)
	m.baseId = 1

	return m
}

func (manager *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := manager.baseId
	manager.baseId ++

	log.Printf("%d#0: *%d client connected, client: %s, server: %s, request: '%s %s %s'",
		os.Getpid(), id, r.RemoteAddr, manager.config.Listen,
		r.Method, r.URL.Path, r.Proto)

	if manager.config.Method != r.Method {
		log.Printf("%d#0: *%d method %s not matched, drop request %s, client: %s, server: %s, request: '%s %s %s'",
			os.Getpid(), id, r.Method, r.URL.Path,
		    r.RemoteAddr, manager.config.Listen,
		    r.Method, r.URL.Path, r.Proto)
		return
	}

	if form := manager.methods[r.Method](r); form != nil {
		if call := form["call"]; call != nil {
			if request := manager.scenes[call[0]]; request != nil {
				stream := NewStream(r, form)
				stream.streams = manager.streams

				// log.Info("service='%s' host='%s' app:='%s' name='%s'", stream.service, stream.host, stream.app, stream.name)
				scenes := request.New(manager)

				log.Printf("%d#0: *%d %s: service='%s' host='%s' app:='%s' name='%s', client: %s, server: %s, request: '%s %s %s'\n",
					os.Getpid(), id, scenes.Name(),
					stream.service, stream.host, stream.app, stream.name,
					r.RemoteAddr, manager.config.Listen,
					r.Method, r.URL.Path, r.Proto)

				scenes.Run(stream)
				scenes.Done(w)
			}
		}
	}
}
