package storage_save_slack

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name string `json:"name"`
}

func StorageImageSend(ctx context.Context, e GCSEvent) error {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("twitterapi-key.json"))
	if err != nil {
		log.Fatal(err)
	}

	obj := client.Bucket(os.Getenv("BUCKET_NAME")).Object(e.Name)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	log.Printf("%s : %s\n", obj.BucketName(), obj.ObjectName())
	slackPost(obj.ObjectName(), reader)
	return nil
}

func slackPost(fileName string, reader io.Reader) {
	url := "https://slack.com/api/files.upload?token=" + os.Getenv("SLACK_ACCESS_TOKEN") + "&channels=" + os.Getenv("SLACK_CHANNEL_ID")
	body := &bytes.Buffer{}

	mw := multipart.NewWriter(body)
	fw, _ := mw.CreateFormFile("file", fileName)
	io.Copy(fw, reader)

	contentType := mw.FormDataContentType()
	mw.Close()

	res, err := http.Post(url, contentType, body)
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}
