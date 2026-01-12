package manifest

import (
	"github.com/cruciblehq/protocol/internal/helpers"
	"github.com/cruciblehq/protocol/pkg/codec"
)

// Loads and parses a manifest file.
//
// The path parameter specifies the full path to the manifest file. The file
// format is inferred from the extension (.yaml, .json, .toml). The function
// reads and unmarshals the file contents according to the [Manifest] structure.
// The structure is expected to conform to the Crucible manifest schema,
// identified by "field" struct tags. Returns the parsed [Manifest] on success,
// or an error if the file could not be read or parsed.
func Read(path string) (*Manifest, error) {

	// Decode file into raw map
	var raw map[string]any
	if _, err := codec.DecodeFile(path, "field", &raw); err != nil {
		return nil, helpers.Wrap(ErrManifestReadFailed, err)
	}

	// Decode into Manifest struct
	var m Manifest
	if err := decodeManifest(raw, &m); err != nil {
		return nil, helpers.Wrap(ErrManifestReadFailed, err)
	}

	return &m, nil
}

// Decodes a raw map into a [Manifest] structure.
//
// The raw parameter is a map representing the unmarshaled content. The manifest
// parameter is a pointer to the [Manifest] structure where the decoded data
// should be stored. The function first decodes common fields into the manifest,
// then resolves the resource type to determine the concrete manifest type.
func decodeManifest(raw map[string]any, manifest *Manifest) error {

	// Decode common fields
	if err := codec.DecodeMap(raw, "field", manifest); err != nil {
		return err
	}

	// Resolve type-specific config
	configs := map[string]any{
		"widget":  &Widget{},
		"service": &Service{},
	}

	target, ok := configs[manifest.Resource.Type]
	if !ok {
		return ErrUnknownResourceType
	}

	// Decode type-specific config
	if err := codec.DecodeMap(raw, "field", target); err != nil {
		return err
	}

	// Assign to manifest
	manifest.Config = target

	return nil
}
