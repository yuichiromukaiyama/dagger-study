package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

var (
	hostDir string = "../../"
	workDir string = "/app"
)

func main() {

	args := os.Args

	if len(args) != 2 {
		panic("no arguments specified")
	}

	registryPath := args[1]
	ctx := context.Background()
	client, err := dagger.Connect(context.Background(), dagger.WithLogOutput(os.Stdout))

	if err != nil {
		panic(err)
	}

	defer client.Close()

	// docker buildx imagetools inspect node:18-alpine
	platform := dagger.Platform("linux/amd64")
	filepath, _ := filepath.Abs(hostDir)

	test := client.Container(dagger.ContainerOpts{Platform: platform}).
		From("node:18-alpine").
		WithDirectory(workDir, client.Host().Directory(filepath)).
		WithWorkdir(workDir).
		WithExec([]string{"npm", "ci"}).
		WithExec([]string{"npm", "test"})

	_, err = test.Stderr(ctx)
	log, _ := test.Stdout(ctx)
	fmt.Println(log)

	if err != nil {
		panic(err)
	}

	ref := client.Container(dagger.ContainerOpts{Platform: platform}).
		From("node:18-alpine").
		WithDirectory(workDir, client.Host().Directory(filepath)).
		WithWorkdir(workDir).
		WithExec([]string{"npm", "ci"}).
		WithDefaultArgs(dagger.ContainerWithDefaultArgsOpts{Args: []string{"npm", "start"}}).
		WithExposedPort(3000)

	_, err = test.Stderr(ctx)
	log, _ = test.Stdout(ctx)
	fmt.Println(log)

	if err != nil {
		panic(err)
	}

	published, err := ref.Publish(ctx, registryPath)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Published image to: %s\n", published)
}
