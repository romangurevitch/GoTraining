# Module 1: Modern API Engineering Principles

This module covers the foundational concepts for building production-ready APIs and platform tools.

## Topics

- **[API Fundamentals](api-fundamentals/)** — REST vs. RPC, idempotency, statelessness, contract-first vs. code-first design
- **[Security & Observability](security-and-observability/)** — AuthN/AuthZ, policy as code, structured logging with `slog`, distributed tracing
- **[API Lifecycle & Deployment](api-lifecycle-and-deployment/)** — versioning (v1/v2), graceful sunsetting, containerisation, cloud deployment
- **[The Agentic Future](the-agentic-future/)** — designing APIs for AI consumption, Model Context Protocol (MCP), tool discovery

## Resources

- [Google API Design Guide](https://cloud.google.com/apis/design)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Go slog documentation](https://pkg.go.dev/log/slog)

## Putting Theory into Practice

The principles covered in this module are foundational to building production-grade services. You will apply these concepts in the following quests:

- **[Go Bank Transfer Quest](../challenges/bank/README.md)** — Implement an idempotent API handler with tracing and security.
- **[Durable Transfer Quest](../challenges/temporal/README.md)** — Build a robust distributed transaction that survives failures and integrates human-in-the-loop approvals.

## Bonus
- **[API Design](api-design/)** — What to do and not to do, with real-world examples of poor API design decisions
