package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"os"
)

type Scenes interface {
	New() interface{}
	Name() string
	Run(stream *Stream, w http.ResponseWriter, r *http.Request)
}

type MethodFuncType func(*http.Request) map[string][]string

type Manager struct {
	config    *Config
	methodFunc map[string]MethodFuncType
	scenesFunc map[string]interface{}
	streams    ServiceMapType
	baseId     uint64
}

func NewManager(config *Config) (*Manager) {
	m := new(Manager)
	m.config = config
	m.methodFunc = make(map[string]MethodFuncType)
	m.methodFunc["GET"] = func(r *http.Request) map[string][]string {
		return r.URL.Query()
	}

	m.methodFunc["POST"] = func(r *http.Request) map[string][]string {
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

	m.scenesFunc = make(map[string]interface{})
	m.scenesFunc["publish"] = PublishScenes{}
	m.scenesFunc["publish_done"] = PublishDoneScenes{}
	m.scenesFunc["play"] = PlayScenes{}
	m.scenesFunc["update_publish"] = UpdatePublishScenes{}

	m.streams = make(ServiceMapType)
	m.baseId = 1

	return m
}

func (manager *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := manager.baseId
	manager.baseId ++

	fmt.Printf("%s %d#0: *%d client connected, client: %s, server: %s, request: '%s %s %s'\n",
		time.Now().Format("2006/01/02 15:04:05"), os.Getpid(), id,
		r.RemoteAddr, manager.config.Listen,
		r.Method, r.URL.Path, r.Proto)

	if manager.config.Method != r.Method {
		fmt.Printf("%s %d#0: *%d method %s not matched, drop request %s, client: %s, server: %s, request: '%s %s %s'\n",
			time.Now().Format("2006/01/02 15:04:05"), os.Getpid(), id, r.Method, r.URL.Path,
		    r.RemoteAddr, manager.config.Listen,
		    r.Method, r.URL.Path, r.Proto)
		return
	}

	if form := manager.methodFunc[r.Method](r); form != nil {
		if call := form["call"]; call != nil {
			if scenes := manager.scenesFunc[call[0]]; scenes != nil {
				stream := NewStream(r, form)
				stream.streams = manager.streams

				fmt.Printf("%s %d#0: *%d %s: service='%s' host='%s' app:='%s' name='%s', client: %s, server: %s, request: '%s %s %s'\n",
					time.Now().Format("2006/01/02 15:04:05"), os.Getpid(), id, (scenes.(Scenes)).Name(),
					stream.service, stream.host, stream.app, stream.name,
					r.RemoteAddr, manager.config.Listen,
					r.Method, r.URL.Path, r.Proto)

				(((scenes.(Scenes)).New()).(Scenes)).Run(stream, w, r)
			}
		}
	}
}
