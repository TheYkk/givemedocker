package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2"
	"log"
	"math/rand"
	"strconv"
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
		port := 30000 + rand.Intn(2767)
		// Create a container
		ctx := context.Background()
		create, err := cli.ContainerCreate(ctx,
			&container.Config{
				Cmd:   []string{"dockerd", "--host=tcp://0.0.0.0:2375"},
				Image: "docker:dind",
			},
			&container.HostConfig{
				Privileged: true,
				PortBindings: nat.PortMap{
					"2375/tcp": {
						{
							HostPort: strconv.Itoa(port),
						},
					},
				},
			},
			nil,
			nil,
			"dockerverbana"+strconv.Itoa(port))
		if err != nil {
			log.Println(err)
		}
		log.Println(create)
		err = cli.ContainerStart(ctx, create.ID, types.ContainerStartOptions{})
		if err != nil {
			log.Println(err)
		}
		return c.SendString(strconv.Itoa(port))
	})
	log.Fatal(app.Listen(":8585"))
}
