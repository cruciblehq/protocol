package blueprint

import "github.com/cruciblehq/protocol/pkg/codec"

// Defines a system composition.
//
// A blueprint orchestrates how resources are deployed, declaring which
// resources should be included.
type Blueprint struct {

	// The blueprint version.
	//
	// This is required and must be the first declaration in the blueprint.
	// This value dictates how the rest of the blueprint is interpreted.
	Version int `field:"version"`

	// Lists services to be deployed in this system.
	//
	// Each service instance is exposed through the gateway at its prefix.
	Services []Service `field:"services"`
}

// Represents a service instance within a blueprint.
//
// A service instance defines how a specific version of a service should be
// deployed and exposed within the system. Multiple instances of the same
// service can be deployed with different identifiers.
type Service struct {

	// Unique identifier for this service instance.
	//
	// This ID is used to track the service through plan generation and
	// deployment and should remain stable across blueprint versions.
	ID string `field:"id"`

	// Reference to the service resource.
	//
	// This follows the Crucible reference format, including namespace, name,
	// and version constraint (e.g., "cruciblehq/hub ^1.0.0").
	Reference string `field:"reference"`

	// API prefix for this service.
	//
	// All service endpoints are exposed under this prefix through the system
	// gateway. Prefixes must not conflict or nest with other service prefixes
	// (e.g., "/api/hub" and "/api/hub/users" would conflict).
	Prefix string `field:"prefix"`
}

// Loads a blueprint from a file.
//
// The path parameter specifies the full path to the blueprint file. The file
// format is inferred from the extension (.yaml, .json, .toml).
func Read(path string) (*Blueprint, error) {
	var bp Blueprint
	if _, err := codec.DecodeFile(path, "field", &bp); err != nil {
		return nil, err
	}
	return &bp, nil
}
