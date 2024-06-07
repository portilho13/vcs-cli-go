package helpers

import (
	"os"
	"strings"
)

func GetLocalPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func GetPathName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	splitedDir := strings.Split(dir, "/")
	return splitedDir[len(splitedDir)-1], nil
}