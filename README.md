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

This project strictly adheres to the [Go Modern Standards Guide](file:///Users/edphillips/projects/new/go_bdd_reference/shared/GoModernStandardsGuide.md). Core principles include:

1.  **Explicit Dependency Injection**: Avoid `func init()` and global state. All dependencies (DBs, clients, loggers) are passed explicitly to constructors.
2.  **Error Wrapping & Context**: Never just return a raw error. Use `%w` to wrap errors with context: `fmt.Errorf("user.Register(%s): %w", username, err)`.
3.  **Hermetic Environments**: The system must be able to boot and verify itself entirely within Docker, ensuring a "zero-drift" developer experience.
4.  **Standard Layout**: Implementation details are confined to `internal/` packages to prevent leakage and ensure a clean public API.

---

## **Testing Patterns**

### **Table-Driven Tests (TDT)**
For internal logic and delivery handlers, we use Table-Driven Tests. This pattern centralizes test logic, making it easy to add new cases (edge cases, error states) without duplicating boilerplate.

```go
func TestHealthHandler(t *testing.T) {
    tests := []struct {
        name           string
        method         string
        expectedStatus int
    }{
        {"Valid GET", http.MethodGet, http.StatusOK},
        {"Invalid POST", http.MethodPost, http.StatusMethodNotAllowed},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... Arrange, Act, Assert
        })
    }
}
```

### **The AAA Pattern**
We structure all tests following the **Arrange, Act, Assert** pattern to ensure clarity:
1.  **Arrange**: Set up the SUT (System Under Test) and any mocks/stubs.
2.  **Act**: Execute the primary function or API call.
3.  **Assert**: Verify the results using `testify/assert` or `require`.

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
