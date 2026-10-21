//go:build !windows

package clievaluate

import "errors"

func processIdentity(int) (string, error) {
	return "", errors.New("coordinator leases require the Windows host")
}
