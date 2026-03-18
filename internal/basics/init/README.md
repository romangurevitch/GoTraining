# The init() Function

`init()` runs automatically when a package is loaded, before `main()`.

## Execution Order

1. Package-level variables initialised
2. `init()` functions run (in file order)
3. `main()` runs

## Pitfalls

- `init()` cannot be called explicitly
- Side effects in `init()` (DB connections, global state) make testing harder
