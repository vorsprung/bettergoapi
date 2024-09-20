package bettergoapi

import (

	//"math"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSame(t *testing.T) {
	m1 := Monitor{PronounceableName: "Froggy"}
	m2 := Monitor{PronounceableName: "Froggy"}
	diff := EzDiff(m1, m2)
	eq := `{
    "pronounceable_name": "Froggy"
  }`
	assert.Equal(t, eq, diff)
}

func TestDifferent(t *testing.T) {
	m1 := Monitor{PronounceableName: "Froggy"}
	m2 := Monitor{PronounceableName: "Doggy"}
	eq := "{\n    \"pronounceable_name\": \"Doggy\"\n  }"
	diff := EzDiff(m1, m2)
	assert.Equal(t, eq, diff)
}

func TestFromCSVRecord(t *testing.T) {

	item := []string{"http://www.google.com", "Google", "HTTP", "Google", "1", "200", "5", "us", "eu", "au", "ap"}
	m := fromCSVRecord(item)
	assert.Equal(t, "Google", m.PronounceableName, "PronounceableName found")
	assert.Equal(t, "HTTP", m.MonitorType, "MonitorType found")
	assert.Equal(t, "Google", m.RequiredKeyword, "RequiredKeyword found")
	assert.Equal(t, 1, m.CheckFrequency, "CheckFrequency found")
	assert.Equal(t, 5, m.RequestTimeout, "RequestTimeout found")
	assert.Equal(t, 4, len(m.Regions), "Regions found")
	item[4] = "a" // check frequency not an integer
	m = fromCSVRecord(item)
	assert.Equal(t, m, Monitor{}, "check frequency not an integer")
	item[4] = "1" // check frequency is an integer
	item[6] = "a" // request timeout not an integer
	m = fromCSVRecord(item)
	assert.Equal(t, m, Monitor{}, "request timeout not an integer")
}

func TestSelectiveCompare(t *testing.T) {
	m1 := Monitor{PronounceableName: "Froggy"}
	m2 := Monitor{PronounceableName: "Froggy"}
	assert.True(t, SelectiveCompare(m1, m2), "Exact match")
	m2.PronounceableName = "Doggy"
	assert.False(t, SelectiveCompare(m1, m2), "PronounceableName does not match")
	m2.PronounceableName = "Froggy"
	m1.Regions = []string{"us"}
	m2.Regions = []string{"eu"}
	assert.False(t, SelectiveCompare(m1, m2), "different regions")
}
