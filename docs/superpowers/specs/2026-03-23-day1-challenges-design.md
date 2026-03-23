# Day 1 Challenges — Design Spec

**Date:** 2026-03-23
**Scope:** Go Training Day 1 — basic challenges for participants transferring from Python/other languages
**Status:** Approved for planning

---

## Relationship to Existing Plan

`docs/plan/challenges.md` contains a comprehensive vision: **105 "Detective Scenario" fixme challenges** across 21 topics. That plan is the long-term library. This spec defines the **Day 1 curated subset**: the 10 challenges participants work through during the afternoon sessions.

The detective brief framing, scenario naming conventions ("The Dead Map", "The Phantom Append"), and educational philosophy from `challenges.md` are adopted wholesale here. The per-challenge package structure (vs the flat single-package approach in `challenges.md`) is chosen because participants need to run each challenge independently for a clear sense of accomplishment.

**Scenarios adopted directly from `challenges.md`:**
- "The Dead Map" → Challenge 02
- "The Phantom Append" / "The Slice Reset" → Challenge 03
- "The Interface Lie" → Challenge 06
- "The Snapshot Method" → Challenge 01 (renamed "The Frozen Account" for banking domain fit)
- "The Sentinel Check" / "The Ignored Error" → Challenge 04
- "The Ignored Cancel" → Challenge 07

---

## Problem Statement

The existing challenges (01–09) are inadequate: dry, ambiguous, and don't teach meaningful Go concepts. They use `return 0` stubs that silently pass wrong values, lack narrative context, and don't reflect real Go gotchas that Python developers hit.

The new set must be **engaging, fun, and provide meaningful learning**. Each challenge is a **Detective Brief** — a small mystery that reveals a Go mental model mismatch. Participants don't just fix code; they build intuition for Go's runtime behavior.

---

## Design Principles

1. **Detective Brief framing** — each challenge presents a mystery: "Alice's balance never changes. The code compiles. What's broken?" Participants are detectives, not students.
2. **Fail fast, fail clearly** — `challenge_test.go` is pre-written and already failing. Participants see a red test immediately, fix it, see green. One challenge = one "aha" moment.
3. **Implme stubs use `panic("implement me")`** — not `return 0` or `return nil`. Wrong placeholder values mask real failures.
4. **Fixme code compiles and runs but behaves wrong** — no compile errors in the test file; participants run tests, see a clean `FAIL` message, read the code, find the bug.
5. **Banking domain flavor** — short backstory connects challenges to CBA's domain and builds anticipation for Day 2's bank challenge.
6. **Each challenge is its own package** — run independently with `go test -v ./internal/challenges/basics/NN-name/...`
7. **Progressive unlock** — challenges are grouped by lecture slot. Participants work through the current track while concepts are being taught. Fast finishers move to bonus challenges.
8. **Tests are read-only** — students edit only the designated file. The test file must not require modification.
9. **No external dependencies** — stdlib only for all Day 1 challenges.

---

## Challenge Format Types

- **fixme** — broken code that compiles and runs, but tests fail at assertion (or panic, caught by test-level `recover`). Student edits `challenge.go` to fix the bug.
- **implme** — function stubs with `panic("implement me")`. Student implements the logic. Tests drive the implementation.
- **testme** (bonus only) — logic is already implemented with a hidden bug. Student writes table-driven tests that expose the bug. Labelled explicitly in the README.

---

## Directory Structure

```
internal/challenges/basics/
├── 01-frozen-account/          # fixme  | receivers: value vs pointer mutation
├── 02-dead-map/                # fixme  | maps: nil map write (with recover)
├── 03-phantom-append/          # fixme  | slices: append semantics, value receiver compound
├── 04-silent-errors/           # implme | errors: multi-return, sentinel, errors.Is
├── 05-fee-calculator/          # implme | interfaces: satisfaction, pointer receiver method set
├── 06-interface-lie/           # fixme  | interfaces: nil concrete in interface != nil
├── 07-ignored-cancel/          # implme | context: WithCancel, goroutine cleanup, defer
├── 08-account-greeter/         # implme | packages: exported vs unexported, basic test
├── 09-interest-bug/            # testme | testing: table-driven, t.Helper — BONUS
└── 10-transfer-error/          # implme | errors: custom types, errors.As — BONUS
```

---

## Track 1: Go Building Blocks (unlocked at 13:30)

*Lecture covers: pointers, structs, slices, errors*
*Core: ~31 min | Bonus time if needed: 09 and 10 are accessible at any point*

---

### Challenge 01: The Frozen Account *(fixme)*

**Detective Brief:** *"Alice deposits $100. The code compiles. The function runs. Her balance is still $0. The teller is baffled. Find the bug."*

**Concept:** Value receivers operate on a copy of the struct. In Python, `self` is always a reference to the object. In Go, `func (a Account)` copies the struct — mutations vanish.

**The bug:**
```go
type Account struct { Balance int64 }

// BUG: value receiver — mutates a copy, not the original
func (a Account) Deposit(amount int64) {
    a.Balance += amount
}
```

**The fix:** Change `(a Account)` → `(a *Account)`. One token. Total "aha."

**Test:** Creates `&Account{Balance: 0}`, calls `Deposit(100)`, asserts `account.Balance == 100` → fails (stays 0 until fixed).

**`defer` cameo:** Test uses `defer` for cleanup — introduces the keyword without making it the lesson.

**Key lesson:** Go structs are values. `func (a Account)` is a snapshot method. Use `func (a *Account)` to mutate.

**Source in `challenges.md`:** Topic 10, scenario 1 — "The Snapshot Method".

**Estimated time:** 8 min

---

### Challenge 02: The Dead Map *(fixme)*

**Detective Brief:** *"The transaction logger runs without crashing — until it does. The map was declared. Was it ever born?"*

**Concept:** `var m map[string]int64` declares a nil map. Reads return zero silently. Writes panic. Python's `{}` gives a ready dict; Go's map literal doesn't.

**The bug:**
```go
type Ledger struct {
    entries map[string]int64 // never initialized
}

func (l *Ledger) Record(id string, amount int64) {
    l.entries[id] = amount // PANIC: assignment to entry in nil map
}
```

**Test pattern:** Uses `defer recover()` to catch the panic and print a friendly, actionable message:
```go
func TestRecord(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Fatalf("nil map panic — initialize the map first.\n  hint: make(map[string]int64)\n  got: %v", r)
        }
    }()
    l := &Ledger{}
    l.Record("tx-1", 10000)
    // assert entries["tx-1"] == 10000
}
```

**The fix:** Add `NewLedger() *Ledger { return &Ledger{entries: make(map[string]int64)} }` and use it.

**Key lesson:** `var m map[T]V` = nil. `make(map[T]V)` = live. Python `{}` = Go `map[T]V{}`.

**Source in `challenges.md`:** Topic 4, scenario 1 — "The Dead Map".

**Estimated time:** 5 min

---

### Challenge 03: The Phantom Append *(fixme)*

**Detective Brief:** *"The fraud team adds suspects to their watch list. The list is always empty. Three suspects. Zero entries. The append ran. Nothing stuck."*

**Concept:** Two compounding bugs — classic Go trap for Python developers who know `list.append()` mutates in place.

**The bug:**
```go
type WatchList struct { suspects []string }

// BUG 1: value receiver — works on a copy
// BUG 2: append return value ignored (moot since BUG 1 applies, but instructive)
func (w WatchList) AddSuspect(id string) {
    w.suspects = append(w.suspects, id)
}
```

**Test:** Adds 3 suspects, asserts `len == 3` → gets 0.

**The fix:** `func (w *WatchList) AddSuspect(id string) { w.suspects = append(w.suspects, id) }`

**Key lesson:** Python's `list.append()` mutates the list. Go's `append()` returns a new slice header. Always capture the return. Always use a pointer receiver if you need the struct to change.

**Source in `challenges.md`:** Topic 1, scenario 1 — "The Phantom Append" + Topic 10, scenario 1 (value receiver compounding).

**Estimated time:** 8 min

---

### Challenge 04: Silent Errors *(implme)*

**Detective Brief:** *"The payment processor silently swallowed a failed withdrawal. No exception. No crash. No trace. Compliance is not happy."*

**Concept:** Go functions return errors as values. Callers must check them. Sentinel errors let callers distinguish failure cases without string matching.

**Sentinel errors:**
- `ErrInsufficientFunds` — balance < amount
- `ErrNegativeAmount` — amount ≤ 0

**Stub:**
```go
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrNegativeAmount    = errors.New("negative amount")

// Withdraw deducts amount from balance.
// Returns ErrNegativeAmount if amount <= 0.
// Returns ErrInsufficientFunds if balance < amount.
func Withdraw(balance, amount int64) (int64, error) {
    panic("implement me")
}
```

**Test covers:** success, insufficient funds (`errors.Is`), negative amount (`errors.Is`).

**Key lesson:** In Python `raise ValueError`. In Go `return 0, ErrNegativeAmount`. `errors.Is` traverses wrapped errors — never use `err.Error() == "some string"`.

**Source in `challenges.md`:** Topic 13, scenarios 1 & 5 ("The Ignored Error", "The Fragile String") + Topic 14, scenario 1 ("The Sentinel Check").

**Estimated time:** 10 min

---

## Track 2: Interfaces, Receivers & Context (unlocked at 14:30)

*Lecture covers: interfaces, receivers, context*

---

### Challenge 05: The Fee Calculator *(implme)*

**Detective Brief:** *"Two account types need to report monthly fees. One compiles. The other won't budge. The interface is the same — why does Go refuse one of them?"*

**Concept:** Interface satisfaction is implicit and checked at compile time. A pointer receiver method is only in `*T`'s method set — passing `T` (value) to an interface that needs `*T` won't compile.

**Code structure:**
```go
// FeeCalculator is defined here — consumer side.
type FeeCalculator interface {
    MonthlyFee() int64
}

// SavingsAccount — use a VALUE receiver (student implements this first, it "just works")
type SavingsAccount struct{ Balance int64 }

// PremiumAccount — must use a POINTER receiver (student discovers why)
type PremiumAccount struct{ Balance int64 }

// TotalFees — already provided, student does not modify
func TotalFees(accounts []FeeCalculator) int64 { ... }
```

The test passes `SavingsAccount{}` (value) and `&PremiumAccount{}` (pointer) into `[]FeeCalculator`. The test file compiles cleanly. If the student uses a value receiver for `PremiumAccount`, *their* file won't compile — the error is in `challenge.go`.

**Key lesson:** Python duck typing = runtime. Go interfaces = compile time. `func (p PremiumAccount) MonthlyFee()` doesn't satisfy the interface for `*PremiumAccount`.

**Source in `challenges.md`:** Topic 11, scenario 4 — "The Pointer Requirement".

**Estimated time:** 10 min

---

### Challenge 06: The Interface Lie *(fixme)*

**Detective Brief:** *"The AML checker says it found no error. But the error handler fires anyway. The checker was skipped. The error is nil — or is it?"*

**Concept:** Returning a typed nil (`(*AMLChecker)(nil)`) as an `error` interface produces a non-nil interface. The interface holds `(type=*AMLChecker, value=nil)` — which is not `== nil`.

**The bug:**
```go
func runAMLCheck(skip bool) error {
    var checker *AMLChecker
    if skip {
        return checker // BUG: typed nil — interface is NOT nil
    }
    return checker.Run()
}

// Caller:
if err := runAMLCheck(true); err != nil {
    // This branch fires even when skip=true
}
```

**The fix:** `return nil` — untyped nil. The interface has no type, so it's truly nil.

**Test:** Calls `runAMLCheck(true)`, asserts `err == nil` → fails until fix.

**Key lesson:** An `error` interface is `(type, value)`. A typed nil pointer in an interface is NOT nil. Always return bare `nil` from functions returning interface types.

**Source in `challenges.md`:** Topic 11, scenario 1 — "The Interface Lie". Topic 14, scenario 4 — "The Typed Nil Return".

**Estimated time:** 8 min

---

### Challenge 07: The Ignored Cancel *(implme)*

**Detective Brief:** *"The fraud monitor launched 50 goroutines. The request completed. The goroutines are still running. Memory is climbing. The context was cancelled. Nobody listened."*

**Concept:** Go goroutines don't stop automatically when a context is cancelled. You must explicitly `select` on `ctx.Done()`. Always `defer cancel()`.

**Stub (`challenge.go`):**
```go
// Watch monitors accountID and calls report on each check.
// It must stop when ctx is cancelled.
func Watch(ctx context.Context, accountID string, report func(string), stopped chan struct{}) {
    go func() {
        defer close(stopped) // signals the test: goroutine has exited
        for {
            select {
            // TODO: case <-ctx.Done(): return
            // TODO: default: report(accountID)
            }
            panic("implement me")
        }
    }()
}
```

**Race-free test design** — uses a `stopped` channel, no `time.Sleep`:
```go
func TestWatch_StopsOnCancel(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())

    reported := make(chan struct{}, 1)
    stopped := make(chan struct{})

    Watch(ctx, "acc-123", func(string) {
        select { case reported <- struct{}{}: default: }
    }, stopped)

    select {
    case <-reported:
    case <-time.After(time.Second):
        t.Fatal("Watch never called report — did you start a goroutine?")
    }

    cancel()

    select {
    case <-stopped:
    case <-time.After(time.Second):
        t.Fatal("goroutine did not stop after context cancellation")
    }
}
```

**`defer` introduced naturally:** The scaffold shows `defer cancel()` in the test and `defer close(stopped)` in the goroutine. README notes both.

**Key lesson:** Python `asyncio` tasks auto-cancel on exception. Go goroutines don't. Always `select` on `ctx.Done()`. Always `defer cancel()`.

**Source in `challenges.md`:** Topic 20, scenario 3 — "The Ignored Cancel". Topic 20, scenario 2 — "The Timer Leak".

**Estimated time:** 12 min

---

## Track 3: Build, Package & Run (unlocked at 15:20)

*Lecture covers: package layout, testing*

---

### Challenge 08: The Account Greeter *(implme)*

**Detective Brief:** *"Build the greeting service — the simplest package CBA ships. It must be correct, exported correctly, and testable. No excuses."*

**Concept:** Exported vs unexported identifiers are the only access control in Go. Capitalisation is the modifier. The compiler enforces it.

**Editable file:** `greeting.go` (README calls this out — exception to the `challenge.go` convention).

**Stub:**
```go
// formatName concatenates first and last name with a space.
// Unexported — only usable within this package.
func formatName(first, last string) string {
    panic("implement me")
}

// Greet returns a personalised welcome message for a new account.
// Exported — callable from any package.
func Greet(accountID, first, last string) string {
    panic("implement me")
    // Expected: "Hello, <First Last>! Your account <accountID> is ready."
}
```

**Test:** Straight assertions (no table) — this challenge is a model for *writing* tests, not for teaching table-driven tests.

**Key lesson:** Python uses `_name` convention. Go uses capitalisation — `lowercase` = package-private, `Uppercase` = exported. The compiler enforces this; there's no workaround.

**Estimated time:** 10 min

---

### Challenge 09: The Interest Bug *(testme — BONUS)*

**Detective Brief:** *"The interest calculator passed code review. The team wants to ship it. Your job: write the tests that prove it's wrong before it hits production."*

**Concept:** Table-driven tests force you to enumerate edge cases you'd miss with ad-hoc tests. `t.Helper()` points failure messages at the failing table row, not the helper function.

**Code structure:**
- `interest.go` — fully implemented, contains one deliberate bug: negative `rate` returns a positive result (should return `ErrNegativeRate`)
- `interest_test.go` stub — student writes ≥5 table cases; the negative-rate case must be included to trigger the bug
- `checkInterest` helper with `t.Helper()` — pre-written and ready to use

**Key lesson:** Table-driven tests aren't just ergonomic — they make edge case coverage systematic. `t.Helper()` makes failure messages readable.

**Source in `challenges.md`:** Topic 21, scenario 1 — "The Wrong Line".

**Estimated time:** 10 min

---

### Challenge 10: The Transfer Error *(implme — BONUS)*

**Detective Brief:** *"A transfer failed. The log says 'error'. Which account? How much? What rule? The team needs answers, not just 'error'."*

**Concept:** Custom error types carry structured data. `errors.As` extracts the concrete type. `errors.Is` matches by sentinel identity. Both are more robust than string comparison.

**Code structure:**
```go
// TransferError carries structured context about a failed transfer.
type TransferError struct {
    FromID string
    ToID   string
    Amount int64
    Reason string
}

func (e *TransferError) Error() string {
    panic("implement me")
}

// Transfer validates and records a transfer.
// Returns *TransferError for business rule violations.
func Transfer(fromID, toID string, amount int64) error {
    panic("implement me")
}
```

**Test:** Calls `Transfer` with bad data, uses `errors.As` to extract `*TransferError` and asserts its fields; also shows `errors.Is` with a sentinel for comparison.

**Key lesson:** Python exceptions have arbitrary attributes. Go errors are values — use custom types for structured context. `errors.As` unwraps; `errors.Is` matches identity. Neither relies on string comparison.

**Source in `challenges.md`:** Topic 14, scenario 2 — "The Custom Field Access". Topic 13, scenario 5 — "The Fragile String".

**Estimated time:** 12 min

---

## File Structure Per Challenge

```
NN-challenge-name/
├── README.md           # Detective brief, goal, run command, key lesson, hints
├── challenge.go        # Editable file (may be named differently — README specifies)
└── challenge_test.go   # Pre-written, read-only, fails before fix/impl
```

**README sections:**
1. **Detective Brief** — vivid 1-2 sentence mystery hook
2. **The Crime Scene** — what the broken code looks like (for fixme) or what to build (for implme)
3. Run command: `go test -v ./internal/challenges/basics/NN-name/...`
4. **Key lesson** — bold Python vs Go mental model comparison
5. `<details><summary>Hints (click to reveal)</summary>` — 3-5 progressive hints

---

## Fast Finisher Path

Participants who complete Track 1 challenges before the 14:30 Track 2 lecture:
1. Move to **bonus challenges 09 and 10** — no lecture required, they build on already-known concepts
2. Start reading the **bank challenge README** (`internal/challenges/bank/`) to preview Day 2
3. If truly done: browse `internal/basics/` demo code for the upcoming lecture topic

---

## Quality Standards

| Standard | Rule |
|---|---|
| Tests fail before fix | Verify with `go test` on unmodified `challenge.go` before authoring |
| Tests are read-only | Students edit only the explicitly designated file |
| Fixme bugs are subtle | Bug readable but not obvious from a Python mental model |
| No external dependencies | stdlib only for all Day 1 challenges |
| Panic bugs use recover | Test catches panic, prints a friendly, actionable message |
| No compile errors in test files | Compile errors must be in the student's file, never the test file |
| `defer` appears naturally | Introduced in Challenge 07 scaffold without a dedicated challenge |
| File naming documented | Challenge 08 uses `greeting.go` — README and spec both call this out |

---

## Connection to `challenges.md` Long-Term Plan

This Day 1 set is the **entry point** to the full 105-scenario library. After Day 1, the `fixme/` flat-package approach from `challenges.md` can be used as a self-study resource covering all 21 topics. The two approaches are complementary:

| Day 1 Spec | `challenges.md` Library |
|---|---|
| 10 challenges, own packages | 105 scenarios, single flat package |
| Banking domain story | Generic but vivid "detective" names |
| fixme + implme + testme | fixme only |
| Taught in sequence with lectures | Self-study, any order |
| Day 1 scope (pointers–context) | All 21 Go topics including concurrency |

The `challenges.md` library is built in batches per the existing plan. Day 1 challenges are implemented first (this spec).

---

## Out of Scope (Day 1)

- HTTP challenges (covered Day 2)
- Mocking / GoMock (covered Day 2)
- Concurrency deep-dive (separate ConcurrencyWorkshop)
- Generics (advanced, not Day 1)
- Arrays (interesting but rarely a Day 1 gotcha vs Python)
- Goroutines / channels standalone (Challenge 07 introduces context cancellation as a proxy)
- `defer` as a standalone challenge (introduced naturally in Challenges 01 and 07)
