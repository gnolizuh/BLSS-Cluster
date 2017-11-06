package main

import (
	"net/http"
)

type StreamMapType map[string]*Stream
type AppMapType map[string]StreamMapType
type ServiceMapType map[string]AppMapType
