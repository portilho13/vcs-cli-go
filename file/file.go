package file

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/portilho13/vcs-cli-go/repository"
)

func ConvertToBin(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	var result string
    for _, c := range content {
        result += fmt.Sprintf("%08b ", c)
    }
    return result, nil

}

func GenerateHash256(content string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(content))
	if err != nil {
		return "", err
	}
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs), nil
}

func GenerateHashedFile(filePath string, repo repository.Repository) (string, error) {
	content, err := ConvertToBin(filePath)
	if err != nil {
		return "", err
	}
	hash, err := GenerateHash256(content)
	if err != nil {
		return "", err
	}

	genFilePath := repo.LocalPath + "/.vcs/objects/" + hash

	dir := filepath.Dir(genFilePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
    }
	
	file, err := os.OpenFile(genFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return "", err
	}

	return hash, nil
}