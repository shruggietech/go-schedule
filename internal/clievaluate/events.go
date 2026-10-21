package clievaluate

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"slices"
	"time"
)

var eventTypes = []string{"environment_ready", "preflight_failed", "case_started", "case_finished", "case_blocked", "cleanup_failed", "capture_finished"}

func validateEvents(data []byte, r Report) error {
	if len(data) > 0 && data[len(data)-1] != '\n' {
		return errors.New("event journal has a truncated final record")
	}
	envs := map[string]bool{}
	for _, env := range r.Environments {
		envs[env.ID] = true
	}
	cases := map[string]CaseResult{}
	for _, c := range r.Cases {
		cases[caseKey(c.EnvironmentID, c.ID)] = c
	}
	sequences := map[string]int64{}
	last := map[string]time.Time{}
	started := map[string]bool{}
	finished := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 4096), 1<<20)
	for scanner.Scan() {
		line := scanner.Bytes()
		var e Event
		if err := Decode(bytes.NewReader(line), &e); err != nil {
			return err
		}
		if err := ValidateSchema("event", line); err != nil {
			return err
		}
		if e.RunID != r.RunID || e.EnvironmentID == nil || !envs[*e.EnvironmentID] || !slices.Contains(eventTypes, e.Type) {
			return errors.New("event journal identity or type mismatch")
		}
		env := *e.EnvironmentID
		sequences[env]++
		if e.Sequence != sequences[env] {
			return errors.New("event sequence is missing, duplicated or reordered")
		}
		at, err := utc(e.At)
		if err != nil || at.Before(last[env]) {
			return errors.New("event timestamp ordering is invalid")
		}
		last[env] = at
		if e.CaseID != nil {
			key := caseKey(env, *e.CaseID)
			if _, ok := cases[key]; !ok {
				return errors.New("event references an unplanned case")
			}
			switch e.Type {
			case "case_started":
				if started[key] {
					return errors.New("duplicate case start event")
				}
				started[key] = true
			case "case_finished", "case_blocked":
				if !started[key] || finished[key] {
					return errors.New("case completion event lacks one matching start")
				}
				finished[key] = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	for key, c := range cases {
		if c.Observed.ExitCode != nil && (!started[key] || !finished[key]) {
			return fmt.Errorf("executed case has no complete journal transition: %s", key)
		}
	}
	return nil
}
