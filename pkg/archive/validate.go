package archive

import (
	"os"
	"path/filepath"

	"github.com/cruciblehq/protocol/pkg/manifest"
)

const (

	// The required main file for widgets.
	WidgetMainFile = "index.js"

	// The required image file for services.
	ServiceImageFile = "image.tar"
)

// Checks that a widget's dist/ directory contains required files.
func ValidateWidgetStructure(distDir string, m *manifest.Widget) error {
	widgetMain := filepath.Join(distDir, WidgetMainFile)
	if _, err := os.Stat(widgetMain); os.IsNotExist(err) {
		return ErrInvalidStructure
	}
	return nil
}

// Checks that a service's dist/ directory contains required files.
func ValidateServiceStructure(distDir string, m *manifest.Service) error {
	serviceImage := filepath.Join(distDir, ServiceImageFile)
	if _, err := os.Stat(serviceImage); os.IsNotExist(err) {
		return ErrInvalidStructure
	}
	return nil
}
