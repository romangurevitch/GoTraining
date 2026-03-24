# Day 1 Challenges: Detective Briefs

These are the **Go building blocks challenges** for Day 1. Each challenge is a short detective mystery — run the tests, they fail, find the bug or implement the function, make them pass.

## How It Works

1. `cd` into a challenge directory
2. Run: `go test -v ./...`
3. The tests fail — read the error message and the code
4. Fix the bug or implement the function
5. Run again — green means done

## Track 1: Go Building Blocks

| # | Challenge | Type | Concept |
|---|-----------|------|---------|
| 01 | [The Frozen Account](01-frozen-account/) | fixme | Value vs pointer receiver |
| 02 | [The Dead Map](02-dead-map/) | fixme | Nil map initialisation |
| 03 | [The Phantom Append](03-phantom-append/) | fixme | Append semantics + value receiver |
| 04 | [Silent Errors](04-silent-errors/) | implme | Sentinel errors, errors.Is |

## Track 2: Interfaces, Receivers & Context

| # | Challenge | Type | Concept |
|---|-----------|------|---------|
| 05 | [The Fee Calculator](05-fee-calculator/) | implme | Interface + pointer receiver method set |
| 06 | [The Interface Lie](06-interface-lie/) | fixme | Typed nil in interface |
| 07 | [The Ignored Cancel](07-ignored-cancel/) | implme | Context cancellation, goroutines, defer |

## Track 3: Build, Package & Run

| # | Challenge | Type | Concept |
|---|-----------|------|---------|
| 08 | [The Account Greeter](08-account-greeter/) | implme | Exported vs unexported identifiers |

## Bonus (any time — for fast finishers)

| # | Challenge | Type | Concept |
|---|-----------|------|---------|
| 09 | [The Interest Bug](09-interest-bug/) ⭐ | testme | Table-driven tests, t.Helper |
| 10 | [The Transfer Error](10-transfer-error/) ⭐ | implme | Custom error types, errors.As |

## Running a Single Challenge

```bash
cd internal/challenges/basics/01-frozen-account
```
```bash
go test -v ./...
```

## Running All Challenges

```bash
cd internal/challenges/basics
```
```bash
go test -v ./...
```
