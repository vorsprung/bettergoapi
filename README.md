bettergoapi Uptime API Clients
-------------------------
This is a go library for setting up and pausing and unpausing bettergoapiUptime Monitors

This can be used to temporarily pause monitoring during patching

It uses the Uptime API see [Getting started with Uptime API](https://bettergoapistack.com/docs/uptime/api/getting-started-with-uptime-api/)

Example programs are included to set up from a CSV file and to pause and unpause monitors

Programs using the library
--------------------------
| path                            |                            |
| --------------------------------| -------------------------- |
| scripts/alarmtoggle/main.go     | alarmtoggle                |
| scripts/fromCSV/main.go         | fromCSV                    |


How to build alarm toggle for Linux
-----------------------------------
    go mod tidy
    env GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o alarmtoggle scripts/alarmtoggle/main.go

This makes a single file static binary that can be copied directly to Intel / AMD hosts and used
Does not bundle the TEAM_TOKEN (see below).  This must be in the environment where alarmtoggle is used

Running alarm toggle
--------------------
#### To test script is working

    ./alarmtoggle 

Shows verbose dump of json of monitors

#### To store all alarm states locally

    ./alarmtoggle  -store=true

alarm states are stored in /tmp/alarms.json

#### To turn off monitors (pause monitors)

   ./alarmtoggle  -set=off -pattern=testrpt1

output is lines in this format

    2023/04/14 11:53:25 OK paused=true https://romeo.linuxufo.com/state.txt alarm id 1148817

pattern match is on the URL only.  It is a substring match.  alarmtoggle -set= only works if -store=true was previously run.  It uses /tmp/alarms.json as a list of alarms

#### To turn on monitors (unpause monitors)

   ./alarmtoggle  -set=on  -pattern=testrpt1

output is lines in this format

    2023/04/14 11:57:05 OK paused=false https://romeo.linuxufo.com/state.txt alarm id 1148817

pattern match is on the URL only.  It is a substring match.  alarmtoggle -set= only works if -store=true was previously run.  It uses /tmp/alarms.json as a list of alarms

#### To run without building
To run a local copy without building

    go run scripts/alarmtoggle/main.go  -set=on  -pattern=testrpt1

#### Problem "patch failed 401 Unauthorized"
If this happens

    2023/04/14 11:59:39 env var not set and failed to get from ssm
    2023/04/14 11:59:39 patch failed 401 Unauthorized
    2023/04/14 11:59:39 response was {"errors":"Invalid Team API token. How to find your Team API token: https://bettergoapistack.com/docs/uptime/api/getting-started-with-bettergoapi-uptime-api#obtaining-a-bettergoapi-uptime-api-token"}
    2023/04/14 11:59:39 patch failed paused=false https://romeo.linuxufo.com/state.txt alarm id 1148817
    exit status 1

The TEAM_TOKEN is not set, see below

#### alarmtoggle all options

| option   | parameter  | default | use or result                                                         |
| -------- | ---------- | ------- | --------------------------------------------------------------------- |
| -show    | true/false | false   | when true verbose show as json all monitors.  Do not do anything else |
| -store   | true/false | false   | write json of monitors to /tmp/alarms.json                            |
| -set     | on         |         | on monitors force on                                                  |
| -set     | off        |         | off monitors force off                                                |
| -set     | last       |         | use previously stored values for pause state of monitor               |
| -set     | reverse    |         | use opposite of stored values for pause state of monitor              |
| -pattern | string     |         | only set alarms that substring match this url pattern                 |


Dependencies
------------
 * To build, go lang compiler (v 1.20.0 used during development)
 * Access to bettergoapi Uptime requires API Token, see https://bettergoapistack.com/docs/uptime/api/getting-started-with-bettergoapi-uptime-api/#obtaining-a-bettergoapi-uptime-api-token
 * The API Token should be set as a env variable called TEAM_TOKEN
 * There is code to use a SSM parameter called "bettergoapi-monitor-token" to store/retrieve the token but this is not tested
 * The host running the alarmtoggle program needs to have network access to the uptime.betterstack.com api endpoints

Files
-----
| file                |                                                           |
| ------------------- | --------------------------------------------------------- |
| client.go           | golang http client                                        |
| client_test.go      |                                                           |
| get_monitor_test.go |                                                           |
| go.mod              | go module config file                                     |
| go.sum              | go module config file                                     |
| loadcompare.go      | golang compare monitors (not currently used)              |
| loadcompare_test.go |                                                           |
| monitor.go          | golang definition of bettergoapi uptime monitor data structure |
| save.go             | save monitors as json to file or s3                       |
| save_test.go        |                                                           |


Platforms
---------
Tested on Ubuntu 22.04 and MacOS 14.6
go version go 1.22.6
