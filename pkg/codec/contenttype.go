package codec

import (
	"mime"
	"strings"
)

// Represents a content serialization type.
//
// ContentType is the serialization format used for encoding and decoding
// data structures. Supported content types include [ContentTypeJSON],
// [ContentTypeYAML], and [ContentTypeTOML].
type ContentType int

const (
	ContentTypeUnknown ContentType = iota // Unknown or unsupported content type.
	ContentTypeJSON                       // JSON serialization format.
	ContentTypeYAML                       // YAML serialization format.
	ContentTypeTOML                       // TOML serialization format.
)

const (
	ApplicationJSON = "application/json" // Standard JSON media type (RFC 4627, RFC 7159, RFC 8259).
	ApplicationYAML = "application/yaml" // Standard YAML media type (commonly used convention).
	ApplicationTOML = "application/toml" // Standard TOML media type (commonly used convention).

	TextJSON = "text/json" // Legacy JSON media type, sometimes used in older systems.
	TextYAML = "text/yaml" // Legacy YAML media type, sometimes used in older systems.
	TextTOML = "text/toml" // Legacy TOML media type, sometimes used in older systems.

	ApplicationXYAML = "application/x-yaml" // Alternative YAML media type, used by some systems.
	ApplicationXTOML = "application/x-toml" // Alternative TOML media type, used by some systems.

	SuffixJSON = "+json" // JSON syntax suffix (RFC 6839).
	SuffixYAML = "+yaml" // YAML syntax suffix.
	SuffixTOML = "+toml" // TOML syntax suffix.
)

// Returns the MIME type for the format.
//
// This method returns the standard MIME type string associated with the
// content type, in the form "type/subtype". For JSON, YAML, and TOML, this
// returns the preferred application/* media type. Returns an empty string
// for [ContentTypeUnknown].
func (c ContentType) MIMEType() string {
	switch c {
	case ContentTypeJSON:
		return ApplicationJSON
	case ContentTypeYAML:
		return ApplicationYAML
	case ContentTypeTOML:
		return ApplicationTOML
	default:
		return ""
	}
}

// Returns the structured syntax suffix for the content type.
//
// Returns the RFC 6839 structured syntax suffix (e.g., "+json") for the content
// type, which can be appended to vendor-specific media types. This is useful
// when constructing custom media type strings that indicate both a specific
// resource type and its serialization format. Returns an empty string for
// [ContentTypeUnknown].
func (c ContentType) Suffix() string {
	switch c {
	case ContentTypeJSON:
		return SuffixJSON
	case ContentTypeYAML:
		return SuffixYAML
	case ContentTypeTOML:
		return SuffixTOML
	default:
		return ""
	}
}

// Returns a human-readable name for the content type.
//
// This method returns a lowercase string representation of the format name,
// suitable for display purposes, logging, or use in configuration files.
// Returns "unknown" for [ContentTypeUnknown].
func (c ContentType) String() string {
	switch c {
	case ContentTypeJSON:
		return "json"
	case ContentTypeYAML:
		return "yaml"
	case ContentTypeTOML:
		return "toml"
	default:
		return "unknown"
	}
}

// Parses a media type string into a ContentType and media type.
//
// This function recognizes both standard media types (e.g., "application/json",
// "text/yaml") and vendor-specific media types with structured syntax suffixes
// as defined by RFC 6839 (e.g., "application/vnd.example+json"). The parsing
// is case-insensitive, as required by HTTP specifications. The function first
// checks for exact matches against standard media types, then examines the
// structured syntax suffix if present. Media type parameters (such as charset)
// are automatically stripped by the underlying MIME parser.
//
// Returns the ContentType (format), the media type (without suffix and
// parameters), or [ContentTypeUnknown] with [ErrInvalidContentType] if the
// media type string is malformed, or [ErrUnsupportedContentType] if the format
// is not recognized.
func Parse(mediaType string) (ContentType, string, error) {
	parsed, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return ContentTypeUnknown, "", ErrInvalidContentType
	}

	// Check standard media types first
	switch parsed {
	case ApplicationJSON, TextJSON:
		return ContentTypeJSON, parsed, nil
	case ApplicationYAML, ApplicationXYAML, TextYAML:
		return ContentTypeYAML, parsed, nil
	case ApplicationTOML, ApplicationXTOML, TextTOML:
		return ContentTypeTOML, parsed, nil
	}

	// Check for structured syntax suffix (RFC 6839)
	// e.g., application/vnd.example+json
	if idx := strings.LastIndexByte(parsed, '+'); idx > 0 {
		suffix := strings.ToLower(parsed[idx+1:])
		base := parsed[:idx]
		switch suffix {
		case "json":
			return ContentTypeJSON, base, nil
		case "yaml":
			return ContentTypeYAML, base, nil
		case "toml":
			return ContentTypeTOML, base, nil
		}
	}

	return ContentTypeUnknown, "", ErrUnsupportedContentType
}

// Determines the appropriate content type based on the HTTP Accept header.
//
// It performs simplified content negotiation that checks if the Accept header
// matches any supported format. Returns ContentTypeJSON for empty, "*/*",
// unsupported, or unparsable Accept headers. This function does not handle
// quality values (q-parameters) or multiple Accept entries.
func Negotiate(accept string) ContentType {
	if accept == "" || accept == "*/*" {
		return ContentTypeJSON
	}

	ct, _, err := Parse(accept)
	if err != nil {
		return ContentTypeJSON // Default to JSON if parsing fails
	}

	return ct
}
