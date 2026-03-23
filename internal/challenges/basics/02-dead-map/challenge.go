package deadmap

// Ledger records transaction amounts by ID.
type Ledger struct {
	entries map[string]float64
	// BUG: entries is never initialized — it's nil until make() is called.
}

// NewLedger returns a Ledger ready to use.
//
// TODO: initialize the entries map so Record doesn't panic.
func NewLedger() *Ledger {
	return &Ledger{} // BUG: entries map is nil
}

// Record stores an amount for the given transaction ID.
func (l *Ledger) Record(id string, amount float64) {
	l.entries[id] = amount // PANIC when entries is nil
}

// Balance returns the recorded amount for id (0 if not found).
func (l *Ledger) Balance(id string) float64 {
	return l.entries[id]
}
