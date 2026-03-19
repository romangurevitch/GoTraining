# Repository Review & Inconsistencies

This document outlines the inconsistencies, bugs, and areas for improvement identified during the review of the Go Training repository.

## 1. High-Level Inconsistencies

### Missing `internal/basics/README.md` - not complete yet - don't do anything about this. 
- **Issue:** The root `README.md` links to `internal/basics`, but there is no `README.md` in that directory.
- **Impact:** Breaks the navigation flow for students starting Module 2.

### `make test-challenges` Fails - not complete yet - don't do anything about this. 
- **Issue:** Running `make test-challenges` results in `no packages to test`.
- **Impact:** The `internal/challenges/` directory contains only READMEs and no Go packages, making the command useless in its current state.

## 2. Infrastructure & Tooling

### `make db-up` Issues - Fix this issue use docker-compose and add a comment on the readme in makefile about common issues. 
- **Issue:** The command uses `docker compose`, which might be `docker-compose` on some systems. Additionally, port 5432 conflicts are common but not mentioned in the troubleshooting/setup guide.
- **Impact:** First-time users might be blocked at the very beginning of the workshop.

### Non-Portable `tools.mk` - it should download the one that should work for the user but most users will use macs so not a big deal if difficult or complext to solve.
- **Issue:** The `protoc` download target explicitly assumes `osx-aarch_64` (Apple Silicon).
- **Impact:** Installation of tools will fail on Linux or Intel-based Macs.

### `make migrate` is a Stub - look at the training reference project, do we need this at all? We need to populate the db only for the bank challenge. 
- **Issue:** The command only echoes instructions and does not perform actual migrations.
- **Impact:** Inconsistent with the "Getting Started" guide which implies a ready-to-use environment.

## 3. Architecture & Design

### Lack of Atomicity in `Transfer` - we can point that out in the training but there's no need to create complexities for students, the priority is to learn go.
- **Issue:** `BankService.Transfer` in `internal/bank/service/service.go` performs multiple independent repository calls without a database transaction.
- **Impact:** Inconsistent with "production-grade" training goals. If one call fails after the debit, funds are lost.

### Repository Transaction Support - again no need for that
- **Issue:** The `Repository` interface in `internal/bank/repository/repository.go` lacks methods for transaction management (`Begin`, `Commit`, `Rollback`).
- **Impact:** Makes it impossible to implement atomic operations correctly within the service layer.

### Package Naming Conflict - suggest a better go way of doing this.
- **Issue:** The package `pkg/api/error` is named `error`, which is a built-in type in Go.
- **Impact:** While legal (especially with aliasing), it is considered bad practice and can lead to confusion and shadowing.

### Inconsistent Domain Errors - make it consistant
- **Issue:** In `internal/bank/domain/account.go`, `CanPerformTransaction` returns a defined `ErrAccountLocked` for locked accounts but an anonymous `errors.New("account is closed")` for closed ones.
- **Impact:** Prevents consistent error handling via `errors.Is(err, domain.ErrAccountClosed)` in the API layer.

## 4. Code & Documentation Mismatches

### `internal/bank/api/server.go` TODOs - update the readme - they need to reference the account and other pre existing code instead. 
- **Issue:** Quest 2 in `internal/challenges/bank/README.md` instructs to "uncomment" a line for `transferHandler`, but the line does not exist in the source code (it's just a `TODO` comment).
- **Impact:** Confuses students who are looking for a specific line to uncomment.

### `cmd/worker/main.go` Status - it's comming soon, ignore for now.
- **Issue:** The file is a bare stub and is explicitly excluded from the main `make build` command.
- **Impact:** Module 4 (Temporal) feels incomplete or "coming soon" rather than ready.

### Outdated Mocking Documentation - all basic stuff are coming soon. 
- **Issue:** `internal/basics/mocking/README.md` refers to the archived `github.com/golang/mock` and incorrectly lists "not type-safe" as a disadvantage for Mockery. The project has already migrated to `go.uber.org/mock` in `go.mod`.
- **Impact:** Misleads students with outdated library recommendations and incorrect technical comparisons.

### Inconsistent Module Naming - Resolve it it should be Modern API Engineering Principles
- **Issue:** Root `README.md` refers to Module 1 as "Modern API Engineering Principles", while `internal/fundamentals/README.md` calls it "Cloud API Engineering Principles".
- **Impact:** Minor branding/naming inconsistency.

---
*Note: This review was conducted on March 19, 2026. Environment-specific versions (like Go 1.26.1) are confirmed as current.*


## 5. Summary of Recommended Fixes
1. Create `internal/basics/README.md`.
2. Update `tools.mk` to be platform-agnostic (using `uname`).
3. Implement a proper migration tool or script.
4. Add transaction support to the Repository and Service layers.
5. Fix the Quest 2 instructions to match the actual code state.
6. Define `ErrAccountClosed` in the domain package.
7. Update mocking documentation to reflect `go.uber.org/mock` and Mockery's newer features.
