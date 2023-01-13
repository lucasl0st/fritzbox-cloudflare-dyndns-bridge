package main

import (
	"context"
	"encoding/json"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open(ConfigFile)

	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(f)

	if err != nil {
		log.Fatal(err)
	}

	var config Config

	err = json.Unmarshal(b, &config)

	if err != nil {
		log.Fatal(err)
	}

	users := make(map[string]User)

	for _, user := range config.Users {
		api, err := cloudflare.NewWithAPIToken(user.CloudflareApiKey)

		if err != nil {
			log.Fatal(err)
		}

		users[user.Username] = User{
			ConfigUser: user,
			cloudflare: *api,
		}
	}

	api := Api{
		users: users,
		ctx:   context.Background(),
		cache: make(map[string]string),
	}

	r := gin.Default()

	api.RegisterEndpoints(r)

	err = r.Run()

	if err != nil {
		log.Fatal(err)
	}
}
