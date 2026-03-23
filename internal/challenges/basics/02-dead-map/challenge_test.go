package deadmap

import "testing"

func TestRecord_StoresEntry(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf(
				"nil map panic — the entries map was never initialized.\n"+
					"  Fix: use NewLedger() or initialize entries with make(map[string]float64).\n"+
					"  Panic: %v",
				r,
			)
		}
	}()

	l := NewLedger()
	l.Record("tx-1", 100.0)

	got := l.Balance("tx-1")
	if got != 100.0 {
		t.Fatalf("expected balance 100.0 for tx-1, got %.2f", got)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("nil map panic: %v\n  Did you use NewLedger()?", r)
		}
	}()

	l := NewLedger()
	l.Record("tx-1", 50.0)
	l.Record("tx-2", 75.0)

	if l.Balance("tx-1") != 50.0 {
		t.Errorf("tx-1: expected 50.0, got %.2f", l.Balance("tx-1"))
	}
	if l.Balance("tx-2") != 75.0 {
		t.Errorf("tx-2: expected 75.0, got %.2f", l.Balance("tx-2"))
	}
}
