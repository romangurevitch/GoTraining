# 🔍 Case 05: The Counter

## The Detective Brief
Two counters implement the same `Increment()` method — same logic, same code: `c.Count++`. But one counter actually increments and the other stays at zero. Why?

## The Crime Scene
```bash
cd internal/challenges/basics/05-fee-calculator
```
```bash
go test -v ./...
```

## Your Mission
1. Run the tests — they panic (`implement me`).
2. Implement `Increment` for both `ValueCounter` and `PointerCounter` — just `c.Count++`.
3. Run the tests again. One counter increments. The other doesn't. Why?
4. **Experiment:** Change `PointerCounter`'s receiver from `*PointerCounter` to `PointerCounter`. What happens to the tests?

## Key Lesson
> **Value vs Pointer receiver:** A value receiver operates on a **copy** — mutations are discarded when the method returns. A pointer receiver operates on the **original** — mutations stick. This also affects interface satisfaction: a method on `*T` means only `*T` satisfies the interface, not `T`.
>
> **Python vs Go:** In Python, `self` is always a reference — mutations always affect the original object. In Go, you must choose explicitly between value and pointer receivers, and the wrong choice silently drops your mutations.

<details>
<summary>Hints (click to reveal)</summary>

1. Both methods have the same body: `c.Count++`.
2. `ValueCounter` uses a value receiver — the increment happens on a copy and is discarded.
3. `PointerCounter` uses a pointer receiver — the increment modifies the original struct.
4. The test for `ValueCounter` asserts that `Count` is still `0` — that's the point.
5. `IncrementAll` is already implemented — don't modify it.
</details>

## Your Next Step
Interfaces can be tricky, but sometimes they flat out lie to you. It's time to investigate the case of the nil that isn't nil.
Head over to **[Case 06: The Interface Lie](../06-interface-lie/README.md)**.
