//go:build !windows

package service

import (
	"github.com/kardianos/service"
)

func platformQueryState(name string) (State, error) {
	prog := &program{}
	svc, err := service.New(prog, baseConfig("", nil))
	if err != nil {
		return StateUnknown, err
	}
	st, err := svc.Status()
	if err != nil {
		return StateUnknown, err
	}
	switch st {
	case service.StatusRunning:
		return StateRunning, nil
	case service.StatusStopped:
		return StateStopped, nil
	default:
		return StateUnknown, nil
	}
}
