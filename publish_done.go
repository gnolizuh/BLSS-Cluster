package main

import (
	"net/http"
)

type PublishDoneScenes struct {
	BaseScenes
}

func (scenes *PublishDoneScenes) New(m *Manager) Scenes {
	p := new(PublishDoneScenes)
	p.status = http.StatusOK
	p.manager = m
	p.headers = make(map[string]string)
	return p
}

func (scenes *PublishDoneScenes) Name() string {
	return "PublishDoneScenes"
}

func (scenes *PublishDoneScenes) Run(s *Stream) {
	// TO DO
	// defer distlock.Unlock()
	// distlock.Lock()

	manager := scenes.manager
	redis_client := manager.redisClient

	// check stream already exist or not?
	if err := redis_client.Del(s.getUniqueKey()).Err(); err != nil {
		// already exist..
		scenes.status = http.StatusForbidden
	}
}
