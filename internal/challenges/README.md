# Challenges

This directory contains all student exercises for the Go Training workshop.

## Structure

```
challenges/
├── basics/
│   ├── fixme/    # Find and fix the bugs
│   └── implme/   # Implement the function
└── bank/         # Go Bank service quests
```

## basics/fixme

Short, focused exercises where buggy code is provided. Your task: identify the problem and fix it.
Inspired by the [ConcurrencyWorkshop](https://github.com/) fixme pattern.

## basics/implme

Exercises with `panic("implement me!")` stubs. Your task: implement the function to make the tests pass.

## bank/

The main challenge — implement the full Go Bank service layer:
1. Domain models (`internal/bank/domain/`)
2. Store / repository (`internal/bank/store/`)
3. Service / business logic (`internal/bank/service/`)
4. HTTP API (`internal/bank/api/`)

Run tests with: `make test-bank`
