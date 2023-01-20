package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const baseApi = "https://api.chainsafe.io/api/v1"

func main() {
	bucketID := os.Getenv("CS_STORAGE_BUCKET_ID")
	if bucketID == "" {
		log.Panicf("CS Storage bucket id missing")
	}

	storageToken := os.Getenv("CS_STORAGE_API_KEY")
	if storageToken == "" {
		log.Panicf("CS Storage api key missing")
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Panicf("can't get PWD: %s", err)
	}

	// configuration file path
	env := os.Getenv("ENV")
	fileName := fmt.Sprintf("config.%s.json", strings.ToLower(env))
	filePath := filepath.Join(pwd, fileName)

	err = uploadToBucket(bucketID, fileName, filePath, storageToken)
	if err != nil {
		log.Panicf("Error on uploading to CS storage: %v", err)
	}

	cid, err := getFileCid(bucketID, fileName, storageToken)
	if err != nil {
		log.Panicf("Error on uploading to CS storage: %v", err)
	}

	fmt.Printf("Uploaded %s to CS Storage, file CID: %s\n", fileName, cid)
}

type requestBody struct {
	Path   string `json:"path"`
	Source string `json:"source"`
}

func getFileCid(bucketID string, fileName string, bearerToken string) (string, error) {
	url := fmt.Sprintf("%s/bucket/%s/file", baseApi, bucketID)
	var err error

	client := &http.Client{}

	// Create the JSON encoded body
	body := &requestBody{Path: fileName, Source: bucketID}
	jsonValue, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	// Add headers to the request
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+bearerToken)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var fileInfo FileInfo
	err = json.NewDecoder(resp.Body).Decode(&fileInfo)
	if err != nil {
		return "", err
	}

	return fileInfo.Content.Cid, err
}

type FileInfo struct {
	Content struct {
		Name        string `json:"name"`
		Cid         string `json:"cid"`
		Size        int    `json:"size"`
		ContentType string `json:"content_type"`
		Version     int    `json:"version"`
		CreatedAt   int    `json:"created_at"`
	} `json:"content"`
	Persistent struct {
		Uploaded   int      `json:"uploaded"`
		SavedTime  int      `json:"saved_time"`
		StoredSize int      `json:"stored_size"`
		StoredCid  string   `json:"stored_cid"`
		Filecoin   struct{} `json:"filecoin"`
	} `json:"persistent"`
	Messages []string `json:"messages"`
}

func uploadToBucket(bucketID string, fileName string, filePath string, bearerToken string) error {
	url := fmt.Sprintf("%s/bucket/%s/upload", baseApi, bucketID)
	var err error

	// Open the file to be uploaded
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new buffer to send the file in
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a new field for the file path
	field, err := writer.CreateFormField("path")
	if err != nil {
		return err
	}
	_, err = field.Write([]byte(""))
	if err != nil {
		return err
	}

	// Create a new part for the file
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}

	// Copy the file into the part
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return err
	}

	resp, err := executeRequest("POST", url, body, writer.FormDataContentType(), bearerToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func executeRequest(
	method string, url string, body *bytes.Buffer, contentType string, bearerToken string,
) (*http.Response, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new POST request with the file as the body
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Add any necessary headers to the request
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+bearerToken)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}
