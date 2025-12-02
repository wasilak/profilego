# Generic Go Library Structure Steering

## Overview
This document provides a flexible structure template that can be adapted for various Go library projects.

## Generic Project Structure
```
project-name/
├── cmd/                 # CLI command implementations (if applicable)
├── pkg/                 # Core library packages
│   ├── [feature1]/     # Feature-specific packages
│   ├── [feature2]/     # Feature-specific packages
│   └── utils/           # Utility functions
├── internal/            # Private library code (if needed)
├── docs/                # Documentation
├── examples/            # Usage examples
├── test/                # Test utilities and mocks
├── go.mod/go.sum        # Go module files
├── README.md            # Main documentation
└── LICENSE              # License file
```

## Package Design Guidelines
- Each package should have a single, well-defined responsibility
- Package names should be short, lowercase, and singular
- Public APIs should be minimal and stable
- Internal implementation details should be unexported
- Use `internal/` directory for code that shouldn't be imported by other projects

## Import Structure
- Standard library imports first
- Third-party imports grouped together
- Project imports grouped together
- Blank lines separating import groups

## Naming Conventions
- Use MixedCaps for exported names (PascalCase)
- Use camelCase for unexported names
- Keep names short but descriptive
- Avoid stuttering (e.g., use `server.New` not `server.NewServer`)
- Follow Go naming conventions consistently

## API Design Guidelines
- Use consistent function signatures
- Return meaningful error responses
- Follow Go interface design best practices
- Maintain backward compatibility in public APIs
- Provide comprehensive documentation

## Configuration Management
- Support multiple configuration methods when applicable
- Use environment variables for runtime configuration
- Provide clear validation for configuration values
- Use struct tags for configuration parsing

## Error Handling Structure
- Use consistent error wrapping with `fmt.Errorf` and `%w`
- Provide contextual error messages
- Implement proper error handling for all external operations
- Use sentinel errors for expected error conditions

## Testing Structure
- Unit tests in `*_test.go` files
- Integration tests in separate test files when needed
- Test files in the same directory as the code they test
- Use table-driven tests where appropriate
- Maintain high test coverage for critical paths

## Documentation Structure
- Inline Go documentation for public APIs
- Markdown documentation in `docs/` directory
- Example usage in `examples/` directory
- README files at each major directory level
- Clear API documentation and usage examples

## Build and Deployment Structure
- Simple build process using `go build`
- CI/CD configuration files
- Release artifacts management
- Version management

## Tooling Recommendations
- Use `ripgrep` (rg) for fast code searching
- Use `fd` for efficient file finding
- Use `golangci-lint` for comprehensive linting
- Use `go vet` and `staticcheck` for code analysis

## Adaptation Guidelines
This structure should be adapted for specific projects by:
1. Adjusting directory structure based on project needs
2. Adding project-specific directories
3. Modifying based on library vs. application focus
4. Including domain-specific organization requirements
