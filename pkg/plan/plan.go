package plan

import (
	"github.com/cruciblehq/protocol/pkg/codec"
)

// Represents a deployment plan.
//
// Specifies what resources will be deployed and the infrastructure configuration
// required to run them. Generated during the planning phase by resolving
// references, allocating infrastructure, and determining routing.
type Plan struct {
	Version      int           `field:"version"`
	Services     []Service     `field:"services"`
	Compute      []Compute     `field:"compute"`
	Environments []Environment `field:"environments"`
	Bindings     []Binding     `field:"bindings"`
	Gateway      Gateway       `field:"gateway"`
}

// Represents a service in the deployment plan.
//
// Contains the resolved reference with exact version and digest.
type Service struct {
	ID        string `field:"id"`
	Reference string `field:"reference"`
}

// Represents a compute resource in the deployment plan.
//
// Defines the compute instance to provision. This only describes the
// infrastructure resource (what to allocate), not what runs on it.
type Compute struct {
	ID           string `field:"id"`
	Provider     string `field:"provider"`
	InstanceType string `field:"instance_type"`
}

// Represents an environment configuration.
//
// Defines a set of environment variables that can be associated with deployments.
// Environments are declared separately and referenced by deployments.
type Environment struct {
	ID        string            `field:"id"`
	Variables map[string]string `field:"variables"`
}

// Represents a binding of a service to compute infrastructure.
//
// Associates a service with a compute instance and optional environment
// configuration. Multiple bindings of the same service enable replicas.
// Multiple services on the same compute enable co-location.
type Binding struct {
	Service     string `field:"service"`
	Compute     string `field:"compute"`
	Environment string `field:"environment,omitempty"`
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
	Pattern string `field:"pattern"`
	Service string `field:"service"`
}

// Saves the plan to a file.
//
// The file format is inferred from the path extension (.json, .yaml, .toml).
// The indent parameter controls whether JSON output should be pretty-printed.
func (p *Plan) Write(path string, indent bool) error {
	return codec.EncodeFile(path, "field", indent, p)
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
