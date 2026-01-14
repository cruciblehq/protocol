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

### [`pkg/manifest`](pkg/manifest)

Resource manifest parsing and validation. Looks for Crucible manifest files
under `path/.cruciblerc/manifest.yaml` and parses and validates according to
Crucible's expected format.

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

Creation and extraction of zstd-compressed tar archives. Used for handling
Crucible resource archives.

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

### [`pkg/codec`](pkg/codec)

Domain types and serialization with support for format-agnostic serialization,
including JSON, YAML, and TOML.

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

### [`pkg/registry`](pkg/registry)

Artifact registry implementation with hierarchical storage for versioned
resources. Provides a complete registry interface with namespace, resource,
version, and channel management.

The package includes two implementations:

- **SQLRegistry**: Local registry backed by SQLite
- **Client**: Remote registry accessed via HTTP

#### Local Registry (SQLRegistry)

```go
import (
    "github.com/cruciblehq/protocol/pkg/registry"
    "database/sql"
)

// Create a new local registry
db, _ := sql.Open("sqlite3", "registry.db?_foreign_keys=on")
reg, err := registry.NewSQLRegistry(db, "/path/to/archives", logger)

// Create namespace
ns, err := reg.CreateNamespace(ctx, registry.NamespaceInfo{
    Name:        "myorg",
    Description: "My organization",
})

// Create resource
res, err := reg.CreateResource(ctx, "myorg", registry.ResourceInfo{
    Name:        "mywidget",
    Type:        "widget",
    Description: "My widget",
})

// Create version
ver, err := reg.CreateVersion(ctx, "myorg", "mywidget", registry.VersionInfo{
    String: "1.0.0",
})

// Upload archive
file, _ := os.Open("widget.tar.zst")
defer file.Close()
ver, err = reg.UploadArchive(ctx, "myorg", "mywidget", "1.0.0", file)

// Create channel
ch, err := reg.CreateChannel(ctx, "myorg", "mywidget", registry.ChannelInfo{
    Name:        "stable",
    Version:     "1.0.0",
    Description: "Stable channel",
})

// Download archive
reader, err := reg.DownloadArchive(ctx, "myorg", "mywidget", "1.0.0")
defer reader.Close()
```

#### Remote Registry (Client)

```go
import "github.com/cruciblehq/protocol/pkg/registry"

// Create a client for a remote registry
client, err := registry.NewClient("https://hub.example.com", nil)

// Client implements the same Registry interface
ns, err := client.CreateNamespace(ctx, registry.NamespaceInfo{
    Name:        "myorg",
    Description: "My organization",
})

// List resources
resources, err := client.ListResources(ctx, "myorg")

// Read specific version
ver, err := client.ReadVersion(ctx, "myorg", "mywidget", "1.0.0")

// Upload archive to remote registry
file, _ := os.Open("widget.tar.zst")
defer file.Close()
ver, err = client.UploadArchive(ctx, "myorg", "mywidget", "1.0.0", file)

// Download from remote registry
reader, err := client.DownloadArchive(ctx, "myorg", "mywidget", "1.0.0")
defer reader.Close()

// Handle errors
_, err = client.ReadNamespace(ctx, "nonexistent")
if err != nil {
    if registryErr, ok := err.(*registry.Error); ok {
        if registryErr.Code == registry.ErrorCodeNotFound {
            // Handle not found
        }
    }
}
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
