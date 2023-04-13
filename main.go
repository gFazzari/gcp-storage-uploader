package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var uploader *ClientUploader
var filename string
var serviceAccount string
var projectID string
var bucketName string

func main() {
	flag.StringVar(&filename, "filename", "/home/gfazzari/sealed-secret/sealed-secret.yml", "Path to file which has to be uploaded")
	flag.StringVar(&serviceAccount, "service_account", "/service-account/sealed-secret.json", "Path to the service account")
	flag.StringVar(&projectID, "projectID", "testingGCP", "GCP's Project ID")
	flag.StringVar(&bucketName, "bucketName", "uploads", "GCP's bucket name")

	flag.Parse()

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", serviceAccount)

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}

	uploader = &ClientUploader{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
		uploadPath: "uploads/",
	}

	err = uploader.UploadFile(filename)
	if err != nil {
		log.Fatalf("Cannot upload file: %v", err)
	}

	log.Println("Successfully upload file.")
}

// UploadFile uploads an object
func (c *ClientUploader) UploadFile(filename string) error {
	// Prepare context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Read file
	file, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		log.Fatalf("File cannot be read: %v", err)
	}

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + filepath.Base(filename)).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("error during copy --> %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("error during closing --> %v", err)
	}

	// All ok
	return nil
}
