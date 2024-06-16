package repository

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func CreateDirectoryTree(path string, repo *Repository) (*DirectoryTree, error) {
	tree := NewDirectoryTree()

	info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    if !info.IsDir() {
        fileBlobHash, err := GenerateHashedFile(path, *repo)
        if err != nil {
            log.Printf("Error generating file blob hash for %s: %v", path, err)
            return nil, err
        }
        tree.Insert(info.Name(), DirectoryTree{File: &fileBlobHash}, path)
        return &tree, nil
    }

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
			fileBlobHash, err := GenerateHashedFile(fullPath, *repo)
			if err != nil {
				log.Printf("Error generating file blob hash for %s: %v", fullPath, err)
				continue
			}
			tree.Insert(name, DirectoryTree{File: &fileBlobHash}, fullPath)
		}
	}
	return &tree, nil
}

func UpdateDirectoryTree(path string, repo *Repository, tree *DirectoryTree) (*DirectoryTree, error) {
	info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    if !info.IsDir() {
        fileBlobHash, err := GenerateHashedFile(path, *repo)
        if err != nil {
            log.Printf("Error generating file blob hash for %s: %v", path, err)
            return nil, err
        }
        tree.Insert(info.Name(), DirectoryTree{File: &fileBlobHash}, path)
        return tree, nil
    }

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
			fileBlobHash, err := GenerateHashedFile(fullPath, *repo)
			if err != nil {
				log.Printf("Error generating file blob hash for %s: %v", fullPath, err)
				continue
			}
			tree.Insert(name, DirectoryTree{File: &fileBlobHash}, fullPath)
		}
	}
	return tree, nil
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
