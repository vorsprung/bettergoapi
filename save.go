package bettergoapi

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type MyUploader interface {
	Upload(*s3manager.UploadInput) (*s3manager.UploadOutput, error)
}

type MyDownloader interface {
	Download(io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error)
}

// call as SaveToS3(monitor, "monitors.json", s3manager.NewUploader(getAWSsession()))
func SaveToS3(monitor []Monitor, path string, uploader MyUploader) error {
	data, _ := json.MarshalIndent(monitor, "  ", " ")
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("bettergoapi-monitor"),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})
	return err
}

// call as LoadFromS3("monitors.json", s3manager.NewDownloader(getAWSsession()))
func LoadFromS3(path string, downloader MyDownloader) ([]Monitor, error) {
	var monitor []Monitor
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(buf, &s3.GetObjectInput{
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
