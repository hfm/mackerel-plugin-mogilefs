package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGraphDefinition(t *testing.T) {
	var mogilefs MogilefsPlugin

	graphdef := mogilefs.GraphDefinition()
	if len(graphdef) != 5 {
		t.Errorf("GetTempfilename: %d should be 5", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var mogilefs MogilefsPlugin
	stub := `uptime 35235
pending_queries 5
processing_queries 38
bored_queryworkers 3
queries 117
work_queue_for_delete 0
work_sent_to_delete 6
.
`

	mogilefsStats := bytes.NewBufferString(stub)

	stats, err := mogilefs.parseStats(mogilefsStats)
	fmt.Println(stats)
	if err != nil {
		t.Errorf("%v", err)
	}
}
