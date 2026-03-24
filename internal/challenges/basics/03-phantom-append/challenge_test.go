package phantomappend

import "testing"

func TestAddSuspect_StoresSuspect(t *testing.T) {
	wl := &WatchList{}
	wl.AddSuspect("acc-001")

	if wl.Count() != 1 {
		t.Fatalf(
			"expected 1 suspect after AddSuspect, got %d.\n"+
				"  Hint: check the receiver type AND whether append's return value is captured.",
			wl.Count(),
		)
	}
}

func TestAddSuspect_StoresMultipleSuspects(t *testing.T) {
	wl := &WatchList{}
	wl.AddSuspect("acc-001")
	wl.AddSuspect("acc-002")
	wl.AddSuspect("acc-003")

	if wl.Count() != 3 {
		t.Fatalf("expected 3 suspects, got %d", wl.Count())
	}
}

func TestAddSuspect_ContainsSuspect(t *testing.T) {
	wl := &WatchList{}
	wl.AddSuspect("acc-007")

	if !wl.Contains("acc-007") {
		t.Fatal("expected WatchList to contain acc-007 after AddSuspect")
	}
}
