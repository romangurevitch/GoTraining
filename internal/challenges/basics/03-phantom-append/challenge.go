package phantomappend

// WatchList tracks suspicious account IDs.
type WatchList struct {
	suspects []string
}

// AddSuspect adds an account ID to the watch list.
//
// BUG 1: value receiver — this method operates on a copy of WatchList.
// BUG 2: even if the receiver were a pointer, append must be assigned back.
//
// TODO: fix both bugs so suspects are actually stored.
func (w WatchList) AddSuspect(id string) {
	w.suspects = append(w.suspects, id) // nolint:staticcheck // intentional bug for challenge
}

// Count returns the number of suspects on the list.
func (w *WatchList) Count() int {
	return len(w.suspects)
}

// Contains reports whether id is on the watch list.
func (w *WatchList) Contains(id string) bool {
	for _, s := range w.suspects {
		if s == id {
			return true
		}
	}
	return false
}
