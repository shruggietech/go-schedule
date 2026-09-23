//go:build linux

package service

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
)

func TestInstalledConfigArg(t *testing.T) {
	cases := []struct {
		argv []string
		want string
		err  bool
	}{
		{[]string{"/usr/local/bin/goschedd"}, "", false},
		{[]string{"/usr/local/bin/goschedd", "--config", "/etc/goschedule/custom config.json"}, "/etc/goschedule/custom config.json", false},
		{[]string{"/usr/local/bin/goschedd", "--config=/etc/goschedule/custom config.json"}, "/etc/goschedule/custom config.json", false},
		{[]string{"/usr/local/bin/goschedd", "--config"}, "", true},
		{[]string{"/usr/local/bin/goschedd", "--config="}, "", true},
		{[]string{"/usr/local/bin/goschedd", "--config", "relative.json"}, "", true},
		{[]string{"/usr/local/bin/goschedd", "--config", "/a", "--config", "/b"}, "", true},
		{[]string{"/usr/local/bin/goschedd", "--config=/a", "--config", "/b"}, "", true},
	}
	for _, tc := range cases {
		got, err := configArg(tc.argv)
		if got != tc.want || (err != nil) != tc.err {
			t.Errorf("configArg(%q) = %q, %v; want %q, err=%v", tc.argv, got, err, tc.want, tc.err)
		}
	}
}

func TestSystemdExecVariant(t *testing.T) {
	want := []systemdExec{{Path: "/usr/local/bin/goschedd", Argv: []string{"/usr/local/bin/goschedd", "--config", "/etc/custom config.json"}}}
	property := dbus.MakeVariant(want)
	var got []systemdExec
	if err := property.Store(&got); err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("decode ExecStart: %v, got %v", err, got)
	}
}

func TestSystemdExecArgsLiveIfAvailable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	argv, err := systemdExecArgs(ctx, "dbus.service")
	if err != nil && (strings.Contains(err.Error(), "connect to system service manager") || strings.Contains(err.Error(), "load dbus.service definition")) {
		t.Skipf("no systemd bus or dbus service: %v", err)
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(argv) == 0 {
		t.Fatal("empty effective ExecStart argv")
	}
}
