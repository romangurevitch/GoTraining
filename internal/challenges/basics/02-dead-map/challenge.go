package deadmap

// Ledger records transaction amounts by ID.
type Ledger struct {
	entries map[string]int64
	// BUG: entries is never initialized — it's nil until make() is called.
}

// NewLedger returns a Ledger ready to use.
//
// TODO: initialize the entries map so Record doesn't panic.
func NewLedger() *Ledger {
	return &Ledger{entries: make(map[string]int64)}
}

// Record stores an amount for the given transaction ID.
func (l *Ledger) Record(id string, amount int64) {
	l.entries[id] = amount // PANIC when entries is nil
}

// Balance returns the recorded amount for id (0 if not found).
func (l *Ledger) Balance(id string) int64 {
	return l.entries[id]
}
