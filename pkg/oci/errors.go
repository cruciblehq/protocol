package oci

import "errors"

var (
	ErrInvalidImage   = errors.New("invalid OCI image")
	ErrSinglePlatform = errors.New("single-platform image")
)
