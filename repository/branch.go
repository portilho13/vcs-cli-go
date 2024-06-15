package repository

import "errors"

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

func UpdateRepoBranch(repo Repository, branch Branch) (*Repository, error) {
	for i, b := range repo.Branch {
		if b.Name == branch.Name {
			repo.Branch[i] = branch
			return &repo, nil
		}
	}
	return &Repository{}, errors.New("Branch not found")
}