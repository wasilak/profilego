# Generic Go Library Steering Documents

This directory contains generic steering documents that can be adapted for any Go library development project. These documents provide flexible templates and best practices that work across multiple projects.

## Available Documents

### 1. [Generic Development Principles](generic-development-principles.md)
Core principles for Go library development including:
- Greenfield development approach
- Code quality standards
- Tooling preferences (ripgrep, fd, etc.)
- When to apply different approaches

### 2. [Generic Product Steering](generic-product-steering.md)
Product-level guidance including:
- Core product principles
- Target use cases
- Product boundaries
- Success metrics
- Adaptation guidelines

### 3. [Generic Structure Steering](generic-structure-steering.md)
Project structure and organization including:
- Flexible project structure template
- Package design guidelines
- Import and naming conventions
- API design principles
- Documentation structure
- Tooling recommendations

### 4. [Generic Technical Steering](generic-technical-steering.md)
Technical implementation guidance including:
- Architecture principles
- Technology stack recommendations
- Development practices
- Quality gates and checklists
- Performance and security considerations
- Tooling best practices with ripgrep/fd examples

## How to Use These Documents

1. **Start with the principles**: Read [Generic Development Principles](generic-development-principles.md) first
2. **Adapt for your project**: Copy and modify the documents for your specific needs
3. **Use tooling recommendations**: Implement the suggested tools (ripgrep, fd) in your workflow
4. **Follow quality gates**: Use the provided checklists for consistent quality

## Key Features

- **Tooling Focus**: Emphasis on modern tools like `ripgrep` and `fd` for better performance
- **Flexibility**: Designed to work across different Go library projects
- **Best Practices**: Incorporates Go idioms and modern development practices
- **Quality Assurance**: Includes comprehensive quality gates and checklists

## Usage Examples

```bash
# Fast code search using ripgrep
rg "functionName" --type=go

# Efficient file finding with fd
fd -e go -x go vet

# Quality checks
go vet ./...
go test ./... -race
```

These documents provide a solid foundation for Go library development while allowing flexibility for project-specific requirements.
