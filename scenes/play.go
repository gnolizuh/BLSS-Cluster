package scenes

import (
	"net/http"
	"github.com/go-redis/redis"
	"fmt"
)

type PlayScenes struct {
	BaseScenes
}

func (scenes *PlayScenes) NewPlayScenes(m *Manager) *BaseScenes {
	p := new(PlayScenes)
	p.status = http.StatusOK
	p.manager = m
	p.headers = make(map[string]string)
	return p
}

func (scenes *PlayScenes) Name() string {
	return "PlayScenes"
}

func (scenes *PlayScenes) Run(s *Stream) {
	// TO DO
	// defer distlock.Unlock()
	// distlock.Lock()

	manager := scenes.manager
	redis_client := manager.redisClient

	// get stream info.
	addr, err := redis_client.Get(s.getUniqueKey()).Result()
	if err == redis.Nil {
		scenes.status = http.StatusNotFound
	} else if err != nil {
		scenes.status = http.StatusNotFound
		panic(err)
	} else {
		daddr := fmt.Sprintf("%s:%d", s.localAddr, s.localPort)
		if addr != daddr {
			scenes.status = http.StatusFound
			scenes.headers["Location"] = fmt.Sprintf("rtmp://%s/%s/%s/%s", addr, s.host, s.app, s.name)
		}
	}
}
