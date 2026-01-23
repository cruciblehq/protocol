package plan

import (
	"github.com/cruciblehq/protocol/pkg/codec"
	"github.com/cruciblehq/protocol/pkg/reference"
)

// Represents a deployment plan.
//
// Specifies what resources will be deployed and the infrastructure configuration
// required to run them. Generated during the planning phase by resolving
// references, allocating infrastructure, and determining routing.
type Plan struct {
	Version        int            `field:"version"`
	Services       []Service      `field:"services"`
	Gateway        Gateway        `field:"gateway"`
	Infrastructure Infrastructure `field:"infrastructure"`
}

// Represents a service in the deployment plan.
//
// Contains the resolved reference with exact version and digest.
type Service struct {
	ID        string              `field:"id"`
	Reference reference.Reference `field:"reference"`
}

// Represents the infrastructure configuration for deployment.
//
// Specifies which provider the system will be deployed to. Provider-specific
// configuration is resolved from provider profiles.
type Infrastructure struct {
	Provider string `field:"provider"`
}

// Represents the API gateway configuration.
//
// Defines how external requests are routed to deployed services. For now,
// nginx is used as the gateway implementation with a default listen address.
type Gateway struct {
	Routes []Route `field:"routes"`
}

// Represents a routing rule in the gateway.
//
// Maps request patterns to service instances. The pattern supports path-based
// routing and can be extended to support additional matching criteria.
type Route struct {
	Pattern   string `field:"pattern"`
	ServiceID string `field:"service_id"`
}

// Saves the plan to a file.
//
// The file format is inferred from the path extension (.json, .yaml, .toml).
func (p *Plan) Write(path string) error {
	return codec.EncodeFile(path, "field", p)
}

// Loads a plan from a file.
//
// The file format is inferred from the path extension (.json, .yaml, .toml).
func Read(path string) (*Plan, error) {
	var p Plan
	if _, err := codec.DecodeFile(path, "field", &p); err != nil {
		return nil, err
	}
	return &p, nil
}
