package mpmogilefs

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
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
func (m *MogileFSPlugin) mogilefsGraphDef() map[string]mp.Graphs {
	labelPrefix := strings.Title(strings.Replace(m.MetricKeyPrefix(), "mogilefs", "MogileFS", -1))
	return map[string]mp.Graphs{
		"mogilefs.stats.activity": {
			Label: labelPrefix + " tracker activity",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "pending_queries", Label: "Pending queries"},
				{Name: "processing_queries", Label: "Processing queries", Stacked: true},
				{Name: "bored_queryworkers", Label: "Bored queryworkers", Stacked: true},
			},
		},
		"mogilefs.stats.queries": {
			Label: labelPrefix + " tracker queries",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "queries", Label: "queries"},
			},
		},
		"mogilefs.stats.times_out_of_qworkers": {
			Label: labelPrefix + " times out of querieworkers",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "times_out_of_qworkers", Label: "time out of queryworkers"},
			},
		},
		"mogilefs.stats.work_queue_for": {
			Label: labelPrefix + " work_queue_for",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "work_queue_for_delete", Label: "work_queue_for_delete"},
				{Name: "work_queue_for_fsck", Label: "work_queue_for_fsck"},
				{Name: "work_queue_for_replicate", Label: "work_queue_for_replicate"},
			},
		},
		"mogilefs.stats.work_sent_to": {
			Label: labelPrefix + " work_sent_to",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "work_sent_to_delete", Label: "work_sent_to_delete"},
				{Name: "work_sent_to_fsck", Label: "work_sent_to_fsck"},
				{Name: "work_sent_to_replicate", Label: "work_sent_to_replicate"},
			},
		},
	}
}

// MogileFSPlugin mackerel plugin for MoglieFS.
type MogileFSPlugin struct {
	Target   string
	Tempfile string
	prefix   string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (m *MogileFSPlugin) MetricKeyPrefix() string {
	if m.prefix == "" {
		m.prefix = "mogilefs"
	}
	return m.prefix
}

func (m *MogileFSPlugin) parseStats(conn io.Reader) (map[string]float64, error) {
	scanner := bufio.NewScanner(conn)
	stats := make(map[string]float64)
	var err error

	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)
		if s == "." {
			return stats, nil
		}

		res := strings.Split(s, " ")
		stats[res[0]], err = strconv.ParseFloat(res[1], 64)
		if err != nil {
			return nil, err
		}
	}

	if err = scanner.Err(); err != nil {
		return stats, err
	}
	return nil, nil
}

// FetchMetrics interface for mackerelplugin.
func (m *MogileFSPlugin) FetchMetrics() (map[string]float64, error) {
	conn, err := net.Dial("tcp", m.Target)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(conn, "!stats")
	stats, err := m.parseStats(conn)

	return stats, err
}

// GraphDefinition interface for mackerelplugin.
func (m *MogileFSPlugin) GraphDefinition() map[string]mp.Graphs {
	return m.mogilefsGraphDef()
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "127.0.0.1", "Hostname")
	optPort := flag.String("port", "7001", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optMetricKeyPrefix := flag.String("metric-key-prefix", "proxysql", "metric key prefix")

	flag.Parse()

	var mogilefs MogileFSPlugin

	mogilefs.Target = net.JoinHostPort(*optHost, *optPort)
	mogilefs.prefix = *optMetricKeyPrefix
	helper := mp.NewMackerelPlugin(&mogilefs)
	helper.Tempfile = *optTempfile

	helper.Run()
}
