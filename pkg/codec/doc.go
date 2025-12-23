// Package types defines domain types and serialization utilities.
//
// This package provides the core domain types used throughout the Crucible
// ecosystem for representing registry entities:
//
//   - [Namespace]: Organizational units for grouping related resources
//   - [Resource]: Publishable artifacts with multiple versions
//   - [Version]: Immutable published versions with semantic versioning
//   - [Channel]: Mutable pointers to versions (e.g., "stable", "latest")
//   - [Error]: Structured error responses from registry APIs
//
// Each type has an associated media type constant (e.g., [MediaTypeNamespace])
// used for content negotiation in registry APIs.
//
// The package provides format-agnostic serialization supporting JSON, YAML,
// and TOML through the [ContentType] enum. Use [Encode] and [Decode] for
// byte slice operations, or [EncodeFile] and [DecodeFile] for file I/O with
// automatic format detection from file extensions.
//
// The key parameter in encode/decode functions specifies which struct tag to
// use for field mapping (e.g., "field", "json", "yaml"). This allows a single
// struct to support multiple serialization strategies.
//
// Supported file extensions:
//   - JSON: .json
//   - YAML: .yaml, .yml
//   - TOML: .toml
//
// Examples:
//
// Working with domain types:
//
//	ns := types.Namespace{
//	    Name:        "myorg",
//	    Description: "My organization's resources",
//	    Resources: []types.Resource{
//	        {
//	            Name: "my-widget",
//	            Type: "widget",
//	        },
//	    },
//	}
//
// Encoding and decoding data:
//
//	type Config struct {
//	    Name    string `field:"name"`
//	    Version int    `field:"version"`
//	}
//
//	cfg := Config{Name: "app", Version: 1}
//
//	// Encode to bytes
//	data, err := types.Encode(types.ContentTypeJSON, "field", cfg)
//
//	// Decode from bytes
//	var decoded Config
//	err = types.Decode(types.ContentTypeJSON, "field", &decoded, data)
//
//	// File operations with automatic format detection
//	err = types.EncodeFile("config.yaml", "field", cfg)
//	ct, err := types.DecodeFile("config.yaml", "field", &decoded)
//
// The package defines sentinel errors for common failure modes:
//   - [ErrInvalidContentType]: Invalid MIME type string
//   - [ErrUnsupportedContentType]: Content type not supported
//   - [ErrEncodingFailed]: Serialization failure
//   - [ErrDecodingFailed]: Deserialization failure
//
// These errors are wrapped with additional context using the crex package
// error conventions.
package codec
