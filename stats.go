package main

import (
	"bufio"
	"io"
	"strings"
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
