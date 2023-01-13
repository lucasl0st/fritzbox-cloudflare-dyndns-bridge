package main

import "github.com/cloudflare/cloudflare-go"

type User struct {
	ConfigUser
	cloudflare cloudflare.API
}
