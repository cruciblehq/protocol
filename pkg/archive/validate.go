package archive

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cruciblehq/protocol/pkg/manifest"
)

var (
	ErrInvalidStructure = errors.New("invalid resource structure")
)

// Checks that a widget's dist/ directory contains required files.
func ValidateWidgetStructure(distDir string, m *manifest.Widget) error {
	widgetMain := filepath.Join(distDir, "index.js")
	if _, err := os.Stat(widgetMain); os.IsNotExist(err) {
		return ErrInvalidStructure
	}
	return nil
}

// Checks that a service's dist/ directory contains required files.
func ValidateServiceStructure(distDir string, m *manifest.Service) error {
	serviceImage := filepath.Join(distDir, "image.tar")
	if _, err := os.Stat(serviceImage); os.IsNotExist(err) {
		return ErrInvalidStructure
	}
	return nil
}
