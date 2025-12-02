# Generic Go Library Development Principles

## Overview
This document establishes fundamental principles for Go library development that can be applied across multiple projects. These principles focus on clean, maintainable, and efficient Go code development.

## Core Development Principles

### Greenfield Development Approach
- Implement functionality as new code from the ground up
- No need to consider backward compatibility for new features
- Focus on optimal solutions without legacy constraints
- Implement breaking changes when they provide better architecture

### Code Quality Standards
- Follow Go best practices and idioms
- Use modern Go features (generics, structured logging, etc.)
- Prioritize clean, readable, and maintainable code
- Implement comprehensive error handling

## Tooling Preferences
- **Search**: Use `ripgrep` (rg) instead of grep when available
- **File finding**: Use `fd` instead of find when available
- **Code analysis**: Use `go vet`, `staticcheck`, and `golangci-lint`
- **Testing**: Use Go's built-in testing with table-driven tests

## Application to Specifications
When creating specifications for new features:
- Write requirements assuming greenfield development
- Design without migration constraints for new functionality
- Create tasks focused on new code implementation
- Include backward compatibility considerations only for existing features

## Benefits of This Approach
- Cleaner, more maintainable code architecture
- No technical debt from legacy implementation decisions
- Ability to implement optimal solutions without constraints
- Faster development cycles without complex migration logic

## When to Apply
This principle applies to:
- New feature development where functionality doesn't exist
- New library components and modules
- New API endpoints and interfaces
- New utility functions and packages

## When to Consider Other Approaches
For modifications to existing functionality:
- Consider backward compatibility
- Implement migration strategies when needed
- Include fallback mechanisms for critical features
