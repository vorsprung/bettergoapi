package better

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func fromCSVRecord(item []string) Monitor {
	var regions_found []string
	for ix, region := range []string{"us", "eu", "au", "ap"} {
		if item[ix+7] != "" {
			regions_found = append(regions_found, region)
		}
	}
	check_freq, errf := strconv.Atoi(item[4])
	request_timeout, errt := strconv.Atoi(item[6])

	if errf != nil {
		log.Printf("field for checkfrequency will not convert to integer %s", item[4])
		return Monitor{}
	}
	if errt != nil {
		log.Printf("field for request_timeout will not convert to integer %s", item[6])
		return Monitor{}
	}
	return Monitor{
		URL:               item[0],
		PronounceableName: item[1],
		MonitorType:       item[2],
		RequiredKeyword:   item[3],
		CheckFrequency:    check_freq,
		//ExpectedStatusCodes: item[5],
		RequestTimeout: request_timeout,
		Regions:        regions_found,
	}
}

func SelectiveCompare(m1 Monitor, m2 Monitor) bool {
	if m1.URL == m2.URL &&
		m1.PronounceableName == m2.PronounceableName &&
		m1.MonitorType == m2.MonitorType &&
		m1.RequiredKeyword == m2.RequiredKeyword &&
		m1.CheckFrequency == m2.CheckFrequency &&
		m1.RequestTimeout == m2.RequestTimeout &&
		len(m1.Regions) == len(m2.Regions) {
		for ix, region := range m1.Regions {
			if region != m2.Regions[ix] {
				return false
			}

		}
		return true
	}
	return false
}

func EzDiff(m1 Monitor, m2 Monitor) string {
	dmp := diffmatchpatch.New()
	json1, _ := json.MarshalIndent(&m1, "  ", "  ")
	json2, _ := json.MarshalIndent(&m2, "  ", "  ")
	diffs := dmp.DiffMain(string(json1), string(json2), false)
	return dmp.DiffText2(diffs)
}
