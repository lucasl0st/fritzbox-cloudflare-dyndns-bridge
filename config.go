package main

const ConfigFile = "config.json"

type Config struct {
	Users []ConfigUser `json:"users"`
}

type ConfigUser struct {
	Username         string   `json:"username"`
	Password         string   `json:"password"`
	AllowedDomains   []string `json:"allowedDomains"`
	CloudflareApiKey string   `json:"cloudflareApiKey"`
}
