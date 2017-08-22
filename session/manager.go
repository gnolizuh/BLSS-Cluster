package session

import (
	"net/http"
	"fmt"
	"../utils"
)

type Manager struct {
	config *utils.Config
}

func NewManager(config *utils.Config) (*Manager) {
	return &Manager{config}
}

func (manager *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
}
