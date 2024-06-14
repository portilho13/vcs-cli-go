package repository

import (
	"compress/zlib"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ConvertToBin reads the file at filePath and converts its content to a binary string representation.
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

// GenerateHash256 creates a SHA-256 hash of the provided content.
func GenerateHash256(content string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(content))
	if err != nil {
		return "", err
	}
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs), nil
}

// GenerateHashedFile generates a binary representation of the file at filePath,
// hashes the content, and stores it in the repository's .vcs/objects directory.
func GenerateHashedFile(filePath string, repo Repository) (string, error) {
	content, err := ConvertToBin(filePath)
	if err != nil {
		return "", err
	}
	hash, err := GenerateHash256(content)
	if err != nil {
		return "", err
	}

	genFilePath := filepath.Join(repo.LocalPath, ".vcs", "objects", hash)

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

func CompressFolder(sourceDir, targetDir string) error {

	outFile, err := os.Create(targetDir)
	if err != nil {
		return err
	}
	
	defer outFile.Close()

	zw := zlib.NewWriter(outFile)

	defer zw.Close()


	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		inFile, err := os.Open(path)
		if err != nil {
			return err
		}

		defer inFile.Close()

		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		_, err = zw.Write([]byte(relativePath + "\n"))
		if err != nil {
			return err
		}

		_, err = io.Copy(zw, inFile)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = zw.Flush()
	if err != nil {
		return err
	}

	return nil
}
