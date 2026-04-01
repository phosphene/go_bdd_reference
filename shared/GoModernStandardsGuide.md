# **Go Modern Standards: Syntax, Design, & Testing (2026)**

This guide provides a reference for writing "learned," production-grade Go code. It prioritizes maintainability, verifiability, and performance in that order.

## **1\. Modern Syntax & Idioms (Go 1.22+)**

### **Range Over Integers & Functions**

Go 1.22+ simplified loops and iterators. Avoid the verbose for i := 0; i \< n; i++ when possible.

// Modern loop  
for i := range 5 {  
    fmt.Println(i) // 0, 1, 2, 3, 4  
}

### **Error Wrapping & Inspection**

Always provide context to errors. Use %w for wrapping and errors.Is or errors.As for checking.

if err \!= nil {  
    return fmt.Errorf("user.Register(%s): %w", username, err)  
}

// Checking errors  
if errors.Is(err, sql.ErrNoRows) { ... }

### **Functional Options Pattern**

For complex struct initialization, use functional options instead of a massive constructor or multiple New... functions.

type Server struct { port int }  
type Option func(\*Server)

func WithPort(p int) Option {  
    return func(s \*Server) { s.port \= p }  
}

func NewServer(opts ...Option) \*Server {  
    s := \&Server{port: 8080} // Default  
    for \_, opt := range opts { opt(s) }  
    return s  
}

## **2\. Design Patterns & Layout**

### **Standard Project Layout**

Follow the [project-layout](https://github.com/golang-standards/project-layout) pattern used in the implementation\_plan.md:

* /cmd: Entry points (keep these thin).  
* /internal: Private code; prevents other projects from importing your implementation details.  
* /pkg: Public libraries meant for external use.

### **Decoupled Business Logic (Hexagonal Lite)**

Core logic must be decoupled from delivery mechanisms (HTTP/CLI) and storage (SQL/NoSQL).

* **Domain Models:** Defined in the internal/domain or internal/entity package. They should have no struct tags (like json: or db:) if possible.  
* **Services:** Implement business rules using only domain models and interfaces for persistence.  
* **DTOs (Data Transfer Objects):** Only the "delivery" layer (the handlers) should know about JSON tags or template specific shapes.

// ❌ BAD: Business logic mixed with DB tags  
type User struct {  
    ID   int    \`db:"id" json:"user\_id"\`   
    Name string \`db:"name" json:"full\_name"\`  
}

// ✅ GOOD: Clean Domain Model  
package domain  
type User struct {  
    ID   int  
    Name string  
}

### **"Dumb" Templates**

Templates should contain **zero** business logic. They are strictly for presentation.

* **No Fancy Tagging:** Avoid using custom template functions to perform data lookups or formatting that should have happened in the Go code.  
* **Prepared View Models:** The handler should prepare a specific struct that matches the template's needs exactly, so the template only needs to perform simple {{ .Field }} access and {{ range }} loops.  
* **Logic in Go, Data in Templates:** If you need to check if user.IsOver18(), calculate that boolean in the Go handler and pass IsAdult: true to the template.

## **3\. Essential Linters (golangci-lint)**

Use a .golangci.yml file to enforce standards. Below is a "Golden Config" for 2026\.

run:  
  timeout: 5m  
  tests: true

linters:  
  enable:  
    \- staticcheck \# The gold standard for static analysis  
    \- revive      \# Modern replacement for golint  
    \- goerr113    \# Ensures errors are wrapped correctly  
    \- gosec       \# Security scanner for vulnerabilities  
    \- nilaway     \# Google's tool for preventing nil panics  
    \- bodyclose   \# Ensures HTTP response bodies are closed  
    \- unparam     \# Finds unused function parameters  
    \- misspell    \# Catches typos in comments

linters-settings:  
  revive:  
    rules:  
      \- name: exported  
        arguments: \["checkPrivateReceivers"\]

## **4\. Common Go Anti-Patterns to Avoid**

### **The "Log and Return" Double-Whammy**

**Anti-Pattern:** Logging an error and then returning it up the stack.

// ❌ BAD: Resulting logs will be noisy and redundant  
if err \!= nil {  
    log.Printf("failed to save: %v", err)  
    return err  
}

**Fix:** Annotate the error with context and return it. Only log once at the top level (e.g., in a middleware or main.go).

### **Goroutine Leaks**

**Anti-Pattern:** Starting a goroutine without a clear exit strategy or context cancellation.

// ❌ BAD: This goroutine might run forever if the channel is never closed  
go func() {  
    for msg := range msgChan {  
        process(msg)  
    }  
}()

**Fix:** Always use a context.Context or a quit channel to ensure background tasks terminate when the parent operation ends.

### **Interface Over-Abstraction**

**Anti-Pattern:** Creating interfaces for every struct before you have a second implementation.

// ❌ BAD: "Premature Abstraction"  
type UserService interface {  
    Get(id int) \*User  
}

**Fix:** Start with concrete types. Only introduce an interface when you actually need to decouple (e.g., for testing with Testcontainers or multiple providers).

### **Using init() for Logic**

**Anti-Pattern:** Using func init() to set up database connections or global state.

// ❌ BAD: Hidden side effects make testing nearly impossible  
func init() {  
    db, \_ \= sql.Open("postgres", "...")  
}

**Fix:** Use explicit initialization and constructor injection. This is the cornerstone of testable systems.

## **5\. Advanced Testing Practices**

### **Table-Driven Tests (TDT) with Parallelism**

Combine subtests with t.Parallel() to speed up suites and catch race conditions.

func TestCalculate(t \*testing.T) {  
    tests := \[\]struct {  
        name     string  
        input    int  
        expected int  
    }{  
        {"positive", 2, 4},  
        {"negative", \-1, 1},  
    }

    for \_, tt := range tests {  
        t.Run(tt.name, func(t \*testing.T) {  
            t.Parallel()  
            res := Calculate(tt.input)  
            if res \!= tt.expected {  
                t.Errorf("got %d, want %d", res, tt.expected)  
            }  
        })  
    }  
}

### **The AAA Pattern (Arrange, Act, Assert)**

1. **Arrange:** Set up mocks, fakes, or Testcontainers.  
2. **Act:** Call the system under test (SUT).  
3. **Assert:** Use testify/assert or testify/require for clear failure messages.

### **Hermetic E2E with Testcontainers**

Always favor **Real Infrastructure** (Postgres/Redis) over mocks for integration points.

* **Benefit:** Catches real issues like SQL syntax or driver-specific behavior.  
* **Cost:** Higher execution time, mitigated by Go's excellent caching and t.Parallel().

## **6\. Coding Proverbs**

* "Don't design with interfaces, discover them."  
* "Clear is better than clever."  
* "Errors are values."  
* "Don't communicate by sharing memory; share memory by communicating."