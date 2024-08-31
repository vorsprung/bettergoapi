package main

import (
	"better"
	"encoding/csv"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

func main() {
	dry := flag.Bool("dry", false, "dry run")

	flag.Parse()
	filename := os.Args[1]
	//filename := "../../ht.csv"
	log.Printf("reading from file %s", filename)
	f, _ := os.Open(filename)
	r := csv.NewReader(f)
	r.Read() //skip header
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal("ReadAll ", err)
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}}
	var createdMonitors [][]string
	for _, monitorData := range records {
		var monitor better.Monitor

		monitor.URL = monitorData[0]
		monitor.MonitorType = monitorData[2]
		monitor.PronounceableName = monitorData[1]
		monitor.CheckFrequency, _ = strconv.Atoi(monitorData[4])
		monitor.RequestTimeout, _ = strconv.Atoi(monitorData[6])
		monitor.RequiredKeyword = monitorData[3]
		monitor.Call = false
		monitor.Sms = false
		monitor.Email = true
		monitor.Paused = aws.Bool(true)
		monitor.Push = true

		//monitor.SetType("monitor")

		var res better.Monitor

		if *dry {
			res.ID = "1234"
			err = nil
		} else {
			res, err = better.PutMonitor(client, monitor)
		}
		createdMonitors = append(createdMonitors, []string{monitor.URL, monitor.PronounceableName, res.ID})
		if err != nil || res.ID == "" {
			log.Printf("problem with %s - %v", monitor.URL, err)
		} else {
			log.Printf("made %s for %s ", res.ID, monitor.URL)
		}

	}
	csvOutFile := "out" + filename + ".csv"
	outFileHandle, err := os.Create(csvOutFile)
	if err != nil {
		log.Printf("problem with csv write %s %v", csvOutFile, err)
	}
	w := csv.NewWriter(outFileHandle)
	w.WriteAll(createdMonitors) // calls Flush internally
}
