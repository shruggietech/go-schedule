//go:build windows

// Package secretstore protects daemon-owned secret values with platform facilities.
package secretstore

import (
	"encoding/base64"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const cryptProtectUIForbidden = 0x1

// Protect encrypts a value for the current Windows service identity with DPAPI.
func Protect(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	input := []byte(value)
	in := windows.DataBlob{Size: uint32(len(input)), Data: &input[0]}
	var out windows.DataBlob
	if err := windows.CryptProtectData(&in, nil, nil, 0, nil, cryptProtectUIForbidden, &out); err != nil {
		return "", fmt.Errorf("protect secret with DPAPI: %w", err)
	}
	defer func() { _, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data))) }()
	protected := unsafe.Slice(out.Data, out.Size)
	return "dpapi:" + base64.RawStdEncoding.EncodeToString(protected), nil
}

// Unprotect decrypts a value for the current Windows service identity.
func Unprotect(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	const prefix = "dpapi:"
	if len(value) <= len(prefix) || value[:len(prefix)] != prefix {
		return "", fmt.Errorf("protected secret has an invalid DPAPI envelope")
	}
	encoded, err := base64.RawStdEncoding.DecodeString(value[len(prefix):])
	if err != nil {
		return "", fmt.Errorf("decode protected secret: %w", err)
	}
	in := windows.DataBlob{Size: uint32(len(encoded)), Data: &encoded[0]}
	var out windows.DataBlob
	if err := windows.CryptUnprotectData(&in, nil, nil, 0, nil, cryptProtectUIForbidden, &out); err != nil {
		return "", fmt.Errorf("unprotect secret with DPAPI: %w", err)
	}
	defer func() { _, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data))) }()
	plain := unsafe.Slice(out.Data, out.Size)
	return string(plain), nil
}
