package main

import (
	"fmt"

	"github.com/portilho13/vcs-cli-go/args"
	"github.com/portilho13/vcs-cli-go/helpers"
	"github.com/portilho13/vcs-cli-go/repository"
)

type Repository = repository.Repository // type alias for repository.Repository type


var repo Repository

func main() {
	fmt.Println("Hello, World!")
	commandArgs := args.GetArgs()
	switch commandArgs[0] {
		case "init":
			path, err := helpers.GetLocalPath()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			repoName, err := helpers.GetPathName()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			repo = repo.Init(repoName, path, "test", "master")
			fmt.Println(repo)
			err = repository.CreateRepoFolders()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
	}
}