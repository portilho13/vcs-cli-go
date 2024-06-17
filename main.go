package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/portilho13/vcs-cli-go/args"
	"github.com/portilho13/vcs-cli-go/helpers"
	"github.com/portilho13/vcs-cli-go/repository"
)

type Repository = repository.Repository // type alias for repository.Repository type

var repo *Repository
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
		var currentBranch repository.Branch
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if repository.RepoExists(localPath) {
			currentBranch, err = repository.GetCurrentBranch(*repo)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			var dirTree *repository.DirectoryTree
			if currentBranch.DirTree == nil {
				dirTree, err = repository.CreateDirectoryTree(path, repo)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
			} else {
				fmt.Println("Updating directory tree")
				dirTree, err = repository.UpdateDirectoryTree(path, repo, currentBranch.DirTree)
				fmt.Println("Directory tree updated")
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
			}

			currentBranch.DirTree = dirTree
			err = repository.UpdateRepoBranch(repo, currentBranch)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			err = repository.SaveRepository(*repo)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}
	case "commit":
		if commandArgs[1] == "-m" {
			trimmed := strings.TrimSpace(commandArgs[2])
			commit := trimmed
			fmt.Println("Comment: ", commit)

		}
	case "clone":
		link := commandArgs[1]
		localPath, err := helpers.GetLocalPath()
		path = filepath.Join(localPath, "test.txt")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		err = repository.Clone(path, link)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	case "origin":
		remotePath := commandArgs[1]
		repo.RemotePath = remotePath
		if repository.RepoExists(path) {
			err = repository.SaveRepository(*repo)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}
	case "push":
		localPath, err := helpers.GetLocalPath()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		sourceDir := localPath + "/.vcs"
		targetFile := repo.Name + ".zlib"

		fmt.Println("Compressing folder...: ", targetFile)

		err = repository.CompressFolder(sourceDir, targetFile)
		if err != nil {
			fmt.Println("Error compressing folder:", err)
		} else {
			fmt.Println("Folder compressed successfully!")
		}

		err = repository.UploadFile(targetFile, repo.RemotePath)
		if err != nil {
			fmt.Println("Error uploading file:", err)
		} else {
			fmt.Println("File uploaded successfully!")
		}
	case "checkout":
		branch := commandArgs[1]
		localPath, err := helpers.GetLocalPath()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		if repository.RepoExists(localPath) {
			fmt.Println("Checking out branch: ", branch)
			err = repository.UpdateBranchDir(repo, branch, localPath)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			err = repository.SaveRepository(*repo)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}

	default:
		fmt.Println("Invalid command")

	}
	if repository.RepoExists(path) {
		err = repository.SaveRepository(*repo)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
}
