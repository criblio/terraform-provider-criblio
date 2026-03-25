package shared

import (
	"encoding/json"
	"testing"
)

// Regression: real Search dashboards use object-shaped xAxis/yAxis in element config; string-only schema broke JSON decode.
func TestDashboardElementUnionCounterSingleAPIPayload(t *testing.T) {
	data := []byte(`{"config":{"onClickAction":{"type":"Run a new search","search":"dataset=\"Availability_daily_zeus\"\n| timestats availability=round(100*sum(available)/sum(sum),4)"},"style":false,"applyThreshold":true,"colorThresholds":{"thresholds":[{"color":"#45850B","threshold":99.9},{"color":"#EFDB23","threshold":99},{"color":"#B20000","threshold":95}]},"separator":true,"legend":{"position":"Right","truncate":true},"colorPalette":0,"colorPaletteReversed":false,"customData":{"trellis":false,"connectNulls":"Leave gaps","stack":false,"dataFields":["availability"],"seriesCount":1},"xAxis":{"labelOrientation":0,"position":"Bottom"},"yAxis":{"position":"Left","scale":"Linear","splitLine":true},"suffix":"%"},"search":{"type":"inline","query":"dataset=\"availability_daily\" service=\"zeus\"\n| summarize availability=round(100*sum(available)/sum(sum),4)","earliest":"-7d","latest":"now","timezone":"$timerange.timezone$"},"id":"ywatz6fqq","type":"counter.single","layout":{"x":0,"y":0,"w":6,"h":2},"title":"Cloud API"}`)

	var vis DashboardElementVisualization
	if err := json.Unmarshal(data, &vis); err != nil {
		t.Fatalf("visualization: %v", err)
	}
	var u DashboardElementUnion
	if err := json.Unmarshal(data, &u); err != nil {
		t.Fatalf("union: %v", err)
	}
}
