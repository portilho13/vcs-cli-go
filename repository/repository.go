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
	Name string
	LocalPath string
	RemotePath string
	Branch string
}

func (r *Repository) Init(Name string, LocalPath string, RemotePath string, Branch string) (*Repository) {
	return &Repository{
		Name: Name,
		LocalPath: LocalPath,
		RemotePath: RemotePath,
		Branch: Branch,
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
	data := map[string]interface{}{
		"Name":    repo.Name,
		"LocalPath":     repo.LocalPath,
		"RemotePath": repo.RemotePath,
		"Branch":  repo.Branch,
	}

	// Convert the map to a JSON byte slice
	jsonData, err := json.MarshalIndent(data, "", "  ")
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