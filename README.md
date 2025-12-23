# Crucible Protocol

Go implementation of the Crucible protocol specifications, providing core data
structures, parsers, and utilities for the Crucible ecosystem.

## Overview

The `protocol` module defines the types and operations used across Crucible
tools and services. It implements resource references, semantic versioning,
manifest parsing, archive handling, and serialization utilities.

## Packages

### [`pkg/reference`](pkg/reference)

Resource identification and version constraint parsing. Implements the complete
Crucible reference syntax for locating and versioning resources.

**Key features:**
- Parse resource identifiers (scheme, registry, namespace, name)
- Semantic version constraints with operators, ranges, wildcards
- Channel-based references for release tracks
- Content-addressable digests for immutable references

```go
import "github.com/cruciblehq/protocol/pkg/reference"

// Parse a complete reference
ref, err := reference.Parse("widget myorg/mywidget >=1.0.0 <2.0.0", "widget", nil)

// Parse just an identifier
id, err := reference.ParseIdentifier("myorg/mywidget", "widget", nil)

// Parse version constraints
vc, err := reference.ParseVersionConstraint(">=1.0.0 <2.0.0")
ok, err := vc.Matches("1.5.0") // true
```

### [`pkg/types`](pkg/types)

Domain types and serialization utilities for registry entities.

**Key features:**
- Core types: `Namespace`, `Resource`, `Version`, `Channel`, `Error`
- Format-agnostic serialization (JSON, YAML, TOML)
- Media type constants for content negotiation
- File I/O with automatic format detection

```go
import "github.com/cruciblehq/protocol/pkg/types"

// Work with domain types
ns := types.Namespace{
    Name: "myorg",
    Resources: []types.Resource{
        {Name: "mywidget", Type: "widget"},
    },
}

// Encode/decode with struct tags
data, err := types.Encode(types.ContentTypeJSON, "field", ns)
err = types.Decode(types.ContentTypeJSON, "field", &ns, data)

// File operations with auto-detection
err = types.EncodeFile("namespace.yaml", "field", ns)
```

### [`pkg/manifest`](pkg/manifest)

Resource manifest parsing and validation.

**Key features:**
- Parse `.cruciblerc/manifest.yaml` files
- Type-specific configuration (Widget, Service)
- Validation of manifest structure and fields

```go
import "github.com/cruciblehq/protocol/pkg/manifest"

// Read manifest from resource directory
m, err := manifest.Read("/path/to/resource")

// Access type-specific config
switch cfg := m.Config.(type) {
case *manifest.Widget:
    fmt.Println(cfg.Build.Main)
case *manifest.Service:
    fmt.Println(cfg.Build.Main)
}
```

### [`pkg/archive`](pkg/archive)

Secure creation and extraction of zstd-compressed tar archives.

**Key features:**
- Zstd compression for efficient storage
- Path validation to prevent directory traversal
- Symlink rejection for security
- Atomic extraction with automatic cleanup on error

```go
import "github.com/cruciblehq/protocol/pkg/archive"

// Create archive
err := archive.Create("mydir", "output.tar.zst")

// Extract archive
err = archive.Extract("output.tar.zst", "destination")

// Extract from reader
file, _ := os.Open("output.tar.zst")
defer file.Close()
err = archive.ExtractFromReader(file, "destination")
```

### [`pkg/crex`](pkg/crex)

Structured error management for user-facing applications.

**Key features:**
- Rich error context (description, reason, fallback, cause)
- Error classification (User, System, Programming, Bug)
- Structured logging integration via `slog.LogValuer`
- Custom log handler with buffering support
- TTY-aware formatting (pretty vs JSON)

```go
import "github.com/cruciblehq/protocol/pkg/crex"

// Create structured errors
err := crex.UserError("could not save file", "insufficient permissions").
    Fallback("Try running with elevated privileges.").
    Detail("path", filePath).
    Cause(previousError).
    Err()

// Use simple wrapping at package boundaries
return crex.Wrap(ErrInvalidInput, err)
```

## Installation

```bash
go get github.com/cruciblehq/protocol
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/reference/...
```

## License

All rights reserved.
