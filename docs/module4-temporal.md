# Module 4: Temporal Orchestration

**Duration:** Day 2 afternoon (13:30–14:15, demo only)
**Location:** `internal/temporal/`

## Overview

Module 4 is a demo-only session. The instructor shows a working TransferFunds saga using Temporal. Participants receive reference material to explore independently.

## Key Question

*Why not just use goroutines?*

Raw goroutines work for short-lived concurrent work, but fail for long-running operations:
- Process crashes lose in-progress state
- Retries must be coded manually
- Timeout and compensation logic is error-prone

Temporal provides **durable execution** — workflows survive failures and automatically retry.

## Self-Paced Reference

See `internal/temporal/README.md` for concepts, key terms, and learning resources.
