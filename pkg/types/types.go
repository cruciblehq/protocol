package types

import (
	"github.com/cruciblehq/protocol/pkg/reference"
)

const (
	MediaTypeNamespace = "application/vnd.crucible.namespace.v0"    // Media type for Namespace.
	MediaTypeResource  = "application/vnd.crucible.resource.v0"     // Media type for Resource.
	MediaTypeArchive   = "application/vnd.crucible.archive.v0+zstd" // Media type for resource package archives.
	MediaTypeError     = "application/vnd.crucible.error.v0"        // Media type for Error.
)

// Represents an organizational unit for grouping related resources.
//
// Namespaces provide isolation and organization for resources. Namespace names
// are globally unique and serve as the first component of resource identifiers.
//
// The corresponding media type is application/vnd.crucible.namespace.v0.
type Namespace struct {
	Name        string     `field:"name"`
	Description string     `field:"description,omitempty"`
	Resources   []Resource `field:"resources"`
}

// Represents a publishable resource within a namespace.
//
// Resources are the primary artifacts managed by the registry, such as widgets,
// services, or templates. Each resource has a unique name within its namespace
// and can have multiple published versions. Resources track their type (e.g.,
// "widget", "service") and maintain metadata like descriptions that can be
// updated over time without affecting published versions.
//
// The corresponding media type is application/vnd.crucible.resource.v0.
type Resource struct {
	Name        string    `field:"name"`
	Type        string    `field:"type"`
	Description string    `field:"description,omitempty"`
	Versions    []Version `field:"versions"`
	Channels    []Channel `field:"channels"`
}

// Represents an immutable published version of a resource.
//
// Each version corresponds to a specific archive uploaded to the registry,
// identified by its semantic version number and content digest. Versions are
// immutable once published and cannot be modified or deleted. They provide
// the foundation for reproducible builds and dependency resolution.
type Version struct {
	Version reference.Version `field:"version"`
	Digest  reference.Digest  `field:"digest"`
	Size    int64             `field:"size"`
}

// Represents a mutable named pointer to a resource version.
//
// Channels provide stable release tracks like "stable", "beta", or "latest"
// that can be updated over time to point to different versions. Consumers can
// fetch resources by channel name (e.g., :stable) to always get whichever
// version is currently assigned to that channel, without specifying exact
// version numbers. This makes them convenient for tracking release tracks,
// but they should not be used where reproducibility is required. For that
// reason, channels are not allowed when referenced from published resources.
type Channel struct {
	Name        string            `field:"name"`
	Description string            `field:"description,omitempty"`
	Version     reference.Version `field:"version"`
	Digest      reference.Digest  `field:"digest"`
	Size        int64             `field:"size"`
}

// Represents an error response from the registry API.
//
// The Code field contains a machine-readable error identifier that clients
// can use for programmatic error handling. The Message field provides a
// human-readable description of the error.
//
// The corresponding media type is application/vnd.crucible.error.v0.
type Error struct {
	Code    string `field:"code"`
	Message string `field:"message"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Code + ": " + e.Message
}
