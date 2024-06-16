package repository

import (
	"archive/tar"
	"compress/zlib"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func GenerateHash256(content string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(content))
	if err != nil {
		return "", err
	}
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs), nil
}

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

func DecompressZlibFile(inputPath, outputDir string) error {
	// Open the compressed file
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// Create a zlib reader
	zlibReader, err := zlib.NewReader(inFile)
	if err != nil {
		return err
	}
	defer zlibReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(zlibReader)

	outputDir = strings.TrimSuffix(outputDir, filepath.Ext(outputDir))

	tempOutputDir := outputDir

	outputDir = outputDir + "/.vcs"

	// Ensure the destination directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Extract files from the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		// Remove the .zip extension from the header name

		// Determine the file path
		filePath := filepath.Join(outputDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(filePath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// Create the file
			outFile, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			// Copy the file data
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported file type: %v", header.Typeflag)
		}
	}
	repo, err := LoadRepository(tempOutputDir)
	if err != nil {
		return err
	}

	err = MountRepositoryFolder(repo, tempOutputDir)
	if err != nil {
		return err
	}

	return nil
}

func GetFileContent(fileName string, localPath string) (string, error) {
	filePath := filepath.Join(localPath, ".vcs", "objects", fileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	decodedContent, err := ConvertFromBin(string(content))
	if err != nil {
		return "", fmt.Errorf("error converting from binary: %v", err)
	}

	return decodedContent, nil
}

func ConvertFromBin(binaryString string) (string, error) {
    // Split the binary string into individual binary representations of bytes
    binaryBytes := strings.Fields(binaryString)

    var content []byte
    for _, binByte := range binaryBytes {
        if len(binByte) != 8 {
            return "", fmt.Errorf("invalid binary string format: %s", binByte)
        }

        // Parse binary string to integer
        intVal, err := strconv.ParseInt(binByte, 2, 64)
        if err != nil {
            return "", fmt.Errorf("error parsing binary string: %v", err)
        }

        // Append byte value to content
        content = append(content, byte(intVal))
    }

    // Convert byte slice to string
    return string(content), nil
}






func mountDirectoryTree(repoPath string, basePath string, tree *DirectoryTree) error {
	for _, subtree := range tree.Directory {
		fullPath := filepath.Join(basePath, subtree.Path)
		if subtree.Tree.File != nil {
			// It's a file
			dir := filepath.Dir(fullPath)
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directory %s: %v", dir, err)
			}

			content, err := GetFileContent(*subtree.Tree.File, repoPath)
			if err != nil {
				return fmt.Errorf("error getting file content: %v", err)
			}

			file, err := os.Create(fullPath)
			if err != nil {
				return fmt.Errorf("error creating file %s: %v", fullPath, err)
			}
			defer file.Close()

			_, err = file.WriteString(content)
			if err != nil {
				return fmt.Errorf("error writing file %s: %v", fullPath, err)
			}
		} else {
			// It's a directory
			err := os.MkdirAll(fullPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating directory %s: %v", fullPath, err)
			}
			err = mountDirectoryTree(repoPath, fullPath, &subtree.Tree)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func MountRepositoryFolder(repo *Repository, localPath string) error {
	currentBranchName := repo.CurrentBranch

	var currentBranch *Branch
	for _, branch := range repo.Branch {
		if branch.Name == currentBranchName {
			currentBranch = &branch
			break
		}
	}

	if currentBranch == nil {
		return fmt.Errorf("Branch %s not found", currentBranchName)
	}

	dirTree := currentBranch.DirTree
	if dirTree == nil {
		return fmt.Errorf("Branch %s has no directory tree", currentBranchName)
	}
	err := mountDirectoryTree(localPath, localPath, dirTree)
	if err != nil {
		return fmt.Errorf("error mounting directory tree: %v", err)
	}

	return nil
}