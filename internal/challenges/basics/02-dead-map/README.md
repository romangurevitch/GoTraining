# 🔍 Case 02: The Dead Map

## The Detective Brief
The transaction logger runs — until it doesn't. The map was declared. Was it ever born? The logger panics at runtime, but looked fine in code review.

## The Crime Scene
```bash
cd internal/challenges/basics/02-dead-map
```
```bash
go test -v ./...
```

## Your Mission
1. Run the tests — they panic with a friendly message.
2. Read `NewLedger()` in `challenge.go`.
3. Initialise the map.
4. Run the tests again — they should pass.

## Key Lesson
> **Python vs Go:** In Python, `{}` gives you a ready-to-use dict. In Go, declaring a map only reserves the name and type — the value remains `nil`. You must initialise it using `make()` or a composite literal `{}` before you can write to it. Writing to a `nil` map will cause a runtime panic.

<details>
<summary>Hints (click to reveal)</summary>

1. Look at the `NewLedger()` function — what is `entries` set to?
2. In Go, a declared map variable without an assignment is `nil`.
3. Check the differences between `var m map[string]int64` and `m := make(map[string]int64)`.
4. How can you ensure `entries` is not `nil` before it's used?
</details>
