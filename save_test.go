package bettergoapi

import (

	//"math"

	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pieoneers/jsonapi-go"
	"github.com/stretchr/testify/assert"
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
	assert.Containsf(t, err.Error(), "operation not permitted", "file read error is %v", err)
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
	Dummy *s3manager.UploadOutput
	err   error
}

func (u *TestUploader) Upload(*s3manager.UploadInput) (*s3manager.UploadOutput, error) {
	return u.Dummy, u.err
}

func TestSaveS3(t *testing.T) {
	var u *TestUploader = &TestUploader{}
	var path = "testdata/example_monitors_list.json"
	m, _ := LoadFromFile(path)
	res := SaveToS3(m, "foo", u)
	assert.Nil(t, res)
}

//	type MyDownloader interface {
//		Download(io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error)
//	}
type TestDownloader struct {
	Dummy int64
	err   error
}

func (u *TestDownloader) Download(io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error) {
	return u.Dummy, u.err
}

func TestLoadS3(t *testing.T) {
	var u *TestDownloader = &TestDownloader{Dummy: 23}
	_, res := LoadFromS3("foo", u)
	assert.Nil(t, res)

}

func TestLoadS3Bad(t *testing.T) {
	var u *TestDownloader = &TestDownloader{Dummy: 23, err: io.EOF}
	_, res := LoadFromS3("foo", u)
	assert.NotNil(t, res)

}
