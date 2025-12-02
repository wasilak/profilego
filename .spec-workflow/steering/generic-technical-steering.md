# Generic Go Library Technical Steering

## Overview
This document provides technical guidance that can be adapted for various Go library development projects.

## Architecture Principles
- Separation of concerns with clear component responsibilities
- Configuration flexibility over hardcoded behaviors
- Comprehensive error handling with meaningful messages
- Observability through structured logging
- Testability through clear interfaces

## Technology Stack Recommendations
- **Language**: Go (latest stable version)
- **Testing**: Go's built-in testing package
- **Logging**: Standard library `log/slog` package
- **Configuration**: Environment variables and struct-based config
- **API Documentation**: Inline Go documentation
- **Dependency Management**: Go modules

## Code Organization Guidelines
- **cmd/**: CLI command implementations (if applicable)
- **pkg/**: Core library packages
- **internal/**: Private libraries (when needed)
- **docs/**: Documentation files
- **examples/**: Usage examples

## Development Practices
- Follow Go best practices and idioms
- Use structured logging with `slog`
- Include comprehensive tests for new functionality
- Maintain complete documentation for public APIs
- Follow clean interface design principles
- Use consistent error response formats

## Quality Gates - Post-Implementation Checklist
After each implementation task, verify the following:

### Code Quality Checks
```bash
# Check for vet errors (REQUIRED)
go vet ./...

# Format check
go fmt ./...

# Import organization
goimports -w ./...

# Linting (recommended)
golangci-lint run
```

### Testing
```bash
# Run all tests
go test ./... -v

# Run specific package tests
go test ./pkg/[package-name] -v

# Run with race detector for concurrent code
go test -race ./...
```

### Build Verification
```bash
# Build the library
go build ./...
```

**FAIL** any task if:
- `go vet ./...` produces errors
- Tests fail
- Build fails
- API changes break existing functionality

## Performance Considerations
- Minimize memory allocations in hot paths
- Optimize for Go's performance characteristics
- Use efficient data structures
- Implement proper resource management

## Security Considerations
- Input validation for all public APIs
- Proper error handling for security-sensitive operations
- Secure defaults for configuration
- Protection against common attack vectors

## Tooling Best Practices
- **Search**: Prefer `ripgrep` (rg) over grep for faster searches
  ```bash
  rg "pattern" .
  rg -i "case-insensitive" .
  rg -A 3 -B 3 "context" .  # Show context
  ```

- **File finding**: Prefer `fd` over find for better performance
  ```bash
  fd "pattern" .
  fd -e go .  # Find Go files
  fd -H "hidden" .  # Include hidden files
  ```

- **Code analysis**: Use modern Go tooling
  ```bash
  # Fast code search
  rg "functionName" --type=go

  # Find files efficiently
  fd -e go -x go vet

  # Static analysis
  staticcheck ./...
  ```

## Adaptation Guidelines
This technical steering should be adapted for specific projects by:
1. Adding project-specific technology requirements
2. Including domain-specific performance considerations
3. Adjusting based on project complexity and scale
4. Adding project-specific security requirements
