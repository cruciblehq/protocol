package blueprint

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
	// Each service instance declared here is exposed through the gateway at
	// its configured prefix.
	Services []Service `field:"services"`
}
