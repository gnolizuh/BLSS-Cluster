package session

import (
	"net/http"
	"fmt"
	"../utils"
	"../scenes"
	"io/ioutil"
	"strings"
)

type Scenes interface {
	New() interface{}
	Run()
}

type FuncType func(*http.Request) map[string][]string

type Manager struct {
	config *utils.Config
	funcMap map[string]FuncType
	scenesMap map[string]interface{}
}

func NewManager(config *utils.Config) (*Manager) {
	m := new(Manager)
	m.config = config
	m.funcMap = make(map[string]FuncType)
	m.funcMap["GET"] = func(r *http.Request) map[string][]string {
		return r.URL.Query()
	}

	m.funcMap["POST"] = func(r *http.Request) map[string][]string {
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

	m.scenesMap = make(map[string]interface{})
	m.scenesMap["publish"] = scenes.PublishScenes{}

	return m
}

func (manager *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if manager.config.Method != r.Method {
		fmt.Printf("method %s not matched, drop request %s\n",
			r.Method, r.URL.Path)
		return
	}

	if form := manager.funcMap[r.Method](r); form != nil {
		if call := form["call"]; call != nil {
			fmt.Println(call)
			if scenes := manager.scenesMap[call[0]]; scenes != nil {
				(((scenes.(Scenes)).New()).(Scenes)).Run()
			}
		}
	}
}
