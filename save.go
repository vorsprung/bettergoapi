package bettergoapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MyUploader interface {
	Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error)
}

type MyDownloader interface {
	Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, opts ...func(*manager.Downloader)) (n int64, err error)
}

// call as SaveToS3(ctx, monitor, "monitors.json", manager.NewUploader(client))
func SaveToS3(ctx context.Context, monitor []Monitor, path string, uploader MyUploader) error {
	data, _ := json.MarshalIndent(monitor, "  ", " ")
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String("bettergoapi-monitor"),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})
	return err
}

// call as LoadFromS3(ctx, "monitors.json", manager.NewDownloader(client))
func LoadFromS3(ctx context.Context, path string, downloader MyDownloader) ([]Monitor, error) {
	var monitor []Monitor
	buf := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String("bettergoapi-monitor"),
		Key:    aws.String(path),
	})
	if err != nil {
		return monitor, err
	}
	json.Unmarshal(buf.Bytes(), &monitor)
	return monitor, nil
}

func SaveToFile(monitor []Monitor, path string) error {

	var file io.WriteCloser

	var err error
	file, err = os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, _ := json.MarshalIndent(monitor, "  ", " ")

	file.Write(data)

	return nil
}

func LoadFromFile(path string) ([]Monitor, error) {
	var monitor []Monitor
	filebytes, err := os.ReadFile(path)
	if err != nil {
		return monitor, err
	}
	json.Unmarshal(filebytes, &monitor)
	return monitor, nil
}
