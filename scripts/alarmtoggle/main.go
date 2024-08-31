package main

import (
	"better"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	// option 1, store alarms in a set or use the last stored alarm set
	// option 2, show all alarm states on or off.  Do not write
	// option 3, setpolicy, on , off, last, reverse
	// option 4, pattern, only set alarms that substring match this url pattern
	store := flag.Bool("store", false, "store alarms or use last stored")
	show := flag.Bool("show", false, "show alarms, do not write")
	set := flag.String("set", "last", "setpolicy, on , off, last as file, reverse of file")
	pattern := flag.String("pattern", "", "only set alarms that substring match this url pattern")
	flag.Parse()
	// validate the setpolicy flag
	if _, ok := map[string]int{"on": 1, "off": 1, "last": 1, "reverse": 1}[*set]; !ok {
		flag.Usage()
		os.Exit(1)
	}

	// default with no flags is to show the alarms
	if flag.NFlag() == 0 {
		show = new(bool)
		*show = true
	}

	// validate the setpolicy flag
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
	// dump the alarm settings, all fields.  This is for show only or for store only
	if *show || *store {
		allMonitors, err := better.GetMonitors(client)
		if err != nil {
			log.Fatal(err)
		}
		bytesAllMonitors, err := json.MarshalIndent(allMonitors, "  ", "  ")
		if err != nil {
			log.Fatal(err)
		}
		if *show {
			log.Println(string(bytesAllMonitors))
			os.Exit(0)
		}

		// store the current alarm states to a file
		if *store {
			err := better.SaveToFile(allMonitors, "/tmp/alarms.json")
			if err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}

	if *set != "" {
		storedMonitors, err := better.LoadFromFile("/tmp/alarms.json")
		if err != nil {
			log.Fatalf("reading alarms from store file %s", err)
		}
		happy := true
		// what if the alarm is not in the stored file?
		for _, monitor := range storedMonitors {
			if strings.Contains(monitor.URL, *pattern) || *pattern == "" {
				// force on or off if flag is set
				if *set == "on" || *set == "off" {
					*monitor.Paused = (*set == "off")
				}
				// reverse the state of the alarm if flag is set
				if *set == "reverse" {
					*monitor.Paused = !*monitor.Paused
				}
				// default is to set the alarm to the state it was in when the file was saved
				err := putShow(client, monitor)
				if err != nil {
					happy = false
				}
			}
		}
		if happy {
			os.Exit(0)
		}
		os.Exit(1)
	} else {
		log.Println("this shouldn't happen")
		flag.Usage()
		os.Exit(1)
	}
}

func putShow(client *http.Client, monitor better.Monitor) error {
	pauseOnly := better.Monitor{Paused: monitor.Paused, ID: monitor.ID}
	_, err := better.PatchMonitor(client, pauseOnly)
	status := "unknown"
	if err != nil {
		status = err.Error()
	} else {
		status = "OK"
	}
	log.Printf("%s paused=%v %s alarm id %s\n", status, *monitor.Paused, monitor.URL, monitor.ID)
	return err
}
