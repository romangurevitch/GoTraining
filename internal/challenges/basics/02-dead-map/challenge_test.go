package deadmap

import "testing"

func TestRecord_StoresEntry(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf(
				"nil map panic — the entries map was never initialized.\n"+
					"  Fix: use NewLedger() or initialize entries with make(map[string]int64).\n"+
					"  Panic: %v",
				r,
			)
		}
	}()

	l := NewLedger()
	l.Record("tx-1", 100)

	got := l.Balance("tx-1")
	if got != 100 {
		t.Fatalf("expected balance 100 for tx-1, got %d", got)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("nil map panic: %v\n  Did you use NewLedger()?", r)
		}
	}()

	l := NewLedger()
	l.Record("tx-1", 50)
	l.Record("tx-2", 75)

	if l.Balance("tx-1") != 50 {
		t.Errorf("tx-1: expected 50, got %d", l.Balance("tx-1"))
	}
	if l.Balance("tx-2") != 75 {
		t.Errorf("tx-2: expected 75, got %d", l.Balance("tx-2"))
	}
}
