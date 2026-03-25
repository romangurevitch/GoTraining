# Challenges

This directory contains all student exercises for the Go Training workshop.

## Structure

```
challenges/
├── basics/       # Core Go language fundamentals
├── bank/         # Go Bank HTTP & Data layer quests
└── temporal/     # Durable Workflow & Agentic Go quest
```

## basics/

The **[Day 1 Challenges: Detective Briefs](basics/README.md)** are short exercises covering core Go building blocks, interfaces, concurrency, and testing.

Each challenge is a short mystery. You will encounter different types of quests:
- **fixme**: Buggy code is provided. Your task is to identify the problem and fix it.
- **implme**: You'll find `panic("implement me!")` stubs. Your task is to implement the function to make the tests pass.
- **testme**: You'll need to write tests to verify existing logic.

Run tests with: `make test-challenges`

## bank/

The **[Go Bank Transfer Quest](bank/README.md)** is your introduction to building production API handlers. 

You'll implement the `POST /v1/transfers` API endpoint in a pre-scaffolded service, focusing on:
- Idiomatic HTTP handler patterns using Gin.
- OpenTelemetry tracing and structured logging with `slog`.
- JWT authentication and scope-based authorisation.
- Table-driven unit testing for handlers.

## temporal/

The **[Durable Transfer Quest](temporal/README.md)** is the final, high-stakes challenge! 

You will transform the bank transfer into a robust **Distributed Transaction** using Temporal. This module focuses on:
- **Agentic Engineering:** Using specialized AI tools to build complex logic.
- **Workflow Orchestration:** Implementing the Compensation Pattern and Human-in-the-loop approvals.
- **Durable Reliability:** Ensuring idempotency and surviving worker restarts.
- **Full-Lifecycle Testing:** Unit, Integration, and the "Gold Standard" Replay tests.

Evaluate your work using the **[Competition Grading Rubric](temporal/GRADING_PROMPT.md)**.

## Your Exploration Journey

Ready to put your skills to the test? We've prepared a series of challenges that increase in complexity as you progress through the training.

Start with the basics: **[Day 1 Challenges: Detective Briefs](basics/README.md)**.

---
[← Back to Main README](../../README.md)
