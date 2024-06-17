package repository

import (
	"errors"
)

type Branch struct {
	Name    string
	DirTree *DirectoryTree
}

func GetCurrentBranch(repo Repository) (Branch, error) {
	for _, branch := range repo.Branch {
		if branch.Name == repo.CurrentBranch {
			return branch, nil
		}
	}
	return Branch{}, errors.New("Branch not found")
}

func UpdateRepoBranch(repo *Repository, branch Branch) error {
	for i, b := range repo.Branch {
		if b.Name == branch.Name {
			repo.Branch[i] = branch
			return nil
		}
	}
	return errors.New("Branch not found")
}

func UpdateBranchDir(repo *Repository, branch string, localPath string) error {
	var currentBranch Branch
	branchFound := false

	for i := range repo.Branch {
		if repo.Branch[i].Name == branch {
			currentBranch = repo.Branch[i]
			branchFound = true
			break
		}
	}

	if !branchFound {
		currentBranch = CreateBranch(repo, branch)
	}

	err := removeAllEntries(localPath)
	if err != nil {
		return err
	}

	repo.CurrentBranch = currentBranch.Name

	if currentBranch.DirTree != nil {
		err := MountRepositoryFolder(repo, localPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateBranch(repo *Repository, branchName string) Branch {
	branch := Branch{Name: branchName, DirTree: nil}
	repo.Branch = append(repo.Branch, branch)
	return branch
}
