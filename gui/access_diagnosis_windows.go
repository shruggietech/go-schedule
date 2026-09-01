//go:build windows

package gui

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/shruggietech/go-schedule/internal/service"
)

const maxPreferredLength = 0xffffffff

var (
	netapi32                    = windows.NewLazySystemDLL("netapi32.dll")
	procNetLocalGroupGetMembers = netapi32.NewProc("NetLocalGroupGetMembers")
	procNetApiBufferFree        = netapi32.NewProc("NetApiBufferFree")
)

type localGroupMemberInfo0 struct {
	SID *windows.SID
}

func diagnoseAccess() accessDiagnosis {
	diagnosis := accessDiagnosis{}
	status, err := service.Control("status", "", nil)
	if err != nil {
		diagnosis.Service = "unknown"
		diagnosis.Detail = "service status unavailable: " + err.Error()
	} else {
		diagnosis.Service = status
	}

	groupSID, _, _, err := windows.LookupSID("", "goschedadmin")
	if err != nil {
		if errors.Is(err, windows.ERROR_NONE_MAPPED) {
			diagnosis.GroupExists = no
		} else {
			diagnosis.Detail = appendDiagnosis(diagnosis.Detail, "group lookup unavailable: "+err.Error())
		}
		return diagnosis
	}
	diagnosis.GroupExists = yes

	token := windows.GetCurrentProcessToken()
	user, err := token.GetTokenUser()
	if err != nil {
		diagnosis.Detail = appendDiagnosis(diagnosis.Detail, "account SID unavailable: "+err.Error())
		return diagnosis
	}
	member, err := localGroupContainsSID("goschedadmin", user.User.Sid)
	if err != nil {
		diagnosis.Detail = appendDiagnosis(diagnosis.Detail, "account membership unavailable: "+err.Error())
	} else if member {
		diagnosis.AccountMember = yes
	} else {
		diagnosis.AccountMember = no
	}
	tokenMember, err := token.IsMember(groupSID)
	if err != nil {
		diagnosis.Detail = appendDiagnosis(diagnosis.Detail, "token membership unavailable: "+err.Error())
	} else if tokenMember {
		diagnosis.TokenMember = yes
	} else {
		diagnosis.TokenMember = no
	}
	return diagnosis
}

func localGroupContainsSID(group string, want *windows.SID) (bool, error) {
	groupPtr, err := windows.UTF16PtrFromString(group)
	if err != nil {
		return false, err
	}
	var resume uintptr
	for {
		var buffer *localGroupMemberInfo0
		var read, total uint32
		status, _, _ := procNetLocalGroupGetMembers.Call(
			0,
			uintptr(unsafe.Pointer(groupPtr)),
			0,
			uintptr(unsafe.Pointer(&buffer)),
			maxPreferredLength,
			uintptr(unsafe.Pointer(&read)),
			uintptr(unsafe.Pointer(&total)),
			uintptr(unsafe.Pointer(&resume)),
		)
		if buffer != nil {
			members := unsafe.Slice(buffer, read)
			for _, member := range members {
				if member.SID != nil && member.SID.Equals(want) {
					_, _, _ = procNetApiBufferFree.Call(uintptr(unsafe.Pointer(buffer)))
					return true, nil
				}
			}
			_, _, _ = procNetApiBufferFree.Call(uintptr(unsafe.Pointer(buffer)))
		}
		if status == 0 {
			return false, nil
		}
		if status != uintptr(windows.ERROR_MORE_DATA) {
			return false, fmt.Errorf("NetLocalGroupGetMembers: %w", syscall.Errno(status))
		}
	}
}

func appendDiagnosis(existing, addition string) string {
	if strings.TrimSpace(existing) == "" {
		return addition
	}
	return existing + "; " + addition
}
