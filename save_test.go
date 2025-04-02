package bettergoapi

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/vorsprung/jsonapi-go"
)

func TestSaveLoad(t *testing.T) {
	var monitor Monitor = Monitor{}
	var monitorAsReloaded []Monitor = []Monitor{}
	path := "/tmp/save_test.json"
	load_filebytes, _ := os.ReadFile("testdata/example_single_monitor.json")
	// this data is the same as data from API
	jsonapi.Unmarshal(load_filebytes, &monitor)
	SaveToFile([]Monitor{monitor}, path)
	filebytes, err := os.ReadFile(path)
	assert.Nilf(t, err, "file read error is %v", err)
	json.Unmarshal(filebytes, &monitorAsReloaded)
	assert.Equalf(t, monitor.URL, monitorAsReloaded[0].URL, "url %s == %s", monitor.URL, monitorAsReloaded[0].URL)
	v := reflect.ValueOf(monitor)
	typeOfS := v.Type()
	v2 := reflect.ValueOf(monitorAsReloaded[0])
	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfS.Field(i).Name
		if fieldName == "ID" || fieldName == "Type" || v.Field(i).CanInterface() {
			continue
		}

		assert.Equal(t, v.Field(i).Interface(), v2.Field(i).Interface(), "field %s expected \"%v\" got \"%v\"",
			fieldName, v.Field(i).Interface(), v2.Field(i).Interface())
	}
}

func TestSaveBad(t *testing.T) {
	var monitor Monitor = Monitor{}
	var badpath = "/dev/xyz"
	err := SaveToFile([]Monitor{monitor}, badpath)
	operation := strings.Contains(err.Error(), "operation not permitted")
	permission := strings.Contains(err.Error(), "permission denied")
	assert.True(t, operation || permission, "file read error is %v", err)
}

func TestLoad(t *testing.T) {
	var testpath = "/tmp/foo.json"
	// clear last test
	os.Remove(testpath)
	var monitors Monitors = []Monitor{}
	// load data from file as if it was from API
	var path = "testdata/example_monitors_list.json"
	load_filebytes, _ := os.ReadFile(path)
	a, _ := jsonapi.Unmarshal(load_filebytes, &monitors)
	assert.NotNil(t, a)
	// save monitors to file
	_ = SaveToFile(monitors, testpath)
	// get back monitors from file
	m, err := LoadFromFile(testpath)
	assert.Nilf(t, err, "file read error is %v", err)
	assert.Equal(t, 6, len(m), "expected 6 monitors, got %d", len(m))
}

func TestLoadBad(t *testing.T) {
	var badpath = "/dev/xyz"
	_, err := LoadFromFile(badpath)
	assert.Containsf(t, err.Error(), "no such file or directory", "file read error is %v", err)
}

type TestUploader struct {
	Dummy *manager.UploadOutput
	err   error
}

func (u *TestUploader) Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	return u.Dummy, u.err
}

func TestSaveS3(t *testing.T) {
	var u *TestUploader = &TestUploader{}
	var path = "testdata/example_monitors_list.json"
	m, _ := LoadFromFile(path)
	ctx := context.Background()
	res := SaveToS3(ctx, m, "foo", u)
	assert.Nil(t, res)
}

type TestDownloader struct {
	Dummy int64
	err   error
}

func (u *TestDownloader) Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, opts ...func(*manager.Downloader)) (int64, error) {
	return u.Dummy, u.err
}

func TestLoadS3(t *testing.T) {
	var u *TestDownloader = &TestDownloader{Dummy: 23}
	ctx := context.Background()
	_, res := LoadFromS3(ctx, "foo", u)
	assert.Nil(t, res)
}

func TestLoadS3Bad(t *testing.T) {
	var u *TestDownloader = &TestDownloader{Dummy: 23, err: io.EOF}
	ctx := context.Background()
	_, res := LoadFromS3(ctx, "foo", u)
	assert.NotNil(t, res)
}
