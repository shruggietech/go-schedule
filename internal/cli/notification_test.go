package cli

import "testing"

func TestNotificationCommandRegistersCompleteManagementSurface(t *testing.T) {
	root := newNotificationCmd()
	want := map[string]bool{"channel": true, "task": true, "group": true, "deliveries": true}
	for _, child := range root.Commands() {
		delete(want, child.Name())
	}
	if len(want) != 0 {
		t.Fatalf("missing notification commands: %v", want)
	}
	channel := newNotificationChannelCmd()
	channelWant := map[string]bool{"add": true, "list": true, "get": true, "update": true, "enable": true, "disable": true, "rotate": true, "test": true, "rm": true}
	for _, child := range channel.Commands() {
		delete(channelWant, child.Name())
	}
	if len(channelWant) != 0 {
		t.Fatalf("missing channel commands: %v", channelWant)
	}
}

func TestParseNotificationOutcomes(t *testing.T) {
	success, failure, err := parseNotificationOutcomes("success,failure")
	if err != nil || !success || !failure {
		t.Fatalf("success=%t failure=%t err=%v", success, failure, err)
	}
	if _, _, err := parseNotificationOutcomes("maybe"); err == nil {
		t.Fatal("accepted invalid outcome")
	}
}
