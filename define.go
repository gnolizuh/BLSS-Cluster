package main

import (
	"net/http"
)

type StreamMapType map[string]*Stream
type AppMapType map[string]StreamMapType
type ServiceMapType map[string]AppMapType

type ScenesDataType struct {
	request *http.Request
	form map[string][]string
}
