# 📅 Go Training Schedule

> **Important:** This schedule is a **flexible guide**. The timing and depth of each session will be adjusted dynamically based on participant experience levels, the pace of challenge completion, and group interests.

## Day 1: API Architecture and the Go Foundation

Focus: Transitioning from interpretative languages to Go's type system and building production-grade API contracts.

| Time | Topic | Description | Module Link |
| :--- | :--- | :--- | :--- |
| 09:30 | **API Fundamentals & Design** | Core principles of REST vs. RPC, security considerations for modern services, and adopting structured logging with `slog`. | [Module 1](README.md#module-1-modern-api-engineering-principles) |
| 10:30 | **Lifecycle & The Agentic Future** | API evolution, graceful sunsetting strategies, containerisation for AWS (ECS), and the role of MCP in AI-driven development. | [Module 1](README.md#module-1-modern-api-engineering-principles) |
| 11:20 | **Toolchain & Docker** | Mastering Go dependency management (`go mod`), and building efficient, multi-stage Docker images. | [Toolchain & Docker](internal/basics/toolchain/README.md) |
| 12:00 | *Lunch Break* | | |
| 13:00 | **API Trivia** | *Sequence Breaker:* A fun exploration of the "Bad API Hall of Fame" and design anti-patterns. | |
| 13:30 | **Go Building Blocks I** | The mental shift: Memory management with pointers, data structures (structs/slices), and Go's unique error-as-value philosophy. | [Module 2](README.md#module-2-go-language-fundamentals) |
| 14:30 | **Go Building Blocks II** | Advanced foundations: Implicit interface implementation, value vs. pointer receivers, and concurrency control with `context`. | [Module 2](README.md#module-2-go-language-fundamentals) |
| 15:20 | **Hands-on Lab** | **Track 1 & 2:** Solving the "Detective Brief" challenges. Moving from "implement me" to fixing subtle runtime bugs. | [Challenges](README.md#challenges) |

## Day 2: Implementation, Orchestration and AI

Focus: Building scalable services, managing durable state with Temporal, and leveraging AI for engineering efficiency.

| Time | Topic | Description | Module Link |
| :--- | :--- | :--- | :--- |
| 09:30 | **The HTTP & Data Layer** | Building high-performance services with Gin, implementing middleware, and the adapter pattern for clean storage isolation. | [Module 3](README.md#module-3-building-the-data--api-service) |
| 10:30 | **Codebase Walkthrough** | Investigating the "Go Bank" architecture: a deep dive into idiomatic project layout and cross-package dependencies. | [Module 3](README.md#module-3-building-the-data--api-service) |
| 11:00 | **Transfer Funds Quest** | **Major Challenge:** Implementing a production-grade fund transfer endpoint end-to-end (OpenAPI -> Handler -> Service). | [Challenges](README.md#challenges) |
| 12:00 | *Lunch Break* | | |
| 13:00 | **Building Blocks Trivia** | *Sequence Breaker:* A fast-paced quiz on Go syntax, "gotchas," and runtime behaviour. | |
| 13:30 | **Temporal Orchestration** | Reliability at scale: Workflow concepts, signal handling, and the replay model for durable execution. | [Module 4](README.md#module-4-temporal-orchestration) |
| 14:30 | **Testing & Service Isolation** | Robust verification: Unit testing with Mockery, integration testing, and an introduction to Infrastructure as Code. | [Module 2](README.md#module-2-go-language-fundamentals) |
| 15:15 | **Agentic Go & Durable Workflows** | **Hands-On Challenge:** Each participant builds a high-value transfer workflow with human-in-the-loop approval using their AI agent. Work through the quests in order — get as far as you can. Full quest guide remains available for self-study after the session. | [Challenges](internal/challenges/temporal/README.md) |
| 16:00 | **Wrap up & Feedback** | Final open Q&A session and programme evaluation. | |
