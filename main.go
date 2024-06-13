package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/portilho13/vcs-cli-go/args"
	"github.com/portilho13/vcs-cli-go/cloud"
	"github.com/portilho13/vcs-cli-go/helpers"
	"github.com/portilho13/vcs-cli-go/repository"
)

type Repository = repository.Repository // type alias for repository.Repository type


var repo *Repository
var comment *string
var currentBranch repository.Branch

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
				var dirTree *repository.DirectoryTree
				if repo.Branch.DirTree == nil {
					dirTree, err = repository.CreateDirectoryTree(path, repo)
					if err != nil {
						fmt.Println("Error: ", err)
						return
					}
				}

				repo.Branch.DirTree = dirTree
				repository.PrintDirectoryTree(repo.Branch.DirTree, 0)
			}
		case "comment":
			if commandArgs[1] == "-m" {
				trimmed := strings.TrimSpace(commandArgs[2])
				comment = &trimmed
				fmt.Println("Comment: ", *comment)

			}
		case "clone":
			link := commandArgs[1]
			localPath, err := helpers.GetLocalPath()
			path = filepath.Join(localPath, "test.txt")
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			err = cloud.Clone(path, link)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
	}
}