package domain

import (
	"strings"
	"testing"
	"time"
)

func TestCapabilityHierarchy(t *testing.T) {
	levels := []Capability{CapabilityObserve, CapabilityOperate, CapabilityManage, CapabilityEnroll}
	for i, actual := range levels {
		for j, required := range levels {
			if got, want := actual.Allows(required), i >= j; got != want {
				t.Fatalf("%s allows %s = %v, want %v", actual, required, got, want)
			}
		}
	}
	if Capability("owner").Valid() || Capability("owner").Allows(CapabilityObserve) || CapabilityEnroll.Allows(Capability("unknown")) {
		t.Fatal("unknown capabilities must fail closed")
	}
}

func TestActorValidationAndExpiration(t *testing.T) {
	if got, err := NormalizeActorDisplayName("  Desktop client  "); err != nil || got != "Desktop client" {
		t.Fatalf("normalize=%q err=%v", got, err)
	}
	for _, value := range []string{"", "\nclient", strings.Repeat("x", 81)} {
		if _, err := NormalizeActorDisplayName(value); err == nil {
			t.Fatalf("accepted invalid name %q", value)
		}
	}
	now := time.Now()
	future := now.Add(time.Minute)
	past := now.Add(-time.Minute)
	if !(Actor{State: ActorStateActive, ExpiresAt: &future}).ActiveAt(now) {
		t.Fatal("future actor should be active")
	}
	if (Actor{State: ActorStateActive, ExpiresAt: &past}).ActiveAt(now) {
		t.Fatal("expired actor should fail closed")
	}
	if ActorState("pending").Valid() || ActorKind("browser").Valid() || AuditResult("ok").Valid() {
		t.Fatal("unknown enums must be invalid")
	}
}
