package better

import (
	//"encoding/json"
	"os"
	"testing"

	"github.com/pieoneers/jsonapi-go"
	"github.com/stretchr/testify/assert"
)

func TestSingleMonitor(t *testing.T) {
	var monitor Monitor = Monitor{}
	filebytes, err := os.ReadFile("testdata/example_single_monitor.json")
	assert.Nilf(t, err, "file read error is %v", err)
	jsonapi.Unmarshal(filebytes, &monitor)
	assert.Equalf(t, "https://wl4-prod-app1.linuxufo.com/indigo/healthcheck", monitor.URL, "URL is >%v<", monitor.URL)
	assert.Equalf(t, monitor.GetID(), "1083221", "id %s", monitor.GetID)
	assert.Equalf(t, monitor.GetType(), "monitor", "type %s", monitor.GetType)
	assert.IsType(t, Monitor{}, monitor.GetData(), "data")
}

func TestListMonitors(t *testing.T) {
	var monitors Monitors = Monitors{}
	filebytes, err := os.ReadFile("testdata/example_monitors_list.json")
	assert.Nilf(t, err, "file read error is %v", err)
	jsonapi.Unmarshal(filebytes, &monitors)
	assert.Lenf(t, monitors, 6, "correct number of monitors found")

}
