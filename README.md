# ADHD

**A**synchronous **D**istracted **H**yperactive  **D**evelopment

The Go context package for when your goroutines can't sit still and need constant supervision.

> "Just like real ADHD, this package helps you manage multiple things at once while occasionally forgetting what you were doing in the first place." - pixie

## Why ADHD?

Ever tried to manage multiple goroutines and felt like herding cats? Your code jumping between tasks faster than a squirrel on espresso? Welcome to ADHD - the context manager that gets it.

## Features

- **Background()** - Root context for applications
- **TODO()** - Placeholder context for development
- **WithCancel()** - Manual context cancellation
- **WithTimeout()** - Duration-based context expiration
- **WithDeadline()** - Absolute time-based context expiration
- **WithValue()** - Key-value storage with context chaining
- **Select()** & **Race()** - Multi-context racing utilities
- **IsDone()** - Non-blocking context completion check



## Quick Start

```go
package main

import (
    "fmt"
    "time"
    "github.com/pixperk/adhd/adhd"
)

func main() {
    ctx, cancel := adhd.WithTimeout(adhd.Background(), 5*time.Second)
    defer cancel()

    select {
    case <-time.After(3 * time.Second):
        fmt.Println("Work completed")
    case <-ctx.Done():
        fmt.Printf("Context cancelled: %v\n", ctx.Err())
}
```

## API Reference

### Context Creation

```go
ctx := adhd.Background()    // Root context
ctx := adhd.TODO()          // Development placeholder
```

### Cancellation

```go
ctx, cancel := adhd.WithCancel(parent)
cancel() // Cancel the context
```

### Timeouts & Deadlines

```go
ctx, cancel := adhd.WithTimeout(parent, 5*time.Second)
ctx, cancel := adhd.WithDeadline(parent, deadline)
defer cancel()
```

### Value Storage

```go
ctx := adhd.WithValue(parent, "key", "value")
value := ctx.Value("key")
```

### Context Racing

```go
result := <-adhd.Select(ctx1, ctx2, ctx3)
fmt.Printf("Context %d completed first\n", result.Index)

winner := <-adhd.Race(ctx1, ctx2)
fmt.Printf("Winner error: %v\n", winner.Error)
```

### Utilities

```go
if adhd.IsDone(ctx) {
    fmt.Println("Context is done")
}

err := adhd.WaitFor(ctx)
fmt.Printf("Context completed with: %v\n", err)
```

## Error Handling

```go
if ctx.Err() == adhd.ErrCanceled {
    fmt.Println("Context was cancelled")
}

if ctx.Err() == adhd.ErrDeadlineExceeded {
    fmt.Println("Context deadline exceeded")
}
```

## Examples

### HTTP Server with Graceful Shutdown

```go
func main() {
    server := &http.Server{Addr: ":8080"}
    
    shutdownCtx, cancel := adhd.WithCancel(adhd.Background())
    defer cancel()
    
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, os.Interrupt)
        <-sigChan
        cancel()
    }()
    
    go server.ListenAndServe()
    
    <-shutdownCtx.Done()
    server.Shutdown(context.Background())
}
```

### Worker Pool with Timeout

```go
func worker(ctx adhd.ADHD, jobs <-chan int, results chan<- int) {
    for {
        select {
        case job := <-jobs:
            results <- job * 2
        case <-ctx.Done():
            return
        }
    }
}

func main() {
    ctx, cancel := adhd.WithTimeout(adhd.Background(), 10*time.Second)
    defer cancel()
    
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    
    for i := 0; i < 3; i++ {
        go worker(ctx, jobs, results)
    }
    
    for i := 1; i <= 5; i++ {
        jobs <- i
    }
    close(jobs)
    
    for i := 0; i < 5; i++ {
        select {
        case result := <-results:
            fmt.Printf("Result: %d\n", result)
        case <-ctx.Done():
            fmt.Printf("Timeout: %v\n", ctx.Err())
            return
        }
    }
}
```

### Context Racing

```go
func fetchData(ctx adhd.ADHD, apiName string) adhd.ADHD {
    resultCtx, cancel := adhd.WithCancel(ctx)
    go func() {
        defer cancel()
        delay := time.Duration(rand.Intn(1000)) * time.Millisecond
        time.Sleep(delay)
    }()
    return resultCtx
}

func main() {
    api1 := fetchData(adhd.Background(), "API1")
    api2 := fetchData(adhd.Background(), "API2") 
    api3 := fetchData(adhd.Background(), "API3")
    
    result := <-adhd.Select(api1, api2, api3)
    
    apiNames := []string{"API1", "API2", "API3"}
    fmt.Printf("Winner: %s\n", apiNames[result.Index])
}
```

## Testing

```bash
go test ./adhd
go test -v ./adhd
go test -cover ./adhd
```

## Performance

ADHD is designed to be lightweight and efficient:

- Minimal memory overhead
- No goroutine leaks  
- Efficient channel operations
- Thread-safe operations

## ADHD vs Standard Context

| Feature | ADHD | context | 
|---------|------|---------|
| Basic contexts | Yes | Yes |
| Cancellation | Yes | Yes |
| Timeouts | Yes | Yes |
| Values | Yes | Yes |
| Context racing | Yes | No |
| Select utilities | Yes | No |

## Contributing

1. Fork the repository
2. Create your feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request



---

**ADHD**: Because sometimes you need to manage multiple things at once, and that's perfectly okay.
