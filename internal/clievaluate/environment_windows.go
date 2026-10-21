//go:build windows

package clievaluate

import (
	"errors"
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
)

// CurrentLimits verifies the real current Job controls before an environment runs.
func CurrentLimits(_ string, b Budget) ([]NamedValue, error) {
	if err := ValidateBudget(b); err != nil {
		return nil, err
	}
	var limits windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	if err := windows.QueryInformationJobObject(0, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&limits)), uint32(unsafe.Sizeof(limits)), nil); err != nil {
		return nil, err
	}
	var cpu struct {
		Flags uint32
		Rate  uint32
	}
	if err := windows.QueryInformationJobObject(0, windows.JobObjectCpuRateControlInformation, uintptr(unsafe.Pointer(&cpu)), uint32(unsafe.Sizeof(cpu)), nil); err != nil {
		return nil, err
	}
	flags := uint32(windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE | windows.JOB_OBJECT_LIMIT_JOB_MEMORY | windows.JOB_OBJECT_LIMIT_ACTIVE_PROCESS)
	if limits.BasicLimitInformation.LimitFlags&flags != flags || uint64(limits.JobMemoryLimit) != uint64(b.MemoryBytes) || limits.BasicLimitInformation.ActiveProcessLimit != uint32(b.WindowsProcesses) || cpu.Flags&5 != 5 || cpu.Rate != uint32(b.CPUPercent*100) {
		return nil, errors.New("windows worker does not have the requested hard CPU/memory/process limits")
	}
	priority, err := windows.GetPriorityClass(windows.CurrentProcess())
	if err != nil {
		return nil, err
	}
	if priority != windows.BELOW_NORMAL_PRIORITY_CLASS {
		return nil, errors.New("windows worker must use below-normal priority")
	}
	return []NamedValue{{Name: "mechanism", Value: "Windows Job Object"}, {Name: "cpu_percent", Value: strconv.Itoa(b.CPUPercent)}, {Name: "cpu_denominator", Value: "aggregate system CPU"}, {Name: "memory_bytes", Value: strconv.FormatInt(b.MemoryBytes, 10)}, {Name: "processes", Value: strconv.Itoa(b.WindowsProcesses)}, {Name: "priority", Value: "below_normal"}}, nil
}
