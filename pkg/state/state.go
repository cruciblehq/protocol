package state

import (
	"time"

	"github.com/cruciblehq/protocol/pkg/codec"
	"github.com/cruciblehq/protocol/pkg/reference"
)

// Represents the current state of a deployment.
//
// Records what resources have been deployed and their runtime identifiers.
// Used for incremental deployments and resource lifecycle management.
type State struct {
	Version    int        `field:"version"`
	Deployment Deployment `field:"deployment"`
	Services   []Service  `field:"services"`
}

// Represents deployment metadata.
//
// Contains information about when and how the deployment was executed.
// Expandable for future metadata like deploy user, environment, etc.
type Deployment struct {
	DeployedAt time.Time `field:"deployed_at"`
}

// Represents a service that has been deployed.
//
// Tracks the service identity and its provider-specific resource identifier.
type Service struct {
	ID         string              `field:"id"`
	Reference  reference.Reference `field:"reference"`
	ResourceID string              `field:"resource_id"`
}

// Saves the state to a file.
func (s *State) Write(path string) error {
	return codec.EncodeFile(path, "field", s)
}

// Loads a state from a file.
func Read(path string) (*State, error) {
	var s State
	if _, err := codec.DecodeFile(path, "field", &s); err != nil {
		return nil, err
	}
	return &s, nil
}
