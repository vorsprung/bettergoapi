package bettergoapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pieoneers/jsonapi-go"
)

var awsSession *session.Session

type HClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func GetAWS() *session.Session {
	if awsSession == nil {
		awsSession = session.Must(session.NewSession())
	}
	return awsSession
}

func lowerMakeReq(method string, path string, body io.Reader, Token string) (*http.Request, error) {
	req, err := http.NewRequest(method, bettergoapiURL+path, body)
	if req != nil {
		req.Header.Set("Authorization", "Bearer "+Token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("User-Agent", "curl/7.85.0")
	}
	return req, err
}

func makeReq(method string, path string, body io.Reader) (*http.Request, error) {
	Token := os.Getenv("TEAM_TOKEN")
	if Token == "" {
		sess := GetAWS()
		svc := ssm.New(sess)
		param, err := svc.GetParameter(&ssm.GetParameterInput{
			Name:           aws.String("bettergoapi-monitor-token"),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			log.Print("env var not set and failed to get from ssm")
			//log.Fatal(err)
		} else {
			Token = *param.Parameter.Value
		}
	}
	return lowerMakeReq(method, path, body, Token)
}
func GetRemote(path string, cli HClient) ([]byte, error) {

	req, _ := makeReq("GET", path, nil)

	res, err := cli.Do(req)

	if err != nil || res.StatusCode != http.StatusOK {
		return nil, err
	}

	byteData, _ := io.ReadAll(res.Body)
	res.Body.Close()
	return byteData, err

}

func PostRemote(cli HClient, path string, obj []byte) ([]byte, error) {
	req, _ := makeReq("POST", path, bytes.NewBuffer(obj))

	// uncomment this to see the raw http request.  Note this breaks the Do() as the body is consumed
	//req.Write(os.Stdout)
	var res *http.Response

	res, err := cli.Do(req)

	if res != nil && res.Body != nil {
		byteData, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode != http.StatusCreated || err != nil {
			log.Print("post failed ", res.Status)
			log.Print("response was ", string(byteData))
			return nil, errors.New("create failed")
		}
		return byteData, err
	}
	log.Print("nil res")
	return nil, err
}

func PatchRemote(cli HClient, path string, obj []byte) ([]byte, error) {
	req, _ := makeReq("PATCH", path, bytes.NewBuffer(obj))

	// uncomment this to see the raw http request.  Note this breaks the Do() as the body is consumed
	//req.Write(os.Stdout)
	var res *http.Response

	res, err := cli.Do(req)

	if res != nil && res.Body != nil {
		byteData, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode != http.StatusOK || err != nil {
			log.Print("patch failed ", res.Status)
			log.Print("response was ", string(byteData))
			return nil, errors.New("patch failed")
		}
		return byteData, err
	}
	//log.Print("nil res")
	return nil, err
}

// &http.Client{}
func GetMonitors(cli HClient) (Monitors, error) {
	d, err := GetRemote("monitors", cli)
	if err != nil {
		return nil, err
	}
	var monitors Monitors = Monitors{}
	_, err = jsonapi.Unmarshal(d, &monitors)
	if err != nil {
		return nil, err
	}
	return monitors, err
}

func PutMonitor(cli HClient, newMonitor Monitor) (Monitor, error) {
	var newMonitorResult = Monitor{}
	data, err := json.MarshalIndent(newMonitor, "  ", " ")
	if err != nil {
		return newMonitorResult, err
	}
	byteData, err := PostRemote(cli, "monitors", data)

	if err != nil {
		return newMonitorResult, err
	}

	jsonapi.Unmarshal(byteData, &newMonitorResult)
	return newMonitorResult, nil
}

func PatchMonitor(cli HClient, updateMonitor Monitor) (Monitor, error) {
	var patchResult = Monitor{}
	// save the ID to put in the path
	id := updateMonitor.ID
	// clear the updateMonitor.ID so it doesn't get sent to the API in the body
	updateMonitor.ID = ""
	data, err := json.MarshalIndent(updateMonitor, "  ", " ")
	if err != nil {
		return patchResult, err
	}
	// bytedata should contain the updated monitor
	byteData, err := PatchRemote(cli, "monitors/"+id, data)

	// if there is an error, return it and a blank monitor
	if err != nil {
		return patchResult, err
	}

	// if there is no error, unmarshal the result and return it
	jsonapi.Unmarshal(byteData, &patchResult)
	return patchResult, nil
}
