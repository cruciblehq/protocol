package manifest

// Holds configuration specific to service resources.
//
// Service resources are backend components that provide functionality to other
// systems by exposing an API. This structure defines configurations that are
// unique to service resources, such as container image references.
type Service struct {

	// Holds build-related configuration for service resources.
	Build struct {
		Image string `field:"image"` // Path to pre-built OCI image tarball (e.g., "dist/image.tar").
	} `field:"build"`
}
