package main

import (
	"fmt"

	"github.com/portilho13/vcs-cli-go/args"
	"github.com/portilho13/vcs-cli-go/file"
	"github.com/portilho13/vcs-cli-go/helpers"
	"github.com/portilho13/vcs-cli-go/repository"
	"github.com/portilho13/vcs-cli-go/tree"
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
			err = repository.CreateRepoFolders()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			content, err := file.GenerateHashedFile("/Users/marioportilho/Desktop/Coding/vcs-cli-go/teste/teste.c", repo)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			fmt.Println(content)
		case "add":
			path := commandArgs[1]
			localPath, err := helpers.GetLocalPath()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			if repository.RepoExists(localPath) {
				dtree, err := tree.CreateDirectoryTree(path, repo)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}

				tree.PrintDirectoryTree(dtree, 0)
			}
	}
}