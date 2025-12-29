package codec

import (
	"encoding/json"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/cruciblehq/protocol/internal/helpers"
	"github.com/go-viper/mapstructure/v2"
	"gopkg.in/yaml.v3"
)

// Encodes a value to the specified format and writes it to w.
//
// The contentType specifies the output format. The key parameter specifies the
// struct tag to use for field mapping. The v parameter is the value to be
// encoded. The function writes the encoded data to w or returns an error if
// encoding fails or if the content type is unsupported.
func Encode(w io.Writer, contentType ContentType, key string, v any) error {

	// Struct to map
	var raw map[string]any
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &raw,
		TagName: key,
	})
	if err != nil {
		return helpers.Wrap(ErrEncodingFailed, err)
	}

	// Decode
	if err := decoder.Decode(v); err != nil {
		return helpers.Wrap(ErrEncodingFailed, err)
	}

	// Map to the target format
	switch contentType {
	case ContentTypeJSON:
		return encodeJSON(w, raw)
	case ContentTypeYAML:
		return encodeYAML(w, raw)
	case ContentTypeTOML:
		return encodeTOML(w, raw)
	default:
		return ErrUnsupportedContentType
	}
}

// Encodes a value to JSON format.
func encodeJSON(w io.Writer, v any) error {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		return helpers.Wrap(ErrEncodingFailed, err)
	}
	return nil
}

// Encodes a value to YAML format.
func encodeYAML(w io.Writer, v any) error {
	encoder := yaml.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		return helpers.Wrap(ErrEncodingFailed, err)
	}
	return nil
}

// Encodes a value to TOML format.
func encodeTOML(w io.Writer, v any) error {
	encoder := toml.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		return helpers.Wrap(ErrEncodingFailed, err)
	}
	return nil
}

// Decodes data in the specified format into the target.
//
// The contentType specifies the format of the input data. The key parameter
// specifies the struct tag to use for field mapping. The target parameter is a
// pointer to the structure where the decoded data should be stored. The r
// parameter is the reader from which to read the data.
func Decode(r io.Reader, contentType ContentType, key string, target any) error {
	var raw map[string]any

	switch contentType {
	case ContentTypeJSON:
		decoder := json.NewDecoder(r)
		if err := decoder.Decode(&raw); err != nil {
			return helpers.Wrap(ErrDecodingFailed, err)
		}
	case ContentTypeYAML:
		decoder := yaml.NewDecoder(r)
		if err := decoder.Decode(&raw); err != nil {
			return helpers.Wrap(ErrDecodingFailed, err)
		}
	case ContentTypeTOML:
		decoder := toml.NewDecoder(r)
		if _, err := decoder.Decode(&raw); err != nil {
			return helpers.Wrap(ErrDecodingFailed, err)
		}
	default:
		return ErrUnsupportedContentType
	}

	mapDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  target,
		TagName: key,
	})
	if err != nil {
		return helpers.Wrap(ErrDecodingFailed, err)
	}

	if err := mapDecoder.Decode(raw); err != nil {
		return helpers.Wrap(ErrDecodingFailed, err)
	}

	return nil
}
