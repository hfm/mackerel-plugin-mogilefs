package mpmogilefs

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

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

func TestGraphDefinition(t *testing.T) {
	var mogilefs MogilefsPlugin

	graphdef := mogilefs.GraphDefinition()
	if len(graphdef) != 5 {
		t.Errorf("GetTempfilename: %d should be 5", len(graphdef))
	}
}

func TestCLI_Run(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("mackerel-plugin-mogilefs -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("mackerel-plugin-mogilefs version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("Output=%q, want %q", errStream.String(), expected)
	}

}
