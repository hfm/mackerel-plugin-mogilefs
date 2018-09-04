package mpmogilefs

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	flag "github.com/docker/docker/pkg/mflag"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

func (m MogilefsPlugin) parseStats(conn io.Reader) (map[string]interface{}, error) {
	scanner := bufio.NewScanner(conn)
	stats := make(map[string]interface{})

	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)
		if s == "." {
			return stats, nil
		}

		res := strings.Split(s, " ")
		stats[res[0]] = res[1]
	}

	if err := scanner.Err(); err != nil {
		return stats, err
	}

	return nil, nil
}

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK             int = 0
	ExitCodeParseFlagError int = 1 + iota
	ExitCodeError
)

/* graphdef is Graph definition for mackerelplugin.

mogilefs.activity:
	pending_queries ... The number of workers queued queries
	processing_queries ... The number of workers processing requests
	bored_queryworkers ... The number of idle workers

mogilefs.queries:
	queries ... The number of queries requested from MogileFS clients

mogilefs.times_out_of_qworkers:
	times_out_of_qworkers ... The number of pending queries dispite issued the queues

mogilefs.work_queue_for:
	work_queue_for_delete ... The number of active delete workers
	work_queue_for_fsck ... The number of active fsck workers
	work_queue_for_replicate ... The number of active replicate workers

mogilefs.work_sent_to:
	work_sent_to_delete ... The number of processed delete worker
	work_sent_to_fsck ... The number of processed fsck worker
	work_sent_to_replicate ... The number of processed replicate worker
*/
var graphdef = map[string](mp.Graphs){
	"mogilefs.stats.activity": mp.Graphs{
		Label: "MogileFS tracker activity",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "pending_queries", Label: "Pending queries", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "processing_queries", Label: "Processing queries", Diff: false, Type: "uint64", Stacked: true},
			mp.Metrics{Name: "bored_queryworkers", Label: "Bored queryworkers", Diff: false, Type: "uint64", Stacked: true},
		},
	},
	"mogilefs.stats.queries": mp.Graphs{
		Label: "MogileFS tracker queries",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "queries", Label: "queries", Diff: false, Type: "uint64"},
		},
	},
	"mogilefs.stats.times_out_of_qworkers": mp.Graphs{
		Label: "MogileFS times out of querieworkers",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "times_out_of_qworkers", Label: "time out of queryworkers", Diff: false, Type: "uint64"},
		},
	},
	"mogilefs.stats.work_queue_for": mp.Graphs{
		Label: "MogileFS work_queue_for",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "work_queue_for_delete", Label: "work_queue_for_delete", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "work_queue_for_fsck", Label: "work_queue_for_fsck", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "work_queue_for_replicate", Label: "work_queue_for_replicate", Diff: false, Type: "uint64"},
		},
	},
	"mogilefs.stats.work_sent_to": mp.Graphs{
		Label: "MogileFS work_sent_to",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "work_sent_to_delete", Label: "work_sent_to_delete", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "work_sent_to_fsck", Label: "work_sent_to_fsck", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "work_sent_to_replicate", Label: "work_sent_to_replicate", Diff: false, Type: "uint64"},
		},
	},
}

// MogilefsPlugin mackerel plugin for MoglieFS.
type MogilefsPlugin struct {
	Target string
}

// FetchMetrics interface for mackerelplugin.
func (m MogilefsPlugin) FetchMetrics() (map[string]interface{}, error) {
	raddr, err := net.ResolveTCPAddr("tcp", m.Target)
	if err != nil {
		_ = fmt.Errorf("Relosve error: %v\n", err)
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		_ = fmt.Errorf("DialTCP error: %v\n", err)
		return nil, err
	}

	fmt.Fprintln(conn, "!stats")

	return m.parseStats(conn)
}

// GraphDefinition interface for mackerelplugin.
func (m MogilefsPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

// CLI is the object for command line interface.
type CLI struct {
	outStream, errStream io.Writer
}

// Run is to parse flags and Run helper (MackerelPlugin) with the given arguments.
func (c *CLI) Run(args []string) int {
	// Flags
	var (
		host     string
		port     string
		tempfile string
		version  bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)

	flags.StringVar(&host, []string{"H", "host"}, "127.0.0.1", "Host of mogilefsd")
	flags.StringVar(&port, []string{"p", "port"}, "7001", "Port of mogilefsd")
	flags.StringVar(&tempfile, []string{"t", "tempfile"}, "/tmp/mackerel-plugin-mogilefs", "Temp file name")
	flags.BoolVar(&version, []string{"v", "version"}, false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(c.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	// Create MackerelPlugin for MogileFS
	var mogilefs MogilefsPlugin
	mogilefs.Target = net.JoinHostPort(host, port)
	helper := mp.NewMackerelPlugin(mogilefs)
	helper.Tempfile = tempfile

	helper.Run()

	return ExitCodeOK
}
