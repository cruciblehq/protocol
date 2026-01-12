package manifest

const (

	// Distribution directory for service build output
	ServiceDistDirectory = "dist"

	// Distribution file for service image output.
	ServiceDistImage = ServiceDistDirectory + "/" + "image.tar"
)

// Holds configuration specific to service resources.
//
// Service resources are backend components that provide functionality to other
// systems by exposing an API. This structure defines configurations that are
// unique to service resources, such as container image references.
type Service struct {

	// Holds build-related configuration for service resources.
	Build struct {
		Image string `field:"image"` // Path to pre-built OCI image tarball (e.g., "build/image.tar").
	} `field:"build"`
}
