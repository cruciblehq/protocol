package blueprint

// Represents a service instance within a blueprint.
//
// A service instance defines how a specific version of a service should be
// deployed and exposed within the system. Multiple instances of the same
// service can be deployed with different names and prefixes.
type Service struct {

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
