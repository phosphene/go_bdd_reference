# Go BDD Reference: A "Learned" Systems Approach

This repository serves as a reference implementation for modern, production-grade Go services. It prioritizes **Reliability**, **Verifiability**, and **Maintainability** through a strict Gherkin-style BDD framework and hermetic orchestration.

---

## **The "What"**
A minimal, decoupled Go service demonstrating:
*   **Hexagonal Lite Architecture**: Decoupling domain logic from delivery (HTTP) and infrastructure (Docker).
*   **Hermetic E2E Testing**: Using [Testcontainers](https://testcontainers.com/) to programmatically orchestrate ephemeral Docker Compose stacks for every test run.
*   **Gherkin BDD**: Human-readable behavioral specifications using [Godog](https://github.com/cucumber/godog).
*   **Modern Go (1.24+)**: Leveraging the latest stable idioms, iterators, and strict linting.

---

## **The "Why"**
Traditional testing often falls into two traps:
1.  **Mock-Heavy Unit Tests**: They prove the code does what the developer *thinks* it does, but often miss integration failures (network issues, SQL syntax, container misconfiguration).
2.  **Brittle Manual QA**: Slow, non-deterministic, and impossible to scale.

### **Our Testing Philosophy**
We follow a **Tracer Bullet** approach:
*   **Verify the System, Not the Function**: High-level BDD tests prove the system works as a "black box" in a real network environment.
*   **Zero-Drift Infrastructure**: If a test requires a database (e.g., Postgres, Dolt), it must be defined in the test's `docker-compose.yml`. No reliance on external "staging" environments.
*   **Consumer-Driven Interfaces**: We "discover" interfaces when we need to decouple, rather than pre-abstracting everything.

---

## **Best Practices & Standards**

Follow the [Go Modern Standards Guide](file:///Users/edphillips/projects/new/go_bdd_reference/shared/GoModernStandardsGuide.md) in the `shared/` directory for detailed rules on:
1.  **Syntax**: Prefer `for range n` over verbose loops; always wrap errors with context (`fmt.Errorf("...: %w", err)`).
2.  **Design**: Keep templates "dumb," keep constructors explicit (no `init()` logic), and use Functional Options for complex setup.
3.  **Testing**: Use the AAA (Arrange, Act, Assert) pattern and table-driven tests for internal logic.

---

## **Getting Started**

### **Prerequisites**
*   **Go 1.24+**
*   **Docker Daemon** (OrbStack, Docker Desktop, or Colima)

### **Execution**
Run the entire suite (unit + BDD) with a single command:
```bash
make test
```

### **Test Coverage**
Generate a line-by-line coverage report:
```bash
make coverage
```
This will:
1.  Run all tests with the `-coverprofile` flag.
2.  Print a text-based summary to the console.
3.  Generate an interactive HTML report at `coverage.html`.

---

## **Project Layout**
```text
.
├── cmd/server           # Entry point (Keep thin)
├── internal/            # Private logic (Non-importable)
├── features/            # Gherkin .feature specifications
├── test/                # BDD runner and Testcontainers orchestration
├── shared/              # Axiomatic standards and implementation guides
└── Makefile             # Automation for test, lint, and build
```
