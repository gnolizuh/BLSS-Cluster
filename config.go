package main

type Config struct {
	Listen    string `json:"listen"`
	Method    string `json:"method"`
	Log       string `json:"log"`
	RedisAddr string `json:"redis_addr"`
	Expire    int64  `json:"expire"`
	Password  string `json:"password"`
}