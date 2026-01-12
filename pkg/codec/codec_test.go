package codec

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type testStruct struct {
	Name    string `key:"name"`
	Version int    `key:"version"`
	Enabled bool   `key:"enabled"`
}

func TestEncode_JSON(t *testing.T) {
	v := testStruct{Name: "test", Version: 1, Enabled: true}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "key", v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if !strings.Contains(data, `"name"`) {
		t.Error("expected JSON to contain 'name' key")
	}
	if !strings.Contains(data, `"test"`) {
		t.Error("expected JSON to contain 'test' value")
	}
}

func TestEncode_YAML(t *testing.T) {
	v := testStruct{Name: "test", Version: 1, Enabled: true}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeYAML, "key", v); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "name:") {
		t.Error("expected YAML to contain 'name' key")
	}
}

func TestEncode_TOML(t *testing.T) {
	v := testStruct{Name: "test", Version: 1, Enabled: true}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeTOML, "key", v); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "name =") {
		t.Error("expected TOML to contain 'name' key")
	}
}

func TestEncode_UnsupportedContentType(t *testing.T) {
	v := testStruct{Name: "test", Version: 1, Enabled: true}

	var buf bytes.Buffer
	err := Encode(&buf, ContentTypeUnknown, "key", v)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrUnsupportedContentType) {
		t.Errorf("expected ErrUnsupportedContentType, got %v", err)
	}
}

func TestEncode_CustomTag(t *testing.T) {
	type customStruct struct {
		Name string `custom:"custom_name"`
	}

	v := customStruct{Name: "test"}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "custom", v); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), `"custom_name"`) {
		t.Error("expected JSON to contain 'custom_name' key")
	}
}

func TestDecode_JSON(t *testing.T) {
	data := `{"name":"test","version":1,"enabled":true}`

	var target testStruct
	if err := Decode(strings.NewReader(data), ContentTypeJSON, "key", &target); err != nil {
		t.Fatal(err)
	}

	if target.Name != "test" {
		t.Errorf("expected name %q, got %q", "test", target.Name)
	}
	if target.Version != 1 {
		t.Errorf("expected version %d, got %d", 1, target.Version)
	}
	if !target.Enabled {
		t.Error("expected enabled to be true")
	}
}

func TestDecode_YAML(t *testing.T) {
	data := "name: test\nversion: 1\nenabled: true\n"

	var target testStruct
	if err := Decode(strings.NewReader(data), ContentTypeYAML, "key", &target); err != nil {
		t.Fatal(err)
	}

	if target.Name != "test" {
		t.Errorf("expected name %q, got %q", "test", target.Name)
	}
	if target.Version != 1 {
		t.Errorf("expected version %d, got %d", 1, target.Version)
	}
	if !target.Enabled {
		t.Error("expected enabled to be true")
	}
}

func TestDecode_TOML(t *testing.T) {
	data := "name = \"test\"\nversion = 1\nenabled = true\n"

	var target testStruct
	if err := Decode(strings.NewReader(data), ContentTypeTOML, "key", &target); err != nil {
		t.Fatal(err)
	}

	if target.Name != "test" {
		t.Errorf("expected name %q, got %q", "test", target.Name)
	}
	if target.Version != 1 {
		t.Errorf("expected version %d, got %d", 1, target.Version)
	}
	if !target.Enabled {
		t.Error("expected enabled to be true")
	}
}

func TestDecodeMap(t *testing.T) {
	raw := map[string]any{
		"name":    "test",
		"version": 1,
		"enabled": true,
	}

	var target testStruct
	if err := DecodeMap(raw, "key", &target); err != nil {
		t.Fatal(err)
	}

	if target.Name != "test" {
		t.Errorf("expected name %q, got %q", "test", target.Name)
	}
	if target.Version != 1 {
		t.Errorf("expected version %d, got %d", 1, target.Version)
	}
	if !target.Enabled {
		t.Error("expected enabled to be true")
	}
}

func TestDecodeMap_NestedStruct(t *testing.T) {
	type Inner struct {
		Value string `field:"inner_value"`
	}

	type Outer struct {
		Title string `field:"title"`
		Data  Inner  `field:"data"`
	}

	raw := map[string]any{
		"title": "test",
		"data": map[string]any{
			"inner_value": "nested",
		},
	}

	var target Outer
	if err := DecodeMap(raw, "field", &target); err != nil {
		t.Fatal(err)
	}

	if target.Title != "test" {
		t.Errorf("expected title %q, got %q", "test", target.Title)
	}
	if target.Data.Value != "nested" {
		t.Errorf("expected inner value %q, got %q", "nested", target.Data.Value)
	}
}

func TestDecode_UnsupportedContentType(t *testing.T) {
	data := `{"name":"test"}`

	var target testStruct
	err := Decode(strings.NewReader(data), ContentTypeUnknown, "key", &target)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrUnsupportedContentType) {
		t.Errorf("expected ErrUnsupportedContentType, got %v", err)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	data := `{invalid}`

	var target testStruct
	err := Decode(strings.NewReader(data), ContentTypeJSON, "key", &target)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrDecodingFailed) {
		t.Errorf("expected ErrDecodingFailed, got %v", err)
	}
}

func TestDecode_InvalidYAML(t *testing.T) {
	data := ":\ninvalid"

	var target testStruct
	err := Decode(strings.NewReader(data), ContentTypeYAML, "key", &target)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrDecodingFailed) {
		t.Errorf("expected ErrDecodingFailed, got %v", err)
	}
}

func TestDecode_InvalidTOML(t *testing.T) {
	data := "= invalid"

	var target testStruct
	err := Decode(strings.NewReader(data), ContentTypeTOML, "key", &target)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrDecodingFailed) {
		t.Errorf("expected ErrDecodingFailed, got %v", err)
	}
}

func TestDecode_CustomTag(t *testing.T) {
	type customStruct struct {
		Name string `custom:"custom_name"`
	}

	data := `{"custom_name":"test"}`

	var target customStruct
	if err := Decode(strings.NewReader(data), ContentTypeJSON, "custom", &target); err != nil {
		t.Fatal(err)
	}

	if target.Name != "test" {
		t.Errorf("expected name %q, got %q", "test", target.Name)
	}
}

func TestRoundtrip_JSON(t *testing.T) {
	original := testStruct{Name: "test", Version: 42, Enabled: true}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "key", original); err != nil {
		t.Fatal(err)
	}

	var decoded testStruct
	if err := Decode(&buf, ContentTypeJSON, "key", &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Name != original.Name {
		t.Errorf("expected name %q, got %q", original.Name, decoded.Name)
	}
	if decoded.Version != original.Version {
		t.Errorf("expected version %d, got %d", original.Version, decoded.Version)
	}
	if decoded.Enabled != original.Enabled {
		t.Errorf("expected enabled %v, got %v", original.Enabled, decoded.Enabled)
	}
}

func TestRoundtrip_YAML(t *testing.T) {
	original := testStruct{Name: "test", Version: 42, Enabled: true}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeYAML, "key", original); err != nil {
		t.Fatal(err)
	}

	var decoded testStruct
	if err := Decode(&buf, ContentTypeYAML, "key", &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Name != original.Name {
		t.Errorf("expected name %q, got %q", original.Name, decoded.Name)
	}
	if decoded.Version != original.Version {
		t.Errorf("expected version %d, got %d", original.Version, decoded.Version)
	}
	if decoded.Enabled != original.Enabled {
		t.Errorf("expected enabled %v, got %v", original.Enabled, decoded.Enabled)
	}
}

func TestRoundtrip_TOML(t *testing.T) {
	original := testStruct{Name: "test", Version: 42, Enabled: true}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeTOML, "key", original); err != nil {
		t.Fatal(err)
	}

	var decoded testStruct
	if err := Decode(&buf, ContentTypeTOML, "key", &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Name != original.Name {
		t.Errorf("expected name %q, got %q", original.Name, decoded.Name)
	}
	if decoded.Version != original.Version {
		t.Errorf("expected version %d, got %d", original.Version, decoded.Version)
	}
	if decoded.Enabled != original.Enabled {
		t.Errorf("expected enabled %v, got %v", original.Enabled, decoded.Enabled)
	}
}
