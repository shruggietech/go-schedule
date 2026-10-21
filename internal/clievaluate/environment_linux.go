//go:build linux

package clievaluate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CurrentLimits verifies the owned cgroup's real kernel controls before execution.
func CurrentLimits(runID string, b Budget) ([]NamedValue, error) {
	if err := ValidateBudget(b); err != nil {
		return nil, err
	}
	if !validRunID(runID) || os.Geteuid() == 0 {
		return nil, errors.New("linux QA worker requires an existing nonroot user and owned run identity")
	}
	data, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		return nil, err
	}
	group := ""
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "0::") {
			group = strings.TrimPrefix(line, "0::")
		}
	}
	if filepath.Base(group) != "gosched-qa-"+runID+".service" || os.Getenv("INVOCATION_ID") == "" {
		return nil, errors.New("linux worker is not inside its owned transient user cgroup")
	}
	root := filepath.Join("/sys/fs/cgroup", strings.TrimPrefix(group, "/"))
	read := func(name string) (string, error) {
		value, err := os.ReadFile(filepath.Join(root, name))
		return strings.TrimSpace(string(value)), err
	}
	memory, err := read("memory.max")
	if err != nil {
		return nil, err
	}
	tasks, err := read("pids.max")
	if err != nil {
		return nil, err
	}
	swap, err := read("memory.swap.max")
	if err != nil {
		return nil, err
	}
	cpu, err := read("cpu.max")
	if err != nil {
		return nil, err
	}
	parts := strings.Fields(cpu)
	if len(parts) != 2 {
		return nil, errors.New("invalid cgroup CPU control")
	}
	quota, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	period, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	if period <= 0 || quota*100 != int64(b.CPUPercent)*period || memory != strconv.FormatInt(b.MemoryBytes, 10) || tasks != strconv.Itoa(b.LinuxTasks) || swap != "0" {
		return nil, fmt.Errorf("linux worker controls do not match budget: cpu=%s memory=%s tasks=%s swap=%s", cpu, memory, tasks, swap)
	}
	stat, err := os.ReadFile("/proc/self/stat")
	if err != nil {
		return nil, err
	}
	end := strings.LastIndexByte(string(stat), ')')
	if end < 0 {
		return nil, errors.New("cannot read worker priority")
	}
	fields := strings.Fields(string(stat)[end+1:])
	if len(fields) < 17 || fields[16] != "10" {
		return nil, errors.New("linux worker must have the requested nice value 10")
	}
	return []NamedValue{{Name: "mechanism", Value: "cgroup-v2 transient user service"}, {Name: "cgroup", Value: group}, {Name: "invocation_id", Value: os.Getenv("INVOCATION_ID")}, {Name: "cpu_percent", Value: strconv.Itoa(b.CPUPercent)}, {Name: "cpu_denominator", Value: "one logical CPU"}, {Name: "memory_bytes", Value: memory}, {Name: "tasks_including_threads", Value: tasks}, {Name: "swap_bytes", Value: swap}, {Name: "nice", Value: fields[16]}}, nil
}
