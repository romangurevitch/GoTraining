# 🔍 Case 03: The Phantom Append

## The Detective Brief
The fraud team adds three suspects to the watch list. The list is always empty. The append ran. Nothing stuck.

## The Crime Scene
```bash
cd internal/challenges/basics/03-phantom-append
```
```bash
go test -v ./...
```

## Your Mission
1. Run the tests — they fail.
2. Read `AddSuspect` in `challenge.go` — there are **two** bugs.
3. Fix both.
4. Run the tests again.

## Key Lesson
> **Python vs Go:** Python's `list.append()` mutates the list in place. In Go, `append()` may allocate a new backing array and always returns the updated slice — which you must capture. When combined with a **value receiver** (which copies the struct), mutations can vanish in two different ways.

<details>
<summary>Hints (click to reveal)</summary>

1. Bug 1: What type of receiver does `AddSuspect` use? Recall what we learned in Case 01 about value vs. pointer receivers.
2. Bug 2: In Go, `append` returns a new slice header. If you don't assign this return value back to your slice, what happens to the changes?
3. How can you ensure that both the struct mutation and the slice update are preserved?
</details>
