# **Go Best Practices: BDD & E2E Implementation Plan**

This document defines the architecture and implementation roadmap for a Go reference project. The goal is to demonstrate a "learned" systems approach using Gherkin-style BDD and Testcontainers for Docker Compose orchestration.

## **1\. Project Philosophy**

* **Verifiable Architecture:** Code is only as good as its tests. High-level E2E tests prove the system works as a "black box."  
* **Consumer-Driven Interfaces:** Interfaces are defined by the consumer to enable easy mocking and decoupling of legacy systems.  
* **Hermetic Testing:** No reliance on shared databases. Every test suite spins up its own ephemeral infrastructure.  
* **Boring Code:** Prioritize readability and explicit dependency injection over "magic" frameworks or clever optimizations.

## **2\. Directory Structure (Standard Go Layout)**

.  
├── cmd/  
│   └── server/  
│       └── main.go           \# Application entry point  
├── internal/  
│   └── platform/             \# Infrastructure (DB, HTTP clients)  
│   └── user/                 \# Domain logic (Handlers, Services)  
├── features/                 \# Gherkin .feature files  
│   └── user\_registration.feature  
├── test/  
│   ├── bdd\_test.go           \# Godog suite runner & Step Definitions  
│   └── testdata/  
│       └── docker-compose.yml \# Ephemeral test infrastructure  
├── Makefile                  \# Build and test automation  
└── go.mod

## **3\. Core Dependencies**

| Library | Purpose |
| :---- | :---- |
| github.com/cucumber/godog | Gherkin BDD framework for Go. |
| github.com/testcontainers/testcontainers-go | Programmatic Docker lifecycle management. |
| github.com/testcontainers/testcontainers-go/modules/compose | Docker Compose orchestration within Go tests. |
| github.com/stretchr/testify | Assertions and mocking toolkit. |

## **4\. Implementation Steps**

### **Phase 1: Infrastructure Setup**

1. **Define the Compose File:** Create test/testdata/docker-compose.yml defining the application and its dependencies (e.g., PostgreSQL).  
2. **Build the Harness:** Implement the Testcontainers "Compose Module" in test/bdd\_test.go to handle stack.Up() and stack.Down().  
3. **Wait Strategies:** Configure health checks in YAML to ensure the Go test runner waits until the database is ready for connections.

### **Phase 2: BDD Definition**

1. **Write the Feature:** Draft features/user\_registration.feature using Given/When/Then syntax.  
2. **Generate Snippets:** Run godog features/ to generate the initial Go function stubs for step definitions.  
3. **Implement Steps:** Fill in the logic to:  
   * Perform HTTP requests to the dynamically mapped port of the container.  
   * Query the database container to verify state persistence.

### **Phase 3: Application Logic**

1. **The Handler:** Implement a standard net/http handler in internal/user.  
2. **Dependency Injection:** Ensure the database connection is passed into the handler via a constructor (NewHandler(db \*sql.DB)).  
3. **Health Check:** Provide a /health endpoint for the Docker health check to verify readiness.

## **5\. The Test Execution Flow**

1. **Initialization:** go test triggers the TestMain or TestFeatures function.  
2. **Orchestration:** Testcontainers pulls images and executes docker compose up.  
3. **Discovery:** The test queries the Docker API to find which random port the host assigned to the app's 8080\.  
4. **Execution:** Godog parses the Gherkin file and executes the Go steps against the live containers.  
5. **Teardown:** t.Cleanup ensures docker compose down \-v is called, removing all containers and volumes.

## **6\. Comparison: Performance vs. Reliability**

Unlike the "Primeagen style" which focuses on micro-optimizing CPU cycles and Vim speed, this plan prioritizes **System Reliability**.

* **The "Hype" Way:** Manual testing or trusting that "Vim \+ Go" is enough.  
* **The "Learned" Way:** Automating the proof that the code works in a real network environment before a single line is merged.

## **7\. Success Criteria**

* \[ \] go test ./test/... passes with zero manual intervention.  
* \[ \] No local Docker containers are left running after the test suite finishes.  
* \[ \] The .feature file provides a clear, human-readable documentation of the system's behavior.