package repository

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
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

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    fileName := getFileName(resp, url)

    fileLoc := filepath.Join("./", fileName)

    err = DecompressZlibFile(filePath, fileLoc)
    if err != nil {
        return err
    }


    

    return nil
}

func getFileName(resp *http.Response, downloadURL string) string {
    // Check Content-Disposition header
    cd := resp.Header.Get("Content-Disposition")
    if cd != "" {
        _, params, err := mime.ParseMediaType(cd)
        if err == nil {
            if filename, ok := params["filename"]; ok {
                return filename
            }
        }
    }

    // Fallback to extracting from URL
    u, err := url.Parse(downloadURL)
    if err != nil {
        return "unknown_filename"
    }

    return path.Base(u.Path)
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