package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	flag "github.com/tcnksm/mflag"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// graphdef is Graph definition for mackerelplugin.
var graphdef = map[string](mp.Graphs){
	"mogilefsd.queries": mp.Graphs{
		Label: "MogileFS tracker queries",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "queries", Label: "queries", Diff: false, Type: "uint64"},
		},
	},
	"mogilefsd.stats": mp.Graphs{
		Label: "MoglieFS tracker activity",
		Unit:  "integer",
		Stacked: true,
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "pending_queries", Label: "Pending queries", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "processing_queries", Label: "Processing queries", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "bored_queryworkers", Label: "Bored queryworkers", Diff: false, Type: "uint64"},
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

// GraphDefinition interface for mackerelplugin.
func (m MogilefsPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

// Parse flags and Run helper (MackerelPlugin) with the given arguments.
func main() {
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
	if err := flags.Parse(os.Args[1:]); err != nil {
		os.Exit(ExitCodeError)
	}

	// Show version
	if version {
		fmt.Fprintf(os.Stderr, "%s version %s\n", Name, Version)
		os.Exit(ExitCodeOK)
	}

	// Create MackerelPlugin for MogileFS
	var mogilefs MogilefsPlugin
	mogilefs.Target = net.JoinHostPort(host, port)
	helper := mp.NewMackerelPlugin(mogilefs)
	helper.Tempfile = tempfile

	helper.Run()
}
