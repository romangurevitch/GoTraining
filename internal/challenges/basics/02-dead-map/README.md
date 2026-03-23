# 🔍 Case 02: The Dead Map

## The Detective Brief
The transaction logger runs — until it doesn't. The map was declared. Was it ever born? The logger panics at runtime, but looked fine in code review.

## The Crime Scene
```bash
go test -v ./internal/challenges/basics/02-dead-map/...
```

## Your Mission
1. Run the tests — they panic with a friendly message.
2. Read `NewLedger()` in `challenge.go`.
3. Initialize the map.
4. Run the tests again — they should pass.

## Key Lesson
> **Python vs Go:** In Python, `{}` gives you a ready-to-use dict. In Go, `map[string]float64` declares the type but the value is `nil` — you must call `make(map[string]float64)` or use a composite literal `map[string]float64{}` before writing to it.

<details>
<summary>Hints (click to reveal)</summary>

1. Look at the `NewLedger()` function — what is `entries` set to?
2. `var m map[string]float64` declares a nil map. Reads return zero values silently. Writes panic.
3. Fix: `entries: make(map[string]float64)` inside the struct literal in `NewLedger`.
4. Alternatively: `entries: map[string]float64{}` works too.
</details>
