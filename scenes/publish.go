package scenes

import (
	"net/http"
	"time"
)

type PublishScenes struct {
	BaseScenes
}

func (scenes *PublishScenes) New(m *Manager) Scenes {
	p := new(PublishScenes)
	p.status = http.StatusOK
	p.manager = m
	p.headers = make(map[string]string)
	return p
}

func (scenes *PublishScenes) Name() string {
	return "PublishScenes"
}

func (scenes *PublishScenes) Run(s *Stream) {
	// TO DO
	// defer distlock.Unlock()
	// distlock.Lock()

	manager := scenes.manager
	redis_client := manager.redisClient
	expire := time.Duration(manager.config.Expire) * time.Second

	// check stream already exist or not?
	set, err := redis_client.SetNX(s.getUniqueKey(), s.getUniqueVal(), expire).Result()
	if err != nil {
		// something error.
		panic(err)
	} else if !set {
		// already exist.
		scenes.status = http.StatusForbidden
	} else {
		// OK!
	}
}