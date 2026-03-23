# 🔍 Case 03: The Phantom Append

## The Detective Brief
The fraud team adds three suspects to the watch list. The list is always empty. The append ran. Nothing stuck.

## The Crime Scene
```bash
go test -v ./internal/challenges/basics/03-phantom-append/...
```

## Your Mission
1. Run the tests — they fail.
2. Read `AddSuspect` in `challenge.go` — there are **two** bugs.
3. Fix both.
4. Run the tests again.

## Key Lesson
> **Python vs Go:** Python's `list.append()` mutates the list in place. Go's `append()` may allocate a new backing array and always returns the updated slice — you must capture the return value. Combined with a value receiver (which copies the struct), suspects vanish twice over.

<details>
<summary>Hints (click to reveal)</summary>

1. Bug 1: What type of receiver does `AddSuspect` use? What does that mean for mutations?
2. Bug 2: `append` in Go returns the updated slice — if you don't capture the return, the original is unchanged. (Moot here given Bug 1, but both need fixing.)
3. Fix 1: Change `(w WatchList)` to `(w *WatchList)`.
4. Fix 2: Ensure the line reads `w.suspects = append(w.suspects, id)` — it already does.
5. After both fixes, the WatchList will actually remember suspects.
</details>
