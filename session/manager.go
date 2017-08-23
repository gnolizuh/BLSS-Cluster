package session

import (
	"net/http"
	"fmt"
	"../utils"
	"io/ioutil"
	"strings"
)

type FuncType func(*http.Request) map[string][]string

type Manager struct {
	config *utils.Config
	funcMap map[string]FuncType
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
			kv[pair[0]] = []string{pair[1]}
		}

		return kv
	}

	return m
}

func (manager *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if manager.config.Method != r.Method {
		fmt.Printf("method %s not matched, drop request %s\n",
			r.Method, r.URL.Path)
		return
	}

	if form := manager.funcMap[r.Method](r); form != nil {
		fmt.Println(form)
	}
}
