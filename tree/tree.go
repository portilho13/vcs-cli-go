package tree

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/portilho13/vcs-cli-go/file"
	"github.com/portilho13/vcs-cli-go/repository"
)

var FOLDERS_TO_IGNORE = []string{"target", ".git", ".vcs"}
var FILES_TO_IGNORE = []string{"Cargo.lock", "Cargo.toml", ".gitignore", "README.md"}

type DirectoryTree struct {
	File      *string
	Directory map[string]Subtree
}

type Subtree struct {
	Tree DirectoryTree
	Path string
}

func NewDirectoryTree() DirectoryTree {
	return DirectoryTree{
		Directory: make(map[string]Subtree),
	}
}

func (dt *DirectoryTree) Insert(name string, subtree DirectoryTree, path string) {
	dt.Directory[name] = Subtree{
		Tree: subtree,
		Path: path,
	}
}

func CreateDirectoryTree(path string, repo *repository.Repository) (*DirectoryTree, error) {
	tree := NewDirectoryTree()
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(path, name)

		if toIgnore(name) {
			continue
		}

		if entry.IsDir() {
			subtree, err := CreateDirectoryTree(fullPath, repo)
			if err != nil {
				return nil, err
			}
			tree.Insert(name, *subtree, fullPath)
		} else {
			fileBlobHash, err := file.GenerateHashedFile(fullPath, *repo)
			if err != nil {
				log.Printf("Error generating file blob hash for %s: %v", fullPath, err)
				continue
			}
			tree.Insert(name, DirectoryTree{File: &fileBlobHash}, fullPath)
		}
	}
	return &tree, nil
}

func PrintDirectoryTree(tree *DirectoryTree, level int) {
	for name, subtree := range tree.Directory {
		for i := 0; i < level; i++ {
			fmt.Print("  ")
		}
		fmt.Println(name)
		fmt.Println("  Path:", subtree.Path)
		PrintDirectoryTree(&subtree.Tree, level+1)
	}

	if tree.File != nil {
		for i := 0; i < level; i++ {
			fmt.Print("  ")
		}
		fmt.Printf("File -> %s\n", *tree.File)
	}
}

func toIgnore(name string) bool {
	for _, f := range FOLDERS_TO_IGNORE {
		if f == name {
			return true
		}
	}
	for _, f := range FILES_TO_IGNORE {
		if f == name {
			return true
		}
	}
	return false
}
