package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var repoFolders = []string{
	".vcs/objects",
	".vcs/tree",
}

type Repository struct {
	Name       string
	LocalPath  string
	RemotePath string
    CurrentBranch string
	Branch     []Branch
}

func (r *Repository) Init(name, localPath, remotePath, branchName string) *Repository {
    var branches []Branch
	branch := Branch{
		Name:    branchName,
		DirTree: nil,
	}
    branches = append(branches, branch)
	return &Repository{
		Name:       name,
		LocalPath:  localPath,
		RemotePath: remotePath,
        CurrentBranch: branchName,
		Branch:     branches,
	}
}

func CreateRepoFolders() (error){
	for _, folder := range repoFolders {
		err := os.MkdirAll(folder, 0777)
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}
	}
	return nil
}

func RepoExists(path string) bool {
    repoPath := filepath.Join(path, ".vcs")
    _, err := os.Stat(repoPath)
    return !os.IsNotExist(err)
}

func SaveRepository(repo Repository) error {
    path := filepath.Join(repo.LocalPath, ".vcs", "info.json")

    // Convert the Repository struct to a JSON byte slice
    jsonData, err := json.MarshalIndent(repo, "", "  ")
    if err != nil {
        fmt.Println("Error marshalling to JSON:", err)
        return err
    }

    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.Write(jsonData)
    if err != nil {
        return err
    }
    return nil
}


func LoadRepository(path string) (*Repository, error) {
    var repo Repository
    path = filepath.Join(path, ".vcs", "info.json")
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&repo)
    if err != nil {
        return nil, err
    }
    return &repo, nil
}

func CheckRepository(path string) bool {
	repoPath := filepath.Join(path, ".vcs", "info.json")
	_, err := os.Stat(repoPath)
	return !os.IsNotExist(err)
}