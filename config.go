package main

type Config struct {
	Listen   string `json:"listen"`
	Method   string `json:"method"`
	Log      string `json:"log"`
	Redis    string `json:"redis"`
	Expire   int64  `json:"expire"`
	Password string `json:"password"`
}