package codec

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/cruciblehq/protocol/internal/helpers"
	"github.com/go-viper/mapstructure/v2"
	"gopkg.in/yaml.v3"
)

// Encodes a value to the specified format and writes it to w.
//
// The contentType specifies the output format. The key parameter specifies the
// struct tag to use for field mapping. The indent parameter controls whether
// JSON output should be pretty-printed. The v parameter is the value to be
// encoded. The function writes the encoded data to w or returns an error if
// encoding fails or if the content type is unsupported.
func Encode(w io.Writer, contentType ContentType, key string, indent bool, v any) error {
	raw, err := structToMap(v, key)
	if err != nil {
		return helpers.Wrap(ErrEncodingFailed, err)
	}

	// Encode map to target format
	switch contentType {
	case ContentTypeJSON:
		return encodeJSON(w, indent, raw)
	case ContentTypeYAML:
		return encodeYAML(w, raw)
	case ContentTypeTOML:
		return encodeTOML(w, raw)
	default:
		return ErrUnsupportedContentType
	}
}

// Converts a struct to a map using the specified tag name, handling nested
// structs recursively.
func structToMap(v any, tagName string) (map[string]any, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nonStructToMap(v, tagName)
	}

	return structFieldsToMap(val, tagName)
}

// Converts non-struct types to a map using mapstructure.
//
// This is a fallback for types like maps or primitives that mapstructure
// can handle directly without custom reflection logic.
func nonStructToMap(v any, tagName string) (map[string]any, error) {
	var result map[string]any
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &result,
		TagName: tagName,
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(v); err != nil {
		return nil, err
	}
	return result, nil
}

// Converts struct fields to a map using the specified tag name.
//
// Iterates through exported fields and uses their tag values as map keys.
// Handles omitempty semantics and recursively converts nested types.
func structFieldsToMap(val reflect.Value, tagName string) (map[string]any, error) {
	result := make(map[string]any)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !field.IsExported() {
			continue
		}

		tagValue, omitEmpty := parseFieldTag(field.Tag, tagName)
		if tagValue == "" {
			continue
		}

		if omitEmpty && isEmptyValue(fieldValue) {
			continue
		}

		convertedValue, err := convertValue(fieldValue, tagName)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", field.Name, err)
		}

		result[tagValue] = convertedValue
	}

	return result, nil
}

// Parses a field tag and returns the tag value and omitempty flag.
//
// Extracts the field name from the tag and detects the "omitempty" option.
// Returns empty string if the tag is not present.
func parseFieldTag(tag reflect.StructTag, tagName string) (string, bool) {
	tagValue := tag.Get(tagName)
	if tagValue == "" {
		return "", false
	}

	// Split by comma to handle multiple options
	parts := strings.Split(tagValue, ",")
	if len(parts) == 0 {
		return "", false
	}

	// First part is the field name
	fieldName := strings.TrimSpace(parts[0])

	// Check remaining parts for omitempty
	omitEmpty := false
	for i := 1; i < len(parts); i++ {
		if strings.TrimSpace(parts[i]) == "omitempty" {
			omitEmpty = true
			break
		}
	}

	return fieldName, omitEmpty
}

// Converts a reflect.Value to a suitable type for encoding.
//
// Handles pointers, structs, slices, arrays, and maps by recursively applying
// the tag-based conversion. Primitive types are returned as-is.
func convertValue(v reflect.Value, tagName string) (any, error) {
	if !v.IsValid() {
		return nil, nil
	}

	// Dereference pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		return structToMap(v.Interface(), tagName)

	case reflect.Slice, reflect.Array:
		result := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			elem, err := convertValue(v.Index(i), tagName)
			if err != nil {
				return nil, err
			}
			result[i] = elem
		}
		return result, nil

	case reflect.Map:
		result := make(map[string]any)
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()
			// Convert key to string
			keyStr := fmt.Sprintf("%v", key.Interface())
			convertedValue, err := convertValue(value, tagName)
			if err != nil {
				return nil, err
			}
			result[keyStr] = convertedValue
		}
		return result, nil

	default:
		return v.Interface(), nil
	}
}

// Checks if a reflect.Value is considered empty for omitempty.
//
// Follows Go's JSON omitempty semantics: zero values for numbers, false for
// booleans, empty for strings/slices/maps/arrays, and nil for pointers/interfaces.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// Encodes a value to JSON format.
func encodeJSON(w io.Writer, indent bool, v any) error {
	encoder := json.NewEncoder(w)
	if indent {
		encoder.SetIndent("", "  ")
	}
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

// Decodes a raw map into the target structure.
//
// The raw parameter is a map representing unmarshaled content. The key
// parameter specifies the struct tag to use for field mapping. The target
// parameter is a pointer to the structure where the decoded data should be
// stored. This is useful when you already have a map and need to decode it
// into a struct with custom field tags.
func DecodeMap(raw map[string]any, key string, target any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  target,
		TagName: key,
	})
	if err != nil {
		return helpers.Wrap(ErrDecodingFailed, err)
	}

	if err := decoder.Decode(raw); err != nil {
		return helpers.Wrap(ErrDecodingFailed, err)
	}

	return nil
}
