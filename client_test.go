package bettergoapi

import (
	"bytes"
	"io"
	"log"
	"math"

	//"math"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

type MyHttpClient struct {
	getThis  *http.Response
	getError error
}

func (m MyHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.getThis, m.getError
}

func TestClient(t *testing.T) {
	log.SetOutput(io.Discard)
	rc := new(MyHttpClient)
	rc.getThis = &http.Response{}
	rc.getThis.StatusCode = 200
	filebytes, err := os.ReadFile("testdata/example_monitors_list.json")
	assert.Nil(t, err, "file read error is %v", err)
	rc.getThis.Body = io.NopCloser(bytes.NewReader(filebytes))
	res, err := GetMonitors(rc)
	assert.Nil(t, err, "no error expected")
	assert.IsType(t, Monitor{}, res[0], "Monitor data expected")
	assert.Len(t, res, 6, "Multiple data items expected")
}

func TestBadClient(t *testing.T) {
	// bad http code
	log.SetOutput(io.Discard)
	rc := new(MyHttpClient)
	rc.getThis = &http.Response{}
	rc.getThis.StatusCode = 451 // legal reasons
	filebytes, _ := os.ReadFile("testdata/example_monitors_list.json")
	rc.getThis.Body = io.NopCloser(bytes.NewReader(filebytes))
	res, _ := GetMonitors(rc)
	assert.Nil(t, res, "no result expected")
	// bad json data
	rc.getThis = &http.Response{}
	rc.getThis.StatusCode = 200
	filebytes, _ = os.ReadFile("testdata/example_broken_single_monitor.json")
	rc.getThis.Body = io.NopCloser(bytes.NewReader(filebytes))
	_, err := GetMonitors(rc)
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
	// error in bytes reading
	rc.getThis = &http.Response{}
	rc.getThis.StatusCode = 200
	filebytes, _ = os.ReadFile("testdata/example_monitors_list.json")
	rc.getThis.Body = io.NopCloser(bytes.NewReader(filebytes))
	rc.getError = io.ErrNoProgress
	_, err = GetMonitors(rc)
	assert.NotNil(t, err)
	// break NewRequest in GetRemote
	_, err = GetRemote("<http://ww", rc)
	assert.NotNil(t, err)
}

func TestPostClient(t *testing.T) {
	log.SetOutput(io.Discard)
	rc := new(MyHttpClient) // mock client
	nm := &Monitor{}        // input data for new monitor
	rc.getThis = &http.Response{}
	rc.getThis.StatusCode = 201
	filebytes, err := os.ReadFile("testdata/example_single_monitor.json")
	assert.Nil(t, err)
	assert.Len(t, filebytes, 1022)
	rc.getThis.Body = io.NopCloser(bytes.NewReader(filebytes))
	nm.URL = "http://www.google.com"
	res, err := PutMonitor(rc, *nm)
	assert.Nil(t, err)
	assert.Equal(t, "up", res.Status, "status expected")
}

func TestPatchClient(t *testing.T) {
	log.SetOutput(io.Discard)
	rc := new(MyHttpClient)                               // mock client
	nm := &Monitor{ID: "1083221", Paused: aws.Bool(true)} // input data for patch
	filebytes, _ := os.ReadFile("testdata/example_single_monitor.json")
	rc.getThis = &http.Response{}
	rc.getThis.StatusCode = 200
	rc.getThis.Body = io.NopCloser(bytes.NewReader(filebytes))
	res, err := PatchMonitor(rc, *nm)
	assert.Nilf(t, err, "no error expected, res is %v", err)
	assert.NotNil(t, res)
}

func TestPauseMonitor(t *testing.T) {
	rc := new(MyHttpClient)                               // mock client
	nm := &Monitor{ID: "1083221", Paused: aws.Bool(true)} // input data for patch
	filebytes, _ := os.ReadFile("testdata/example_single_monitor.json")
	rc.getThis = &http.Response{}
	rc.getThis.StatusCode = 200
	rc.getThis.Body = io.NopCloser(bytes.NewReader(filebytes))
	_, err := PatchMonitor(rc, *nm)
	assert.Nilf(t, err, "no error expected, err is %v", err)
}

func TestBadPatchClient(t *testing.T) {
	log.SetOutput(io.Discard)
	rc := new(MyHttpClient)                // mock client
	nm := &Monitor{Paused: aws.Bool(true)} // input data for patch
	rc.getThis = &http.Response{}
	rc.getThis.StatusCode = 500
	rc.getError = &http.MaxBytesError{}
	_, err := PatchMonitor(rc, *nm)
	assert.NotNil(t, err)
	// get at error in PatchRemote
	filebytes, _ := os.ReadFile("testdata/example_single_monitor.json")
	rc.getThis.Body = io.NopCloser(bytes.NewReader(filebytes))
	_, err = PatchMonitor(rc, *nm)
	assert.NotNil(t, err)
	// kill marshal
	nm.Port = math.NaN()
	_, err = PatchMonitor(rc, *nm)
	assert.NotNil(t, err)
}
func TestBadPostClient(t *testing.T) {
	log.SetOutput(io.Discard)
	// remote error
	rc := new(MyHttpClient) // mock client
	nm := &Monitor{}        // input data for new monitor
	rc.getThis = &http.Response{}
	rc.getThis.StatusCode = 500
	rc.getError = &http.MaxBytesError{}
	filebytes, err := os.ReadFile("testdata/example_single_monitor.json")
	assert.Nil(t, err)
	assert.Len(t, filebytes, 1022)
	rc.getThis.Body = io.NopCloser(bytes.NewReader(filebytes))
	_, err = PutMonitor(rc, *nm)
	assert.NotNil(t, err)
	// kill marshal
	nm.Port = math.NaN()
	rc.getError = nil
	err = nil
	_, err = PutMonitor(rc, *nm)
	assert.NotNil(t, err)
	// kill makeReq
	nm.Port = 5050
	err = nil
	req, _ := makeReq(">", "[/]foo", nil)
	assert.Nil(t, req)
	// nil result from Do() in PostRemote
	nm.URL = "http://www.google.com"
	rc.getThis = nil
	res, _ := PostRemote(rc, "http://www.google.com", nil)
	assert.Nil(t, res)
}

func TestAWSSession(t *testing.T) {
	var sess session.Session
	awsSession = &sess
	res := GetAWS()
	assert.NotNil(t, res)
}
