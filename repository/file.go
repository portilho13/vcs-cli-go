package repository

import (
	"archive/tar"
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

func CompressFolder(source, destination string) error {
	// Create the destination file
	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Create a zlib writer
	zlibWriter := zlib.NewWriter(destFile)
	defer zlibWriter.Close()

	// Create a tar writer
	tarWriter := tar.NewWriter(zlibWriter)
	defer tarWriter.Close()

	// Walk through the source directory and add files to the tar writer
	return filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create the tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(file[len(source):])

		// Write the header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// If it's a directory, we don't need to write any file data
		if fi.Mode().IsRegular() {
			// Open the file
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			// Copy the file data to the tar writer
			if _, err := io.Copy(tarWriter, f); err != nil {
				return err
			}
		}

		return nil
	})
}
