package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/medyagh/goprettyorgohome/models"
)

// parseJSON is a very forgiving JSON parser.
func ParseJSON(path string) ([]models.TestEvent, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	events := []models.TestEvent{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Go's -json output is line-by-line JSON events
		b := scanner.Bytes()
		if b[0] == '{' {
			ev := models.TestEvent{}
			err = json.Unmarshal(b, &ev)
			if err != nil {
				continue
			}
			events = append(events, ev)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, err
}

// group events by their test name
func ProcessEvents(evs []models.TestEvent) []models.TestGroup {
	gm := map[string]int{}
	groups := []models.TestGroup{}
	for _, e := range evs {
		if e.Test == "" {
			continue
		}
		index, ok := gm[e.Test]
		if !ok {
			index = len(groups)
			groups = append(groups, models.TestGroup{
				TestName: e.Test,
				Start:    e.Time,
			})
			gm[e.Test] = index
		}
		groups[index].Events = append(groups[index].Events, e)
		groups[index].Status = e.Action
	}

	// Hide ancestors
	for k, v := range gm {
		for k2 := range gm {
			if strings.HasPrefix(k2, fmt.Sprintf("%s/", k)) {
				groups[v].Hidden = true
			}
		}
	}

	return groups
}
