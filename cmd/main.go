package main

import (
	"fmt"
	g "github.com/brandonlbarrow/github-go/internal/github"
	"os"
)

func main() {

	creds := g.Credentials{
		Username:    os.Getenv("USERNAME"),
		Password:    os.Getenv("PASSWORD"),
		Multifactor: true,
	}
	basic, err := g.NewHttpBasicProvider(&creds)
	if err != nil {
		_ = fmt.Errorf("failed to initialize provider: %w", err)
		os.Exit(1)
	}

	client := basic.Auth()

	repos, _, err := client.ListRepos()
	if err != nil {
		_ = fmt.Errorf("failed to list repos: %w", err)
		os.Exit(1)
	}
	fmt.Println(len(repos))
}

