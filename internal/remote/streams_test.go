package remote

import "testing"

func TestStreamSetSeparatesGlobalAndPerActorCapacity(t *testing.T) {
	set := newStreamSet(3, 2)
	if !set.acquire("first") {
		t.Fatal("actor did not receive its first stream slot")
	}
	if !set.acquire("first") {
		t.Fatal("actor did not receive its second stream slot")
	}
	if set.acquire("first") {
		t.Fatal("actor exceeded its stream limit")
	}
	if !set.acquire("second") || set.acquire("third") {
		t.Fatal("global stream limit was not enforced")
	}
	set.release("first")
	if !set.acquire("third") {
		t.Fatal("released stream slot was not reusable")
	}
}
