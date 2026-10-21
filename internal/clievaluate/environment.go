package clievaluate

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// ValidateBudget rejects unbounded or host-disruptive worker settings.
func ValidateBudget(b Budget) error {
	if b.CPUPercent < 1 || b.CPUPercent > 50 || b.MemoryBytes < 64<<20 || b.MemoryBytes > 4<<30 || b.WindowsProcesses < 2 || b.WindowsProcesses > 64 || b.LinuxTasks < 32 || b.LinuxTasks > 512 || b.CaseTimeoutMS < 100 || b.CaseTimeoutMS > 120000 || b.TotalTimeoutMS < b.CaseTimeoutMS || b.TotalTimeoutMS > 7200000 || b.ChannelBytes < 1024 || b.ChannelBytes > 1<<20 || b.ArtifactBytes < 1<<20 || b.ArtifactBytes > 256<<20 {
		return errors.New("QA budgets must use bounded positive CPU, memory, processes, deadlines and output limits; start from the documented defaults")
	}
	return nil
}

// WSLWorkerArguments applies Linux controls independently of the Windows launcher.
// All values are separate argv elements; no shell interpolation occurs.
func WSLWorkerArguments(runID, root, executable string, args []string, b Budget) ([]string, error) {
	if err := ValidateBudget(b); err != nil {
		return nil, err
	}
	if !validRunID(runID) {
		return nil, errors.New("invalid owned run identity")
	}
	result := []string{"systemd-run", "--user", "--no-ask-password", "--quiet", "--wait", "--pipe", "--collect", "--service-type=exec", "--expand-environment=no", "--unit=gosched-qa-" + runID + ".service", "--property=CPUQuota=" + strconv.Itoa(b.CPUPercent) + "%", "--property=MemoryMax=" + strconv.FormatInt(b.MemoryBytes, 10), "--property=MemorySwapMax=0", "--property=TasksMax=" + strconv.Itoa(b.LinuxTasks), "--property=Nice=10", "--property=RuntimeMaxSec=" + fmt.Sprintf("%ds", int64((time.Duration(b.TotalTimeoutMS)*time.Millisecond+30*time.Second)/time.Second)), "--property=TimeoutStopSec=5s", "--property=KillMode=control-group", "--property=SendSIGKILL=yes", "--property=Restart=no", "--property=LimitCORE=0", "--working-directory=" + root, "--setenv=GOMAXPROCS=2", "--setenv=GOFLAGS=-p=2", "--setenv=GOMEMLIMIT=512MiB", executable}
	return append(result, args...), nil
}

func validRunID(id string) bool {
	if len(id) != 36 {
		return false
	}
	for i, r := range id {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if r != '-' {
				return false
			}
		} else if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return false
		}
	}
	return true
}
