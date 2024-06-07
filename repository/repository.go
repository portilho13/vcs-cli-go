package repository

import (
	"fmt"
	"os"
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

func (r *Repository) Init(Name string, LocalPath string, RemotePath string, Branch string) (Repository) {
	return Repository{
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