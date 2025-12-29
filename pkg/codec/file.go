package codec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Encodes a value to a file.
//
// The content type is inferred from the file extension. The key parameter
// specifies the struct tag to use for field mapping. The v parameter is the
// value to be encoded. The path parameter is the file path to write to.
func EncodeFile(path, key string, v any) error {
	ct, err := contentTypeFromExtension(path)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrEncodingFailed, err)
	}
	defer f.Close()

	if err := Encode(f, ct, key, v); err != nil {
		return err
	}

	return nil
}

// Decodes a file into the target.
//
// The content type is inferred from the file extension. The key parameter
// specifies the struct tag to use for field mapping. The target parameter is a
// pointer to the structure where the decoded data should be stored. The path
// parameter is the file path to read from.
func DecodeFile(path, key string, target any) (ContentType, error) {
	ct, err := contentTypeFromExtension(path)
	if err != nil {
		return ContentTypeUnknown, err
	}

	f, err := os.Open(path)
	if err != nil {
		return ContentTypeUnknown, fmt.Errorf("%w: %v", ErrDecodingFailed, err)
	}
	defer f.Close()

	return ct, Decode(f, ct, key, target)
}

// Returns the content type for the file extension.
func contentTypeFromExtension(path string) (ContentType, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return ContentTypeJSON, nil
	case ".yaml", ".yml":
		return ContentTypeYAML, nil
	case ".toml":
		return ContentTypeTOML, nil
	default:
		return ContentTypeUnknown, fmt.Errorf("%w: unknown extension %q", ErrUnsupportedContentType, ext)
	}
}
