package oci

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/cruciblehq/protocol/internal/helpers"
	"github.com/cruciblehq/protocol/pkg/archive"
)

const (

	// The media type for an OCI Image Index.
	OCIImageIndexMediaType = "application/vnd.oci.image.index.v1+json"

	// Directory path prefix for SHA256 blobs within an OCI tarball.
	OCIBlobsSHA256Dir = "blobs/sha256"
)

var (

	// ErrInvalidImage indicates the OCI image is malformed or corrupted.
	ErrInvalidImage = errors.New("invalid OCI image")

	// ErrSinglePlatform indicates the image only supports a single platform.
	ErrSinglePlatform = errors.New("single-platform image")
)

// Represents an OCI Image Index structure.
type Index struct {
	SchemaVersion int          `json:"schemaVersion"` // OCI image index schema version
	MediaType     string       `json:"mediaType"`     // OCI media type
	Manifests     []Descriptor `json:"manifests"`     // List of image manifests
}

// Represents an OCI Content Descriptor.
type Descriptor struct {
	MediaType string    `json:"mediaType"`          // OCI media type
	Digest    string    `json:"digest"`             // Content digest
	Size      int64     `json:"size"`               // Size in bytes
	Platform  *Platform `json:"platform,omitempty"` // Platform information
}

// Represents an OCI Platform structure.
type Platform struct {
	Architecture string `json:"architecture"`      // CPU architecture
	OS           string `json:"os"`                // Operating system
	Variant      string `json:"variant,omitempty"` // CPU variant (optional)
}

// Returns the platforms required for universal deployment.
//
// Images must support all these platforms to be considered universally
// deployable within the Crucible ecosystem. The list may be updated in future
// releases to include additional architectures as needed.
func RequiredPlatforms() []string {
	return []string{
		"linux/amd64", // x86_64 servers, most cloud providers
		"linux/arm64", // ARM servers, Apple Silicon, modern cloud instances
	}
}

// Reads and parses the index.json from an OCI tarball.
func ReadIndex(imagePath string) (*Index, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return nil, helpers.Wrap(ErrInvalidImage, err)
	}
	defer f.Close()

	tr := tar.NewReader(f)
	indexData, err := archive.FindInTar(tr, "index.json")
	if err != nil {
		return nil, helpers.Wrap(ErrInvalidImage, err)
	}

	if indexData == nil {
		return nil, helpers.Wrap(ErrInvalidImage, errors.New("index.json not found"))
	}

	var index Index
	if err := json.Unmarshal(indexData, &index); err != nil {
		return nil, helpers.Wrap(ErrInvalidImage, err)
	}

	return &index, nil
}

// Checks if an OCI index points to a nested manifest list.
//
// Docker Buildx creates a nested structure where the top-level index.json
// references another index in the blobs directory containing the actual
// platform-specific manifests.
func IsNestedIndex(index *Index) bool {
	return len(index.Manifests) > 0 && index.Manifests[0].MediaType == OCIImageIndexMediaType
}

// Reads a nested OCI index from the blobs directory.
func ReadNestedIndex(imagePath, digest string) (*Index, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return nil, helpers.Wrap(ErrInvalidImage, err)
	}
	defer f.Close()

	if digest == "" {
		return nil, helpers.Wrap(ErrInvalidImage, errors.New("manifest digest is empty"))
	}

	// Remove "sha256:" prefix if present
	if len(digest) > 7 && digest[:7] == "sha256:" {
		digest = digest[7:]
	}

	tr := tar.NewReader(f)
	blobPath := path.Join(OCIBlobsSHA256Dir, digest)
	nestedData, err := archive.FindInTar(tr, blobPath)
	if err != nil {
		return nil, helpers.Wrap(ErrInvalidImage, err)
	}

	if nestedData == nil {
		return nil, helpers.Wrap(ErrInvalidImage, errors.New("nested manifest list not found"))
	}

	var index Index
	if err := json.Unmarshal(nestedData, &index); err != nil {
		return nil, helpers.Wrap(ErrInvalidImage, err)
	}

	return &index, nil
}

// Extracts valid platform identifiers from an OCI index.
//
// Excludes attestation manifests and descriptors with unknown os/architecture.
// Returns a map where keys are platform strings in "os/arch" format.
func Platforms(index *Index) map[string]bool {
	platforms := make(map[string]bool)
	for _, manifest := range index.Manifests {
		if manifest.Platform == nil {
			continue
		}
		// Skip attestation manifests with unknown os/arch
		if manifest.Platform.OS == "unknown" || manifest.Platform.Architecture == "unknown" {
			continue
		}
		key := fmt.Sprintf("%s/%s", manifest.Platform.OS, manifest.Platform.Architecture)
		platforms[key] = true
	}
	return platforms
}

// Validates that an OCI tarball contains multiple platforms.
//
// Returns ErrSinglePlatform if the image supports only one platform or none
// and ErrInvalidImage if the tarball structure is invalid.
func ValidateMultiPlatform(imagePath string) error {
	index, err := ReadIndex(imagePath)
	if err != nil {
		return err
	}

	// OCI images from buildx have a nested structure - resolve it if needed
	if IsNestedIndex(index) {
		index, err = ReadNestedIndex(imagePath, index.Manifests[0].Digest)
		if err != nil {
			return err
		}
	}

	platforms := Platforms(index)

	if len(platforms) == 0 {
		return helpers.Wrap(ErrSinglePlatform, errors.New("image does not specify any platforms"))
	}

	// Check that all required platforms are present
	var missing []string
	for _, required := range RequiredPlatforms() {
		if !platforms[required] {
			missing = append(missing, required)
		}
	}

	if len(missing) > 0 {
		var supported []string
		for p := range platforms {
			supported = append(supported, p)
		}
		return helpers.Wrap(ErrSinglePlatform, fmt.Errorf(
			"image missing required platforms %v (has %v)",
			missing, supported,
		))
	}

	return nil
}
