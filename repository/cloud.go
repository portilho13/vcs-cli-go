package repository

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Clone(filePath string, url string) error {
    // Create the file
    out, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Check server response
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad status: %s", resp.Status)
    }

    // Get the content length for progress tracking
    contentLength := resp.ContentLength
    if contentLength <= 0 {
        return fmt.Errorf("unable to determine file size")
    }

    // Progress tracking variables
    var downloaded int64
    buf := make([]byte, 1024)
    start := time.Now()

    // Download the file with progress tracking
    for {
        n, err := resp.Body.Read(buf)
        if n > 0 {
            if _, err := out.Write(buf[:n]); err != nil {
                return err
            }
            downloaded += int64(n)
            printProgress(downloaded, contentLength, start)
        }
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }
    }

    fmt.Println("\nDownload completed!")
    return nil
}

// Function to print the download progress
func printProgress(downloaded, contentLength int64, start time.Time) {
    percent := float64(downloaded) / float64(contentLength) * 100
    elapsed := time.Since(start).Seconds()
    speed := float64(downloaded) / 1024 / elapsed
    fmt.Printf("\rDownloading... %.2f%% (%.2f KB/s)", percent, speed)
}

func GetServerStatus(ip string) error {
    _, err := http.Get(ip)
    if err != nil {
        return err
    }
    return nil
}

func UploadFile(filePath string, url string) error {
	// Open the .zlib file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
    io.Copy(part, file)
    writer.Close()

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodPost, url, body)
    req.Header.Add("Content-Type", writer.FormDataContentType())
	if err != nil {
		return err
	}

	// Create a HTTP client
	client := http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		// Read the response body for error details
		var responseBody []byte
		_, err := resp.Body.Read(responseBody)
		if err != nil && err != io.EOF {
			return err
		}
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, responseBody)
	}
    
	return nil
}