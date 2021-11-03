package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/ksuid"
	"log"
	"math/rand"
	"time"

	"github.com/docker/docker/client"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// Custom config
	app := fiber.New(fiber.Config{
		AppName: "Give me a docker",
	})

	app.Get("/give", func(c *fiber.Ctx) error {
		//Generate safe port number from range 30000-32767
		//port := 30000 + rand.Intn(2767)
		//portStr := strconv.Itoa(port)

		kid := ksuid.New().String()
		//kid := "kabaca"

		// Create a container
		ctx := context.Background()
		create, err := cli.ContainerCreate(ctx,
			&container.Config{
				Cmd:   []string{"dockerd", "--host=tcp://0.0.0.0:2375"},
				Image: "docker:dind",
				Labels: map[string]string{
					"traefik.enable": "true",
					"traefik.http.routers." + kid + ".rule":                      "Host(`" + kid + ".theykk.com`)",
					"traefik.http.services." + kid + ".loadbalancer.server.port": "2375",
				},
			},
			&container.HostConfig{
				Privileged: true,
				// Do not open port to public, use proxy instead
				//PortBindings: nat.PortMap{
				//	"2375/tcp": {
				//		{
				//			HostPort: portStr,
				//		},
				//	},
				//},
			},
			nil,
			nil,
			"dockerverbana_"+kid)
		if err != nil {
			log.Println(err)
		}
		log.Println(create)
		err = cli.ContainerStart(ctx, create.ID, types.ContainerStartOptions{})
		if err != nil {
			log.Println(err)
		}
		return c.SendString(kid)
	})
	log.Fatal(app.Listen(":8585"))
}
