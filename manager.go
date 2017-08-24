package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
)

type Scenes interface {
	New() interface{}
	Run(stream *Stream)
}

type MethodFuncType func(*http.Request) map[string][]string

type Manager struct {
	config *Config
	methodFunc map[string]MethodFuncType
	scenesFunc map[string]interface{}
	streams ServiceMapType
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

	m.streams = make(ServiceMapType)

	return m
}

func (manager *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if manager.config.Method != r.Method {
		fmt.Printf("method %s not matched, drop request %s\n",
			r.Method, r.URL.Path)
		return
	}

	if form := manager.methodFunc[r.Method](r); form != nil {
		fmt.Println(form)
		if call := form["call"]; call != nil {
			if scenes := manager.scenesFunc[call[0]]; scenes != nil {
				data := ScenesDataType{r, form}
				stream := NewStream(data)
				stream.streams = manager.streams

				(((scenes.(Scenes)).New()).(Scenes)).Run(stream)
				fmt.Println(manager.streams)
			}
		}
	}
}
