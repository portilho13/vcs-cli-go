package main

import (
	"fmt"
	"strings"

	"github.com/portilho13/vcs-cli-go/args"
	"github.com/portilho13/vcs-cli-go/helpers"
	"github.com/portilho13/vcs-cli-go/repository"
	"github.com/portilho13/vcs-cli-go/tree"
)

type Repository = repository.Repository // type alias for repository.Repository type


var repo *Repository
var dirTree *tree.DirectoryTree
var comment *string

func main() {
	commandArgs := args.GetArgs()
	path, err := helpers.GetLocalPath()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	if repository.CheckRepository(path) {
		repo, err = repository.LoadRepository(path)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
	switch commandArgs[0] {
		case "init": 
			if repository.CheckRepository(path) {
				fmt.Println("Repository already exists")
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

			if !repository.CheckRepository(path) {
				err = repository.SaveRepository(*repo)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
			}
			fmt.Println(repo)

		case "add":
			path := commandArgs[1]
			localPath, err := helpers.GetLocalPath()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			if repository.RepoExists(localPath) {
				dirTree, err = tree.CreateDirectoryTree(path, repo)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}

				tree.PrintDirectoryTree(dirTree, 0)
			}
		case "comment":
			if commandArgs[1] == "-m" {
				trimmed := strings.TrimSpace(commandArgs[2])
				comment = &trimmed
				fmt.Println("Comment: ", *comment)

			}
	}
}