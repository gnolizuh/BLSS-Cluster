package scenes

import (
	"net/http"
	"time"
)

type UpdatePublishScenes struct {
	BaseScenes
}

func (scenes *UpdatePublishScenes) New(m *Manager) Scenes {
	p := new(UpdatePublishScenes)
	p.status = http.StatusOK
	p.manager = m
	p.headers = make(map[string]string)
	return p
}

func (scenes *UpdatePublishScenes) Name() string {
	return "UpdatePublishScenes"
}

func (scenes *UpdatePublishScenes) Run(s *Stream) {
	// TO DO
	// defer distlock.Unlock()
	// distlock.Lock()

	manager := scenes.manager
	redis_client := manager.redisClient
	expire := time.Duration(manager.config.Expire) * time.Second

	// keep stream alive.
	set, err := redis_client.Expire(s.getUniqueKey(), expire).Result()
	if err != nil {
		// something error.
		panic(err)
	} else if !set {
		// stream was expired.
		scenes.status = http.StatusForbidden
	} else {
		// OK!
	}
}