# ProfileGo Examples

This directory contains examples demonstrating different ways to use the ProfileGo library.

## Example Structure

```
examples/
├── deprecated/      # Deprecated API usage (for legacy code maintenance)
│   └── legacy.go
├── basic/           # Basic current API usage
│   └── basic.go
├── advanced/        # Advanced current API usage
│   └── advanced.go
└── README.md        # This file
```

## Available Examples

### 1. Deprecated API Example
**Location**: `deprecated/legacy.go`

This example demonstrates the deprecated API usage for maintaining existing code:

```bash
# Run the deprecated example
go run examples/deprecated/legacy.go
```

**Key Features**:
- Uses `profilego.Config` (deprecated)
- Uses `profilego.Init()` (deprecated)
- Shows legacy configuration format
- Includes deprecation warnings in comments

### 2. Basic Current API Example
**Location**: `basic/basic.go`

This example demonstrates the basic recommended way to use the current API:

```bash
# Run the basic example
go run examples/basic/basic.go
```

**Key Features**:
- Uses `config.Config` (current)
- Uses `profilego.InitWithConfig()` (current)
- Minimal configuration
- Simple usage pattern

### 3. Advanced Current API Example
**Location**: `advanced/advanced.go`

This example demonstrates advanced usage of the current API:

```bash
# Run the advanced example
go run examples/advanced/advanced.go
```

**Key Features**:
- Comprehensive configuration with multiple profile types
- Advanced tagging and attributes
- Profiling state management (start/stop)
- Error handling patterns
- Custom attributes usage

## Usage Recommendations

- **New Projects**: Use the `basic/` or `advanced/` examples as starting points
- **Legacy Maintenance**: Use the `deprecated/` example only for maintaining existing code
- **Migration**: Compare deprecated vs. current examples to understand API changes

## API Comparison

| Feature | Deprecated API | Current API |
|---------|---------------|-------------|
| Config Struct | `profilego.Config` | `config.Config` |
| Init Function | `profilego.Init()` | `profilego.InitWithConfig()` |
| Backend Config | String "pyroscope"/"pprof" | `core.PyroscopeBackend`/`core.PprofBackend` |
| Profile Types | Fixed set | Configurable array |
| Additional Attrs | Supported | Enhanced support |

## Running Examples

Each example can be run independently:

```bash
# Basic example
go run examples/basic/basic.go

# Advanced example
go run examples/advanced/advanced.go

# Deprecated example (for reference only)
go run examples/deprecated/legacy.go
```

The examples demonstrate the evolution of the API and provide clear guidance on which approach to use for different scenarios.
