# Day 1 Challenges Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace ALL existing inadequate Day 1 challenges (the old 01–09 are all superseded) with 10 high-quality "Detective Brief" challenges that teach Go mental model shifts to Python/other-language developers.

**Architecture:** Each challenge is a standalone Go package under `internal/challenges/basics/NN-name/`. Every package contains exactly two files: `challenge.go` (editable, pre-broken or stubbed) and `challenge_test.go` (pre-written, read-only, fails until solved). The test file is the source of truth — write it first, then craft the broken/stub code to match.

**Tech Stack:** Go stdlib only (no testify, no external deps), `go test -v ./internal/challenges/basics/NN-name/...`

---

## Reference Documents

- **Spec:** `docs/superpowers/specs/2026-03-23-day1-challenges-design.md`
- **Long-term plan:** `docs/plan/challenges.md`
- **Module:** `github.com/romangurevitch/go-training`
- **Existing challenges (ALL superseded):** `internal/challenges/basics/01-*` through `09-*` — all nine old directories are replaced. Challenges 06–09 (package-refactoring, mock-payment-gateway, middleware-chain, worker-pool) are Day 2/advanced topics and do not belong in Day 1 basics. They are removed from this directory, not archived.

---

## Conventions (read before touching any file)

### Package names
Each package name is the last segment of the directory, with hyphens removed and made valid Go: `frozenaccount`, `deadmap`, `phantomappend`, `silenterrors`, `feecalculator`, `interfacelie`, `ignoredcancel`, `accountgreeter`, `interestbug`, `transfererror`.

### Test failure messages
Every `t.Fatal` / `t.Errorf` message must:
1. State what was expected vs what happened
2. Give a one-line hint pointing at the bug location
Example: `t.Fatalf("Balance should be 100 after Deposit — got %v. Hint: check the receiver type on Deposit()", account.Balance)`

### Panic recovery template (fixme challenges that panic)
```go
defer func() {
    if r := recover(); r != nil {
        t.Fatalf("<what panicked> — <friendly hint>\n  got panic: %v", r)
    }
}()
```

### README template
```markdown
# 🔍 Case NN: <Title>

## The Detective Brief
<1-2 sentence mystery hook>

## The Crime Scene
`go test -v ./internal/challenges/basics/NN-name/...`

## Your Mission
1. Run the tests — they fail.
2. Read `challenge.go`.
3. Find the bug / implement the function.
4. Run the tests again — they should pass.

## Key Lesson
> **Python vs Go:** <one-sentence mental model comparison>

<details>
<summary>Hints (click to reveal)</summary>

1. <hint 1>
2. <hint 2>
3. <hint 3>
</details>
```

### Verification after every challenge
```bash
# Must see FAIL before fix:
go test -v ./internal/challenges/basics/NN-name/...

# Must see PASS after fix (use reference solution to verify):
# Temporarily apply fix, run, revert
```

---

## File Map

| File | Action | Purpose |
|------|--------|---------|
| `internal/challenges/basics/01-frozen-account/challenge.go` | Create | Fixme: value receiver bug |
| `internal/challenges/basics/01-frozen-account/challenge_test.go` | Create | Pre-written test, read-only |
| `internal/challenges/basics/01-frozen-account/README.md` | Create | Detective brief |
| `internal/challenges/basics/02-dead-map/challenge.go` | Create | Fixme: nil map bug |
| `internal/challenges/basics/02-dead-map/challenge_test.go` | Create | Pre-written test with recover |
| `internal/challenges/basics/02-dead-map/README.md` | Create | Detective brief |
| `internal/challenges/basics/03-phantom-append/challenge.go` | Create | Fixme: append + value receiver |
| `internal/challenges/basics/03-phantom-append/challenge_test.go` | Create | Pre-written test |
| `internal/challenges/basics/03-phantom-append/README.md` | Create | Detective brief |
| `internal/challenges/basics/04-silent-errors/challenge.go` | Create | Implme: sentinel errors |
| `internal/challenges/basics/04-silent-errors/challenge_test.go` | Create | Pre-written test |
| `internal/challenges/basics/04-silent-errors/README.md` | Create | Detective brief |
| `internal/challenges/basics/05-fee-calculator/challenge.go` | Create | Implme: interface + pointer receiver |
| `internal/challenges/basics/05-fee-calculator/challenge_test.go` | Create | Pre-written test |
| `internal/challenges/basics/05-fee-calculator/README.md` | Create | Detective brief |
| `internal/challenges/basics/06-interface-lie/challenge.go` | Create | Fixme: typed nil in interface |
| `internal/challenges/basics/06-interface-lie/challenge_test.go` | Create | Pre-written test |
| `internal/challenges/basics/06-interface-lie/README.md` | Create | Detective brief |
| `internal/challenges/basics/07-ignored-cancel/challenge.go` | Create | Implme: context cancellation |
| `internal/challenges/basics/07-ignored-cancel/challenge_test.go` | Create | Pre-written race-free test |
| `internal/challenges/basics/07-ignored-cancel/README.md` | Create | Detective brief |
| `internal/challenges/basics/08-account-greeter/greeting.go` | Create | Implme: exported vs unexported |
| `internal/challenges/basics/08-account-greeter/greeting_test.go` | Create | Pre-written test |
| `internal/challenges/basics/08-account-greeter/README.md` | Create | Detective brief |
| `internal/challenges/basics/09-interest-bug/interest.go` | Create | Already implemented with hidden bug |
| `internal/challenges/basics/09-interest-bug/interest_test.go` | Create | Stub — student writes table cases |
| `internal/challenges/basics/09-interest-bug/README.md` | Create | Detective brief |
| `internal/challenges/basics/10-transfer-error/challenge.go` | Create | Implme: custom error type |
| `internal/challenges/basics/10-transfer-error/challenge_test.go` | Create | Pre-written test |
| `internal/challenges/basics/10-transfer-error/README.md` | Create | Detective brief |
| `internal/challenges/basics/README.md` | Modify | Update index to list all 10 challenges |
| Old challenges 01–05 | Delete | Remove superseded directories |

---

## Task 1: Challenge 01 — The Frozen Account (fixme)

**Concept:** Value receiver operates on a copy — mutations vanish.

**Files:**
- Create: `internal/challenges/basics/01-frozen-account/challenge_test.go`
- Create: `internal/challenges/basics/01-frozen-account/challenge.go`
- Create: `internal/challenges/basics/01-frozen-account/README.md`

- [ ] **Step 1.1: Create the test file (read-only)**

```go
// internal/challenges/basics/01-frozen-account/challenge_test.go
package frozenaccount

import "testing"

func TestDeposit_UpdatesBalance(t *testing.T) {
	account := &Account{Balance: 0}
	account.Deposit(10000)

	if account.Balance != 10000 {
		t.Fatalf(
			"Balance should be 10000 after Deposit(10000) — got %d.\n"+
				"  Hint: check whether Deposit uses a value or pointer receiver.",
			account.Balance,
		)
	}
}

func TestDeposit_MultipleDeposits(t *testing.T) {
	account := &Account{Balance: 5000}
	account.Deposit(2500)
	account.Deposit(2500)

	if account.Balance != 10000 {
		t.Fatalf(
			"Balance should be 10000 after two Deposit(2500) calls — got %d.",
			account.Balance,
		)
	}
}
```

- [ ] **Step 1.2: Create challenge.go with the bug**

```go
// internal/challenges/basics/01-frozen-account/challenge.go
package frozenaccount

// Account holds a bank account balance.
type Account struct {
	Balance int64
}

// Deposit adds amount to the account balance.
//
// BUG: this method uses a value receiver — it operates on a copy of Account,
// so the original balance is never updated.
//
// TODO: fix the receiver so deposits actually change the balance.
func (a Account) Deposit(amount int64) {
	a.Balance += amount
}
```

- [ ] **Step 1.3: Verify the test fails**

```bash
go test -v ./internal/challenges/basics/01-frozen-account/...
```
Expected: `FAIL` — "Balance should be 10000 after Deposit(10000) — got 0"

- [ ] **Step 1.4: Verify the test passes with the fix applied**

Temporarily change `(a Account)` to `(a *Account)` in `challenge.go`, run tests, then revert.
Expected: `PASS`

- [ ] **Step 1.5: Revert to buggy state and create README**

```markdown
# 🔍 Case 01: The Frozen Account

## The Detective Brief
Alice deposits $100 into her account. The code compiles. The function runs without error. Her balance is still $0. The teller is baffled.

## The Crime Scene
`go test -v ./internal/challenges/basics/01-frozen-account/...`

## Your Mission
1. Run the tests — they fail.
2. Read `challenge.go` — the bug is in the method signature.
3. Fix the bug.
4. Run the tests again — they should pass.

## Key Lesson
> **Python vs Go:** In Python, `self` is always a reference to the actual object. In Go, `func (a Account) Deposit(...)` receives a *copy* of the struct — mutations to `a` vanish when the function returns. Use `func (a *Account) Deposit(...)` to mutate the original.

<details>
<summary>Hints (click to reveal)</summary>

1. Look at the receiver type in the `Deposit` method declaration.
2. In Go, `(a Account)` means the method gets a copy of the struct. `(a *Account)` means it gets a pointer to the original.
3. Change exactly one token in the method signature.
4. After the fix, `account.Deposit(100)` will update the balance of the original `account` variable.
</details>
```

- [ ] **Step 1.6: Commit**

```bash
git add internal/challenges/basics/01-frozen-account/
git commit -m "feat(challenges): add 01-frozen-account — value vs pointer receiver fixme"
```

---

## Task 2: Challenge 02 — The Dead Map (fixme)

**Concept:** Nil map panics on write. `make()` is required.

**Files:**
- Create: `internal/challenges/basics/02-dead-map/challenge_test.go`
- Create: `internal/challenges/basics/02-dead-map/challenge.go`
- Create: `internal/challenges/basics/02-dead-map/README.md`

- [ ] **Step 2.1: Create the test file with panic recovery**

```go
// internal/challenges/basics/02-dead-map/challenge_test.go
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
	l.Record("tx-1", 10000)

	got := l.Balance("tx-1")
	if got != 10000 {
		t.Fatalf("expected balance 10000 for tx-1, got %d", got)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("nil map panic: %v\n  Did you use NewLedger()?", r)
		}
	}()

	l := NewLedger()
	l.Record("tx-1", 5000)
	l.Record("tx-2", 7500)

	if l.Balance("tx-1") != 5000 {
		t.Errorf("tx-1: expected 5000, got %d", l.Balance("tx-1"))
	}
	if l.Balance("tx-2") != 7500 {
		t.Errorf("tx-2: expected 7500, got %d", l.Balance("tx-2"))
	}
}
```

- [ ] **Step 2.2: Create challenge.go with the bug**

```go
// internal/challenges/basics/02-dead-map/challenge.go
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
	return &Ledger{} // BUG: entries map is nil
}

// Record stores an amount for the given transaction ID.
func (l *Ledger) Record(id string, amount int64) {
	l.entries[id] = amount // PANIC when entries is nil
}

// Balance returns the recorded amount for id (0 if not found).
func (l *Ledger) Balance(id string) int64 {
	return l.entries[id]
}
```

- [ ] **Step 2.3: Verify the test fails with a friendly panic message**

```bash
go test -v ./internal/challenges/basics/02-dead-map/...
```
Expected: `FAIL` — "nil map panic — the entries map was never initialized..."

- [ ] **Step 2.4: Verify test passes with the fix**

Temporarily change `NewLedger` to return `&Ledger{entries: make(map[string]int64)}`, run, revert.
Expected: `PASS`

- [ ] **Step 2.5: Revert to buggy state and create README**

```markdown
# 🔍 Case 02: The Dead Map

## The Detective Brief
The transaction logger runs — until it doesn't. The map was declared. Was it ever born? The logger panics at runtime, but looked fine in code review.

## The Crime Scene
`go test -v ./internal/challenges/basics/02-dead-map/...`

## Your Mission
1. Run the tests — they panic with a friendly message.
2. Read `NewLedger()` in `challenge.go`.
3. Initialize the map.
4. Run the tests again — they should pass.

## Key Lesson
> **Python vs Go:** In Python, `{}` gives you a ready-to-use dict. In Go, `map[string]int64` declares the type but the value is `nil` — you must call `make(map[string]int64)` or use a composite literal `map[string]int64{}` before writing to it.

<details>
<summary>Hints (click to reveal)</summary>

1. Look at the `NewLedger()` function — what is `entries` set to?
2. `var m map[string]int64` declares a nil map. Reads return zero values silently. Writes panic.
3. Fix: `entries: make(map[string]int64)` inside the struct literal in `NewLedger`.
4. Alternatively: `entries: map[string]int64{}` works too.
</details>
```

- [ ] **Step 2.6: Commit**

```bash
git add internal/challenges/basics/02-dead-map/
git commit -m "feat(challenges): add 02-dead-map — nil map initialization fixme"
```

---

## Task 3: Challenge 03 — The Phantom Append (fixme)

**Concept:** Value receiver + discarded `append` return = suspects that vanish.

**Files:**
- Create: `internal/challenges/basics/03-phantom-append/challenge_test.go`
- Create: `internal/challenges/basics/03-phantom-append/challenge.go`
- Create: `internal/challenges/basics/03-phantom-append/README.md`

- [ ] **Step 3.1: Create the test file**

```go
// internal/challenges/basics/03-phantom-append/challenge_test.go
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
```

- [ ] **Step 3.2: Create challenge.go with two compounding bugs**

```go
// internal/challenges/basics/03-phantom-append/challenge.go
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
	w.suspects = append(w.suspects, id)
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
```

- [ ] **Step 3.3: Verify test fails**

```bash
go test -v ./internal/challenges/basics/03-phantom-append/...
```
Expected: `FAIL` — "expected 1 suspect after AddSuspect, got 0"

- [ ] **Step 3.4: Verify test passes with fix**

Temporarily change receiver to `(w *WatchList)` and ensure `w.suspects = append(w.suspects, id)` (it already is), run, revert.
Expected: `PASS`

- [ ] **Step 3.5: Revert to buggy state and create README**

```markdown
# 🔍 Case 03: The Phantom Append

## The Detective Brief
The fraud team adds three suspects to the watch list. The list is always empty. The append ran. Nothing stuck.

## The Crime Scene
`go test -v ./internal/challenges/basics/03-phantom-append/...`

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
```

- [ ] **Step 3.6: Commit**

```bash
git add internal/challenges/basics/03-phantom-append/
git commit -m "feat(challenges): add 03-phantom-append — append semantics + value receiver fixme"
```

---

## Task 4: Challenge 04 — Silent Errors (implme)

**Concept:** Error as return value. Sentinel errors. `errors.Is`.

**Files:**
- Create: `internal/challenges/basics/04-silent-errors/challenge_test.go`
- Create: `internal/challenges/basics/04-silent-errors/challenge.go`
- Create: `internal/challenges/basics/04-silent-errors/README.md`

- [ ] **Step 4.1: Create the test file**

```go
// internal/challenges/basics/04-silent-errors/challenge_test.go
package silenterrors

import (
	"errors"
	"testing"
)

func TestWithdraw_Success(t *testing.T) {
	newBalance, err := Withdraw(10000, 4000)
	if err != nil {
		t.Fatalf("expected no error for valid withdrawal, got: %v", err)
	}
	if newBalance != 6000 {
		t.Fatalf("expected new balance 6000, got %d", newBalance)
	}
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	_, err := Withdraw(5000, 10000)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf(
			"expected ErrInsufficientFunds when balance < amount, got: %v\n"+
				"  Hint: use errors.Is, not == or string comparison.",
			err,
		)
	}
}

func TestWithdraw_NegativeAmount(t *testing.T) {
	_, err := Withdraw(10000, -1000)
	if !errors.Is(err, ErrNegativeAmount) {
		t.Fatalf("expected ErrNegativeAmount for amount <= 0, got: %v", err)
	}
}

func TestWithdraw_ZeroAmount(t *testing.T) {
	_, err := Withdraw(10000, 0)
	if !errors.Is(err, ErrNegativeAmount) {
		t.Fatalf("expected ErrNegativeAmount for amount=0, got: %v", err)
	}
}

func TestWithdraw_ExactBalance(t *testing.T) {
	newBalance, err := Withdraw(10000, 10000)
	if err != nil {
		t.Fatalf("expected no error when withdrawing exact balance, got: %v", err)
	}
	if newBalance != 0 {
		t.Fatalf("expected new balance 0, got %d", newBalance)
	}
}
```

- [ ] **Step 4.2: Create challenge.go with stubs**

```go
// internal/challenges/basics/04-silent-errors/challenge.go
package silenterrors

import "errors"

// ErrInsufficientFunds is returned when the withdrawal amount exceeds the balance.
var ErrInsufficientFunds = errors.New("insufficient funds")

// ErrNegativeAmount is returned when the withdrawal amount is zero or negative.
var ErrNegativeAmount = errors.New("negative or zero amount")

// Withdraw deducts amount from balance and returns the new balance.
//
// Returns (0, ErrNegativeAmount) if amount <= 0.
// Returns (0, ErrInsufficientFunds) if balance < amount.
// Returns (balance - amount, nil) on success.
func Withdraw(balance, amount int64) (int64, error) {
	panic("implement me")
}
```

- [ ] **Step 4.3: Verify tests fail with panic**

```bash
go test -v ./internal/challenges/basics/04-silent-errors/...
```
Expected: all tests panic with "implement me"

- [ ] **Step 4.4: Verify tests pass with correct implementation**

Temporarily implement `Withdraw`:
```go
func Withdraw(balance, amount int64) (int64, error) {
    if amount <= 0 {
        return 0, ErrNegativeAmount
    }
    if balance < amount {
        return 0, ErrInsufficientFunds
    }
    return balance - amount, nil
}
```
Run, revert to `panic("implement me")`.
Expected: `PASS`

- [ ] **Step 4.5: Create README**

```markdown
# 🔍 Case 04: Silent Errors

## The Detective Brief
The payment processor silently swallowed a failed withdrawal. No exception. No crash. No trace. The balance never changed. Compliance is not happy.

## The Crime Scene
`go test -v ./internal/challenges/basics/04-silent-errors/...`

## Your Mission
1. Run the tests — they panic (implement me).
2. Implement `Withdraw` in `challenge.go`.
3. Return the right sentinel error for each failure case.
4. Run the tests again — they should pass.

## Key Lesson
> **Python vs Go:** In Python you `raise ValueError(...)`. In Go you `return 0, ErrNegativeAmount`. Errors are values — callers must check them. Use `errors.Is(err, ErrX)` not `err.Error() == "some string"` — the latter breaks as soon as errors get wrapped.

<details>
<summary>Hints (click to reveal)</summary>

1. Check `amount <= 0` first and return `(0, ErrNegativeAmount)`.
2. Check `balance < amount` next and return `(0, ErrInsufficientFunds)`.
3. On success, return `(balance - amount, nil)`.
4. The test uses `errors.Is(err, ErrInsufficientFunds)` — your returned error must be (or wrap) the sentinel.
5. Never ignore the second return value from a Go function that returns `error`.
</details>
```

- [ ] **Step 4.6: Commit**

```bash
git add internal/challenges/basics/04-silent-errors/
git commit -m "feat(challenges): add 04-silent-errors — sentinel errors implme"
```

---

## Task 5: Challenge 05 — The Fee Calculator (implme)

**Concept:** Interface satisfaction. Pointer receiver method set.

**Files:**
- Create: `internal/challenges/basics/05-fee-calculator/challenge_test.go`
- Create: `internal/challenges/basics/05-fee-calculator/challenge.go`
- Create: `internal/challenges/basics/05-fee-calculator/README.md`

- [ ] **Step 5.1: Create the test file**

```go
// internal/challenges/basics/05-fee-calculator/challenge_test.go
package feecalculator

import "testing"

func TestSavingsAccount_MonthlyFee(t *testing.T) {
	// SavingsAccount uses a value receiver — both T and *T satisfy the interface.
	var calc FeeCalculator = SavingsAccount{Balance: 100000}
	if got := calc.MonthlyFee(); got != 500 {
		t.Fatalf("SavingsAccount.MonthlyFee() = %d, want 500", got)
	}
}

func TestPremiumAccount_MonthlyFee(t *testing.T) {
	// PremiumAccount must use a pointer receiver — only *PremiumAccount satisfies the interface.
	// If you use a value receiver, this line will not compile.
	var calc FeeCalculator = &PremiumAccount{Balance: 500000}
	if got := calc.MonthlyFee(); got != 2500 {
		t.Fatalf("PremiumAccount.MonthlyFee() = %d, want 2500", got)
	}
}

func TestTotalFees(t *testing.T) {
	accounts := []FeeCalculator{
		SavingsAccount{Balance: 100000},
		&PremiumAccount{Balance: 500000},
	}
	total := TotalFees(accounts)
	if total != 3000 {
		t.Fatalf("TotalFees() = %d, want 3000", total)
	}
}
```

- [ ] **Step 5.2: Create challenge.go with stubs**

```go
// internal/challenges/basics/05-fee-calculator/challenge.go
package feecalculator

// FeeCalculator calculates the monthly fee for an account.
// Defined on the consumer side — this package decides what it needs.
type FeeCalculator interface {
	MonthlyFee() int64
}

// SavingsAccount charges a flat $5/month fee.
// Use a VALUE receiver — both SavingsAccount and *SavingsAccount will satisfy FeeCalculator.
type SavingsAccount struct {
	Balance int64
}

// MonthlyFee returns the monthly fee for a SavingsAccount.
// TODO: implement — return 500
func (s SavingsAccount) MonthlyFee() int64 {
	panic("implement me")
}

// PremiumAccount charges a flat $25/month fee.
// Use a POINTER receiver — only *PremiumAccount will satisfy FeeCalculator.
// If you use a value receiver here, the test line `var calc FeeCalculator = &PremiumAccount{...}`
// will compile but `var calc FeeCalculator = PremiumAccount{...}` won't — and that's the lesson.
type PremiumAccount struct {
	Balance int64
}

// MonthlyFee returns the monthly fee for a PremiumAccount.
// TODO: implement using a POINTER receiver — return 2500
func (p *PremiumAccount) MonthlyFee() int64 {
	panic("implement me")
}

// TotalFees sums the monthly fees for all accounts.
// Do not modify this function.
func TotalFees(accounts []FeeCalculator) int64 {
	var total int64
	for _, a := range accounts {
		total += a.MonthlyFee()
	}
	return total
}
```

- [ ] **Step 5.3: Verify tests fail with panic**

```bash
go test -v ./internal/challenges/basics/05-fee-calculator/...
```
Expected: panic "implement me" on all tests

- [ ] **Step 5.4: Verify tests pass with implementation**

Implement both `MonthlyFee` methods (return 500 and 2500), run, revert.
Expected: `PASS`

- [ ] **Step 5.5: Create README**

```markdown
# 🔍 Case 05: The Fee Calculator

## The Detective Brief
Two account types need to report monthly fees. One compiles. The other won't budge. The interface is the same — why does Go refuse one of them?

## The Crime Scene
`go test -v ./internal/challenges/basics/05-fee-calculator/...`

## Your Mission
1. Run the tests — they panic (implement me).
2. Implement `MonthlyFee` for `SavingsAccount` with a **value** receiver.
3. Implement `MonthlyFee` for `PremiumAccount` with a **pointer** receiver.
4. Run the tests — they should pass.
5. Experiment: try a value receiver for `PremiumAccount`. What happens?

## Key Lesson
> **Python vs Go:** Python uses duck typing at runtime — if it has the method, it works. Go checks interface satisfaction at **compile time**. A pointer receiver method (`func (p *PremiumAccount) MonthlyFee()`) is only in `*PremiumAccount`'s method set, not `PremiumAccount`'s. Pass `&PremiumAccount{}`, not `PremiumAccount{}`.

<details>
<summary>Hints (click to reveal)</summary>

1. `SavingsAccount.MonthlyFee()` should use `(s SavingsAccount)` — value receiver. Return 500.
2. `PremiumAccount.MonthlyFee()` should use `(p *PremiumAccount)` — pointer receiver. Return 2500.
3. In the test, `SavingsAccount{}` (value) satisfies `FeeCalculator`. `&PremiumAccount{}` (pointer) satisfies it after your fix.
4. Try changing `*PremiumAccount` to `PremiumAccount` in the test — the compiler will reject it. That's the lesson.
5. `TotalFees` is already implemented — don't modify it.
</details>
```

- [ ] **Step 5.6: Commit**

```bash
git add internal/challenges/basics/05-fee-calculator/
git commit -m "feat(challenges): add 05-fee-calculator — interface + pointer receiver implme"
```

---

## Task 6: Challenge 06 — The Interface Lie (fixme)

**Concept:** Typed nil in interface is not nil.

**Files:**
- Create: `internal/challenges/basics/06-interface-lie/challenge_test.go`
- Create: `internal/challenges/basics/06-interface-lie/challenge.go`
- Create: `internal/challenges/basics/06-interface-lie/README.md`

- [ ] **Step 6.1: Create the test file**

```go
// internal/challenges/basics/06-interface-lie/challenge_test.go
package interfacelie

import "testing"

func TestRunAMLCheck_SkipReturnsNil(t *testing.T) {
	err := runAMLCheck(true)
	if err != nil {
		t.Fatalf(
			"runAMLCheck(skip=true) should return nil error, got: %v (%T)\n"+
				"  Hint: a typed nil (*AMLChecker)(nil) stored in an error interface is NOT nil.\n"+
				"  Fix: return bare nil instead of the typed variable.",
			err, err,
		)
	}
}

func TestRunAMLCheck_NoSkipReturnsNil(t *testing.T) {
	err := runAMLCheck(false)
	if err != nil {
		t.Fatalf("runAMLCheck(skip=false) should return nil for a passing check, got: %v", err)
	}
}
```

- [ ] **Step 6.2: Create challenge.go with the bug**

```go
// internal/challenges/basics/06-interface-lie/challenge.go
package interfacelie

import "fmt"

// AMLChecker performs anti-money-laundering checks.
type AMLChecker struct {
	threshold int64
}

// check simulates an AML check (always passes in this stub).
func (c *AMLChecker) check() error {
	return nil // real implementation would inspect transactions
}

// runAMLCheck runs an AML check unless skip is true.
//
// BUG: when skip is true, this returns a typed nil (*AMLChecker)(nil) as an error interface.
// The interface is (type=*AMLChecker, value=nil) — which is NOT == nil.
//
// TODO: make runAMLCheck return a genuinely nil error when skip is true.
func runAMLCheck(skip bool) error {
	var checker *AMLChecker
	if skip {
		return checker // BUG: typed nil — the interface is not nil
	}
	checker = &AMLChecker{threshold: 1000000}
	return checker.check()
}

// Explain prints what the interface holds — useful for debugging.
func Explain(err error) string {
	return fmt.Sprintf("err=%v, type=%T, isNil=%v", err, err, err == nil)
}
```

- [ ] **Step 6.3: Verify test fails**

```bash
go test -v ./internal/challenges/basics/06-interface-lie/...
```
Expected: `FAIL` — "runAMLCheck(skip=true) should return nil error, got: <nil> (*interfacelie.AMLChecker)"

- [ ] **Step 6.4: Verify fix works**

Temporarily change `return checker` to `return nil` in the `if skip` branch, run, revert.
Expected: `PASS`

- [ ] **Step 6.5: Create README**

```markdown
# 🔍 Case 06: The Interface Lie

## The Detective Brief
The AML checker says it found no error. But the error handler fires anyway. The checker was skipped. The variable is nil. Yet `err != nil` is true. This is Go's most famous gotcha.

## The Crime Scene
`go test -v ./internal/challenges/basics/06-interface-lie/...`

## Your Mission
1. Run the tests — they fail, even though the "error" looks nil.
2. Read `runAMLCheck` in `challenge.go`.
3. Understand why returning a typed nil pointer as `error` isn't the same as returning `nil`.
4. Fix it with a one-word change.

## Key Lesson
> **Python vs Go:** A Go interface holds two things: `(type, value)`. When you return `(*AMLChecker)(nil)` as an `error`, the interface is `(type=*AMLChecker, value=nil)` — the type is set, so the interface is NOT nil. Always return bare `nil` from functions returning interface types, never a typed nil variable.

<details>
<summary>Hints (click to reveal)</summary>

1. `var checker *AMLChecker` declares a nil pointer of type `*AMLChecker`.
2. Storing that nil pointer in an `error` interface creates `(type=*AMLChecker, value=nil)`.
3. That interface != nil, because the type field is set.
4. Fix: in the `if skip` branch, write `return nil` — not `return checker`.
5. The `Explain()` function in `challenge.go` can help you see what the interface actually holds.
</details>
```

- [ ] **Step 6.6: Commit**

```bash
git add internal/challenges/basics/06-interface-lie/
git commit -m "feat(challenges): add 06-interface-lie — typed nil in interface fixme"
```

---

## Task 7: Challenge 07 — The Ignored Cancel (implme)

**Concept:** `context.WithCancel`, goroutine lifecycle, `defer`, race-free testing.

**Files:**
- Create: `internal/challenges/basics/07-ignored-cancel/challenge_test.go`
- Create: `internal/challenges/basics/07-ignored-cancel/challenge.go`
- Create: `internal/challenges/basics/07-ignored-cancel/README.md`

- [ ] **Step 7.1: Create the test file**

```go
// internal/challenges/basics/07-ignored-cancel/challenge_test.go
package ignoredcancel

import (
	"context"
	"testing"
	"time"
)

func TestWatch_ReportsAccountID(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // always defer cancel — prevents goroutine leaks

	reported := make(chan string, 1)
	stopped := make(chan struct{})

	Watch(ctx, "acc-123", func(id string) {
		select {
		case reported <- id:
		default:
		}
	}, stopped)

	select {
	case got := <-reported:
		if got != "acc-123" {
			t.Fatalf("expected report for acc-123, got %q", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Watch never called report — did you start a goroutine?")
	}
}

func TestWatch_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	reported := make(chan struct{}, 1)
	stopped := make(chan struct{})

	Watch(ctx, "acc-999", func(string) {
		select {
		case reported <- struct{}{}:
		default:
		}
	}, stopped)

	// Wait for at least one report (confirms the goroutine started)
	select {
	case <-reported:
	case <-time.After(2 * time.Second):
		t.Fatal("Watch never called report — goroutine did not start")
	}

	// Cancel the context
	cancel()

	// Confirm the goroutine stopped (closed the stopped channel)
	select {
	case <-stopped:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("Watch goroutine did not stop after context cancellation — did you select on ctx.Done()?")
	}
}
```

- [ ] **Step 7.2: Create challenge.go with stub**

```go
// internal/challenges/basics/07-ignored-cancel/challenge.go
package ignoredcancel

import "context"

// Watch monitors accountID and calls report for each check cycle.
// It must stop when ctx is cancelled — no goroutine leaks allowed.
//
// stopped is closed by the goroutine when it exits. The test uses this
// to confirm the goroutine actually stopped (race-free, no time.Sleep).
//
// TODO: implement the goroutine inside Watch so that:
//  1. It calls report(accountID) on each iteration.
//  2. It selects on ctx.Done() and returns when the context is cancelled.
//  3. It closes stopped when it exits (use defer — it's already there).
func Watch(ctx context.Context, accountID string, report func(string), stopped chan struct{}) {
	go func() {
		defer close(stopped) // signals the test: this goroutine has exited
		for {
			select {
			case <-ctx.Done():
				// TODO: return here to stop the goroutine
				panic("implement me: return when context is cancelled")
			default:
				// TODO: call report(accountID) here
				panic("implement me: call report on each iteration")
			}
		}
	}()
}
```

**Note:** The stub has `panic("implement me")` in both branches. The goroutine starts, immediately hits `select`, picks either `ctx.Done()` (if already cancelled) or `default` (normally), then panics. The test's `select { case <-reported: ... }` will timeout because `report` is never called before the panic. The panic in the goroutine also surfaces as a test failure. This is the correct implme failing state — panics are visible, not silent timeouts.

- [ ] **Step 7.3: Verify tests fail**

```bash
go test -v -timeout 10s ./internal/challenges/basics/07-ignored-cancel/...
```
Expected: `FAIL` with timeout — "Watch never called report"

- [ ] **Step 7.4: Verify tests pass with implementation**

Temporarily fill in the two TODO cases:
```go
select {
case <-ctx.Done():
    return
default:
    report(accountID)
}
```
Run, revert.
Expected: `PASS`

- [ ] **Step 7.5: Create README**

```markdown
# 🔍 Case 07: The Ignored Cancel

## The Detective Brief
The fraud monitor launched goroutines for every request. The requests completed. The goroutines are still running. Memory climbs. The context was cancelled. Nobody listened.

## The Crime Scene
`go test -v -timeout 10s ./internal/challenges/basics/07-ignored-cancel/...`

## Your Mission
1. Run the tests — they timeout.
2. Read `Watch` in `challenge.go`.
3. Implement the two TODO cases inside the `select`.
4. Run the tests again — they should pass.

## Key Lesson
> **Python vs Go:** Python's `asyncio` tasks can be cancelled automatically. Go goroutines run until they return — you must explicitly `select` on `ctx.Done()` to stop them. Always `defer cancel()` when you create a context: if you forget, goroutines can leak for the lifetime of the process.

<details>
<summary>Hints (click to reveal)</summary>

1. A `select` with no cases blocks forever — the goroutine will never call `report`. Add two cases.
2. `case <-ctx.Done(): return` — exits the goroutine when the context is cancelled.
3. `default: report(accountID)` — runs when the context is not yet cancelled.
4. `defer close(stopped)` is already there — it signals the test that the goroutine has exited. Don't remove it.
5. Notice `defer cancel()` in the test. Always defer context cancellation — even if you expect it to be called anyway.
</details>
```

- [ ] **Step 7.6: Commit**

```bash
git add internal/challenges/basics/07-ignored-cancel/
git commit -m "feat(challenges): add 07-ignored-cancel — context cancellation implme"
```

---

## Task 8: Challenge 08 — The Account Greeter (implme)

**Concept:** Exported vs unexported identifiers. Package structure.

**Files:**
- Create: `internal/challenges/basics/08-account-greeter/greeting_test.go`
- Create: `internal/challenges/basics/08-account-greeter/greeting.go`
- Create: `internal/challenges/basics/08-account-greeter/README.md`

- [ ] **Step 8.1: Create the test file**

```go
// internal/challenges/basics/08-account-greeter/greeting_test.go
package accountgreeter

import "testing"

func TestGreet_ReturnsCorrectMessage(t *testing.T) {
	got := Greet("ACC-001", "Alice", "Smith")
	want := "Hello, Alice Smith! Your account ACC-001 is ready."

	if got != want {
		t.Fatalf("Greet() =\n  %q\nwant\n  %q", got, want)
	}
}

func TestGreet_DifferentName(t *testing.T) {
	got := Greet("ACC-999", "Bob", "Jones")
	want := "Hello, Bob Jones! Your account ACC-999 is ready."

	if got != want {
		t.Fatalf("Greet() =\n  %q\nwant\n  %q", got, want)
	}
}

// formatName is unexported — this test verifies its effect indirectly through Greet.
// You cannot call formatName from outside this package.
func TestGreet_UsesFormatName(t *testing.T) {
	got := Greet("X", "Mary", "Jane")
	if got != "Hello, Mary Jane! Your account X is ready." {
		t.Fatalf("unexpected greeting: %q", got)
	}
}
```

- [ ] **Step 8.2: Create greeting.go with stubs**

```go
// internal/challenges/basics/08-account-greeter/greeting.go
package accountgreeter

import "fmt"

// formatName concatenates first and last name with a space.
// This function is unexported — it can only be used within this package.
// (Lowercase first letter = package-private in Go.)
//
// TODO: implement this function.
func formatName(first, last string) string {
	panic("implement me")
}

// Greet returns a personalised welcome message.
// This function is exported — it can be called from any package.
// (Uppercase first letter = exported/public in Go.)
//
// Expected format: "Hello, <First Last>! Your account <accountID> is ready."
//
// TODO: implement using formatName and fmt.Sprintf.
func Greet(accountID, first, last string) string {
	_ = fmt.Sprintf // hint: use this
	panic("implement me")
}
```

- [ ] **Step 8.3: Verify tests fail with panic**

```bash
go test -v ./internal/challenges/basics/08-account-greeter/...
```
Expected: panic "implement me"

- [ ] **Step 8.4: Verify tests pass**

Implement both functions, run, revert.
Expected: `PASS`

- [ ] **Step 8.5: Create README**

```markdown
# 🔍 Case 08: The Account Greeter

## The Detective Brief
Build the greeting service — the simplest package CBA ships. It must be correct, exported correctly, and testable. One exported function. One unexported helper. No excuses.

## The Crime Scene
`go test -v ./internal/challenges/basics/08-account-greeter/...`

## Your Mission
1. Run the tests — they panic (implement me).
2. Implement `formatName` (unexported helper) in `greeting.go`.
3. Implement `Greet` (exported function) using `formatName` and `fmt.Sprintf`.
4. Run the tests again.

**Note:** The editable file here is `greeting.go`, not `challenge.go`. This mirrors how real Go packages name their files after what they contain.

## Key Lesson
> **Python vs Go:** Python uses `_name` convention for private members — it's advisory, not enforced. In Go, a lowercase first letter (`formatName`) means package-private and the compiler enforces it. An uppercase first letter (`Greet`) means exported and callable from any package. There are no other access modifiers.

<details>
<summary>Hints (click to reveal)</summary>

1. `formatName(first, last string) string` — return `first + " " + last`.
2. `Greet(accountID, first, last string) string` — use `fmt.Sprintf("Hello, %s! Your account %s is ready.", formatName(first, last), accountID)`.
3. Try calling `formatName` from a different package — the compiler will reject it.
4. The test file is in the same package (`package accountgreeter`) so it *can* test unexported behaviour indirectly.
</details>
```

- [ ] **Step 8.6: Commit**

```bash
git add internal/challenges/basics/08-account-greeter/
git commit -m "feat(challenges): add 08-account-greeter — exported vs unexported implme"
```

---

## Task 9: Challenge 09 — The Interest Bug (testme — BONUS)

**Concept:** Table-driven tests. `t.Helper()`. Finding bugs by writing tests.

**Files:**
- Create: `internal/challenges/basics/09-interest-bug/interest.go`
- Create: `internal/challenges/basics/09-interest-bug/interest_test.go`
- Create: `internal/challenges/basics/09-interest-bug/README.md`

- [ ] **Step 9.1: Create interest.go with a hidden bug**

```go
// internal/challenges/basics/09-interest-bug/interest.go
package interestbug

import "errors"

// ErrNegativeRate is returned when the interest rate is negative.
var ErrNegativeRate = errors.New("interest rate cannot be negative")

// Calculate computes compound interest.
//
//	principal: starting amount (must be > 0)
//	rate:      annual interest rate as a decimal (e.g. 0.05 for 5%)
//	years:     number of years
//
// Returns the final amount after compound interest, or an error.
//
// BUG: negative rate is not validated — Calculate returns a result instead of ErrNegativeRate.
// (Do not fix this — your job is to write tests that CATCH this bug.)
func Calculate(principal int64, rate float64, years int) (int64, error) {
	if years < 0 {
		return 0, errors.New("years cannot be negative")
	}
	// BUG: missing: if rate < 0 { return 0, ErrNegativeRate }
	result := float64(principal)
	for i := 0; i < years; i++ {
		result *= 1 + rate
	}
	return int64(result), nil
}
```

- [ ] **Step 9.2: Create interest_test.go stub for students to fill in**

```go
// internal/challenges/basics/09-interest-bug/interest_test.go
package interestbug

import (
	"errors"
	"math"
	"testing"
)

// checkResult is a pre-built test helper. Use it in your table cases.
// t.Helper() ensures that when a test fails, the error points to the
// table row that failed — not to this function.
func checkResult(t *testing.T, got, want int64) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

// TestCalculate exercises Calculate with multiple scenarios.
// Your job: add at least 5 test cases. Include the case that exposes the hidden bug.
//
// The bug: Calculate does not validate negative rates — it should return ErrNegativeRate.
// One of your test cases must catch this bug and fail until interest.go is fixed.
func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		principal int64
		rate      float64
		years     int
		want      int64
		wantErr   error
	}{
		// TODO: add at least 5 test cases.
		// Include:
		//   - a normal case (e.g. 100000 at 5% for 1 year = 105000)
		//   - zero years (result == principal)
		//   - zero rate (result == principal)
		//   - negative years (want error)
		//   - negative rate (want ErrNegativeRate — this one will fail until the bug is fixed!)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.principal, tt.rate, tt.years)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			checkResult(t, got, tt.want)
		})
	}
}
```

- [ ] **Step 9.3: Verify the stub compiles but has no actual test cases yet**

```bash
go test -v ./internal/challenges/basics/09-interest-bug/...
```
Expected: `ok` with 0 tests run (empty table) — participants must add cases.

- [ ] **Step 9.4: Verify the negative-rate case will fail when added**

Add the negative-rate case temporarily to confirm it fails:
```go
{"negative rate", 1000, -0.05, 1, 0, ErrNegativeRate},
```
Run — expect `FAIL`. This confirms the bug is real and the test case works. Remove the temporary case and leave the stub for students.

- [ ] **Step 9.5: Create README**

```markdown
# 🔍 Case 09: The Interest Bug ⭐ BONUS

## The Detective Brief
The interest calculator passed code review. The team wants to ship it. Your job: write the tests that prove it's broken before it hits production accounts.

## The Crime Scene
`go test -v ./internal/challenges/basics/09-interest-bug/...`

## Your Mission
1. Read `interest.go` — understand what `Calculate` should do.
2. Add at least **5 table cases** to `TestCalculate` in `interest_test.go`.
3. One of your cases must expose the hidden bug — watch it fail.
4. (Optional) Fix the bug in `interest.go` and watch your test go green.

**Note:** For this challenge you edit `interest_test.go`, not `challenge.go`.

## Key Lesson
> **Go testing:** Table-driven tests force you to enumerate edge cases you'd miss with ad-hoc tests. `t.Helper()` is already called in `checkResult` — notice how failure messages point to the table row, not the helper function.

<details>
<summary>Hints (click to reveal)</summary>

1. A normal case: principal=1000, rate=0.05, years=1 → want=1050.0
2. Zero years: any principal and rate → want=principal (no compounding)
3. Zero rate: any principal → want=principal (1+0 = 1, no growth)
4. Negative years: `want error` (already validated in `Calculate`)
5. **The bug case**: negative rate (e.g. -0.05) should return `ErrNegativeRate` — but it doesn't. Add this case and watch it fail.
</details>
```

- [ ] **Step 9.6: Commit**

```bash
git add internal/challenges/basics/09-interest-bug/
git commit -m "feat(challenges): add 09-interest-bug — table-driven testme bonus"
```

---

## Task 10: Challenge 10 — The Transfer Error (implme — BONUS)

**Concept:** Custom error types. `errors.As`. `errors.Is`.

**Files:**
- Create: `internal/challenges/basics/10-transfer-error/challenge_test.go`
- Create: `internal/challenges/basics/10-transfer-error/challenge.go`
- Create: `internal/challenges/basics/10-transfer-error/README.md`

- [ ] **Step 10.1: Create the test file**

```go
// internal/challenges/basics/10-transfer-error/challenge_test.go
package transfererror

import (
	"errors"
	"testing"
)

func TestTransfer_Success(t *testing.T) {
	err := Transfer("acc-001", "acc-002", 10000)
	if err != nil {
		t.Fatalf("expected no error for valid transfer, got: %v", err)
	}
}

func TestTransfer_ZeroAmount_ReturnsTransferError(t *testing.T) {
	err := Transfer("acc-001", "acc-002", 0)

	var te *TransferError
	if !errors.As(err, &te) {
		t.Fatalf(
			"expected *TransferError, got: %T %v\n"+
				"  Hint: Transfer must return *TransferError for validation failures.",
			err, err,
		)
	}
	if te.FromID != "acc-001" {
		t.Errorf("TransferError.FromID = %q, want %q", te.FromID, "acc-001")
	}
	if te.Amount != 0 {
		t.Errorf("TransferError.Amount = %d, want 0", te.Amount)
	}
}

func TestTransfer_NegativeAmount_ReturnsTransferError(t *testing.T) {
	err := Transfer("acc-X", "acc-Y", -5000)

	var te *TransferError
	if !errors.As(err, &te) {
		t.Fatalf("expected *TransferError for negative amount, got: %T %v", err, err)
	}
	if te.Amount != -5000 {
		t.Errorf("TransferError.Amount = %d, want -5000", te.Amount)
	}
}

func TestTransfer_SameAccount_ReturnsTransferError(t *testing.T) {
	err := Transfer("acc-001", "acc-001", 10000)

	var te *TransferError
	if !errors.As(err, &te) {
		t.Fatalf("expected *TransferError for same-account transfer, got: %T %v", err, err)
	}
	if te.Reason == "" {
		t.Error("TransferError.Reason should describe why the transfer failed")
	}
}

func TestTransferError_ErrorString(t *testing.T) {
	te := &TransferError{
		FromID: "acc-A",
		ToID:   "acc-B",
		Amount: 50000,
		Reason: "amount must be positive",
	}
	msg := te.Error()
	if msg == "" {
		t.Fatal("TransferError.Error() must return a non-empty string")
	}
}
```

- [ ] **Step 10.2: Create challenge.go with stubs**

```go
// internal/challenges/basics/10-transfer-error/challenge.go
package transfererror

import "fmt"

// TransferError carries structured context about a failed transfer.
// Use errors.As to extract this from a returned error.
type TransferError struct {
	FromID string
	ToID   string
	Amount int64
	Reason string
}

// Error implements the error interface.
// TODO: return a descriptive string, e.g.:
//   "transfer from acc-001 to acc-002 of 10000 failed: <reason>"
func (e *TransferError) Error() string {
	_ = fmt.Sprintf // hint: use this
	panic("implement me")
}

// Transfer validates and records a fund transfer.
//
// Rules:
//   - amount must be > 0, otherwise return *TransferError with Reason "amount must be positive"
//   - fromID and toID must differ, otherwise return *TransferError with Reason "cannot transfer to same account"
//   - on success, return nil
//
// TODO: implement the validation logic.
func Transfer(fromID, toID string, amount int64) error {
	panic("implement me")
}
```

- [ ] **Step 10.3: Verify tests fail with panic**

```bash
go test -v ./internal/challenges/basics/10-transfer-error/...
```
Expected: panic "implement me"

- [ ] **Step 10.4: Verify tests pass with implementation**

Implement `Error()` and `Transfer()`, run, revert.
Expected: `PASS`

- [ ] **Step 10.5: Create README**

```markdown
# 🔍 Case 10: The Transfer Error ⭐ BONUS

## The Detective Brief
A transfer failed. The log says "error". Which account? How much? What rule? The compliance team needs answers, not just a string that says "error".

## The Crime Scene
`go test -v ./internal/challenges/basics/10-transfer-error/...`

## Your Mission
1. Run the tests — they panic (implement me).
2. Implement `TransferError.Error()` — a descriptive string.
3. Implement `Transfer()` — validate amount and account IDs, return `*TransferError` for failures.
4. Run the tests again.

## Key Lesson
> **Python vs Go:** Python exceptions carry arbitrary attributes. Go errors are values — define a struct that implements the `error` interface to carry structured context. Use `errors.As(err, &te)` to extract the concrete type, not type assertions (`err.(*TransferError)`). `errors.As` works through wrapped errors; direct type assertions don't.

<details>
<summary>Hints (click to reveal)</summary>

1. `TransferError.Error()` — use `fmt.Sprintf("transfer from %s to %s of $%.2f failed: %s", ...)`.
2. `Transfer()` — check `amount <= 0` first, then `fromID == toID`.
3. Return `&TransferError{FromID: fromID, ToID: toID, Amount: amount, Reason: "..."}` for each failure.
4. Return `nil` on success.
5. The test uses `errors.As(err, &te)` — this works even if the error is wrapped with `fmt.Errorf("...: %w", te)`.
</details>
```

- [ ] **Step 10.6: Commit**

```bash
git add internal/challenges/basics/10-transfer-error/
git commit -m "feat(challenges): add 10-transfer-error — custom error types implme bonus"
```

---

## Task 11: Clean Up Old Challenges and Update Index

- [ ] **Step 11.1: Remove superseded challenge directories (tracked by git)**

Use `git rm -rf` so deletions are staged and included in the next commit. All nine old directories are superseded — challenges 06–09 (package-refactoring, mock-payment-gateway, middleware-chain, worker-pool) covered Day 2/advanced topics and do not belong in the Day 1 basics set.

```bash
git rm -rf internal/challenges/basics/01-structs-and-pointers
git rm -rf internal/challenges/basics/02-interfaces-and-receivers
git rm -rf internal/challenges/basics/03-error-handling
git rm -rf internal/challenges/basics/04-table-driven-tests
git rm -rf internal/challenges/basics/05-safe-http-client
git rm -rf internal/challenges/basics/06-package-refactoring
git rm -rf internal/challenges/basics/07-mock-payment-gateway
git rm -rf internal/challenges/basics/08-middleware-chain
git rm -rf internal/challenges/basics/09-worker-pool
```

- [ ] **Step 11.2: Update the basics README index**

Update `internal/challenges/basics/README.md` (or `internal/challenges/README.md`) to list all 10 new challenges with their tracks and difficulty.

```markdown
# Day 1 Challenges: Detective Briefs

These are the **Go building blocks challenges** for Day 1. Each challenge is a short detective mystery — run the tests, they fail, find the bug or implement the function, make them pass.

## Track 1: Go Building Blocks (13:30)
| # | Challenge | Type | Concept |
|---|-----------|------|---------|
| 01 | [The Frozen Account](01-frozen-account/) | fixme | Value vs pointer receiver |
| 02 | [The Dead Map](02-dead-map/) | fixme | Nil map initialization |
| 03 | [The Phantom Append](03-phantom-append/) | fixme | Append semantics + value receiver |
| 04 | [Silent Errors](04-silent-errors/) | implme | Sentinel errors, errors.Is |

## Track 2: Interfaces, Receivers & Context (14:30)
| # | Challenge | Type | Concept |
|---|-----------|------|---------|
| 05 | [The Fee Calculator](05-fee-calculator/) | implme | Interface + pointer receiver method set |
| 06 | [The Interface Lie](06-interface-lie/) | fixme | Typed nil in interface |
| 07 | [The Ignored Cancel](07-ignored-cancel/) | implme | Context cancellation, goroutines, defer |

## Track 3: Build, Package & Run (15:20)
| # | Challenge | Type | Concept |
|---|-----------|------|---------|
| 08 | [The Account Greeter](08-account-greeter/) | implme | Exported vs unexported |

## Bonus (any time)
| # | Challenge | Type | Concept |
|---|-----------|------|---------|
| 09 | [The Interest Bug](09-interest-bug/) ⭐ | testme | Table-driven tests, t.Helper |
| 10 | [The Transfer Error](10-transfer-error/) ⭐ | implme | Custom error types, errors.As |

## Running All Challenges
```bash
go test -v ./internal/challenges/basics/...
```

## Running a Single Challenge
```bash
go test -v ./internal/challenges/basics/01-frozen-account/...
```
```

- [ ] **Step 11.3: Commit the cleanup and index update**

Deletions are already staged by `git rm` in Step 11.1. Just add the README changes:

```bash
git add internal/challenges/basics/README.md
git commit -m "chore(challenges): replace old 01-09 challenges with Day 1 detective brief set"
```

---

## Task 12: Final Verification

- [ ] **Step 12.1: Run all new challenges — verify each fails in the expected way**

```bash
# 01–03: fixme — tests should FAIL with descriptive assertion messages
go test -v ./internal/challenges/basics/01-frozen-account/...
# Expected: FAIL "Balance should be 100 after Deposit(100) — got 0.00"

go test -v ./internal/challenges/basics/02-dead-map/...
# Expected: FAIL "nil map panic — initialize the map first..."

go test -v ./internal/challenges/basics/03-phantom-append/...
# Expected: FAIL "expected 1 suspect after AddSuspect, got 0"

# 04, 05, 08, 10: implme — tests should panic with "implement me"
go test -v ./internal/challenges/basics/04-silent-errors/...
go test -v ./internal/challenges/basics/05-fee-calculator/...
go test -v ./internal/challenges/basics/08-account-greeter/...
go test -v ./internal/challenges/basics/10-transfer-error/...
# Expected: each panics with "implement me"

# 06: fixme — tests should FAIL with typed-nil explanation
go test -v ./internal/challenges/basics/06-interface-lie/...
# Expected: FAIL "runAMLCheck(skip=true) should return nil error, got: <nil> (*interfacelie.AMLChecker)"

# 07: implme — goroutine panics; use -race and timeout
go test -v -race -timeout 10s ./internal/challenges/basics/07-ignored-cancel/...
# Expected: goroutine panics with "implement me: call report on each iteration"

# 09: testme EXCEPTION — empty table runs 0 tests, reports ok (this is expected before student adds cases)
go test -v ./internal/challenges/basics/09-interest-bug/...
# Expected: ok (0 tests run) — this is CORRECT. Students add table cases; the negative-rate case then fails.
```

- [ ] **Step 12.2: Verify the module builds cleanly**

```bash
go build ./...
```
Expected: no errors (`panic()` stubs compile cleanly; unused import guards like `_ = fmt.Sprintf` prevent "imported and not used" errors)

- [ ] **Step 12.3: Run Challenge 07 with race detector explicitly**

```bash
go test -v -race -timeout 15s ./internal/challenges/basics/07-ignored-cancel/...
```
Expected: panic with "implement me" (no race reported — the stub panics before any race can occur). After the student implements it, re-run with `-race` to confirm no data races.

- [ ] **Step 12.4: Final commit**

```bash
git add docs/superpowers/specs/ docs/superpowers/plans/
git commit -m "docs: add Day 1 challenge spec and implementation plan"
```
