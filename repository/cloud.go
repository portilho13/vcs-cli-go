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
    
    out, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer out.Close()

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

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
    cd := resp.Header.Get("Content-Disposition")
    if cd != "" {
        _, params, err := mime.ParseMediaType(cd)
        if err == nil {
            if filename, ok := params["filename"]; ok {
                return filename
            }
        }
    }

    u, err := url.Parse(downloadURL)
    if err != nil {
        return "unknown_filename"
    }

    return path.Base(u.Path)
}


func UploadFile(filePath string, url string) error {
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

	req, err := http.NewRequest(http.MethodPost, url, body)
    req.Header.Add("Content-Type", writer.FormDataContentType())
	if err != nil {
		return err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad status: %s", resp.Status)
	}
    
	return nil
}