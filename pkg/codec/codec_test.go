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
	if err := Encode(&buf, ContentTypeJSON, "key", false, v); err != nil {
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
	if err := Encode(&buf, ContentTypeYAML, "key", false, v); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "name:") {
		t.Error("expected YAML to contain 'name' key")
	}
}

func TestEncode_TOML(t *testing.T) {
	v := testStruct{Name: "test", Version: 1, Enabled: true}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeTOML, "key", false, v); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "name =") {
		t.Error("expected TOML to contain 'name' key")
	}
}

func TestEncode_UnsupportedContentType(t *testing.T) {
	v := testStruct{Name: "test", Version: 1, Enabled: true}

	var buf bytes.Buffer
	err := Encode(&buf, ContentTypeUnknown, "key", false, v)
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
	if err := Encode(&buf, ContentTypeJSON, "custom", false, v); err != nil {
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
	if err := Encode(&buf, ContentTypeJSON, "key", false, original); err != nil {
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
	if err := Encode(&buf, ContentTypeYAML, "key", false, original); err != nil {
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
	if err := Encode(&buf, ContentTypeTOML, "key", false, original); err != nil {
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

func TestEncode_NestedStruct(t *testing.T) {
	type Inner struct {
		Value string `field:"inner_value"`
	}
	type Outer struct {
		Title string `field:"title"`
		Data  Inner  `field:"data"`
	}

	v := Outer{
		Title: "test",
		Data:  Inner{Value: "nested"},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if !strings.Contains(data, `"title"`) {
		t.Error("expected JSON to contain 'title' key")
	}
	if !strings.Contains(data, `"data"`) {
		t.Error("expected JSON to contain 'data' key")
	}
	if !strings.Contains(data, `"inner_value"`) {
		t.Error("expected JSON to contain 'inner_value' key")
	}
	if !strings.Contains(data, `"nested"`) {
		t.Error("expected JSON to contain 'nested' value")
	}
}

func TestEncode_SliceOfStructs(t *testing.T) {
	type Item struct {
		ID   string `field:"id"`
		Name string `field:"name"`
	}
	type Container struct {
		Items []Item `field:"items"`
	}

	v := Container{
		Items: []Item{
			{ID: "1", Name: "first"},
			{ID: "2", Name: "second"},
		},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if !strings.Contains(data, `"items"`) {
		t.Error("expected JSON to contain 'items' key")
	}
	if !strings.Contains(data, `"id"`) {
		t.Error("expected JSON to contain 'id' key")
	}
	if !strings.Contains(data, `"first"`) {
		t.Error("expected JSON to contain 'first' value")
	}
	if !strings.Contains(data, `"second"`) {
		t.Error("expected JSON to contain 'second' value")
	}
}

func TestEncode_Map(t *testing.T) {
	type Config struct {
		Settings map[string]string `field:"settings"`
	}

	v := Config{
		Settings: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if !strings.Contains(data, `"settings"`) {
		t.Error("expected JSON to contain 'settings' key")
	}
	if !strings.Contains(data, `"key1"`) {
		t.Error("expected JSON to contain 'key1' key")
	}
	if !strings.Contains(data, `"value1"`) {
		t.Error("expected JSON to contain 'value1' value")
	}
}

func TestEncode_Pointer(t *testing.T) {
	type Data struct {
		Value *string `field:"value"`
	}

	str := "test"
	v := Data{Value: &str}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if !strings.Contains(data, `"value"`) {
		t.Error("expected JSON to contain 'value' key")
	}
	if !strings.Contains(data, `"test"`) {
		t.Error("expected JSON to contain 'test' value")
	}
}

func TestEncode_NilPointer(t *testing.T) {
	type Data struct {
		Value *string `field:"value,omitempty"`
	}

	v := Data{Value: nil}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if strings.Contains(data, `"value"`) {
		t.Error("expected JSON to not contain 'value' key due to omitempty")
	}
}

func TestEncode_OmitEmpty(t *testing.T) {
	type Data struct {
		Name    string `field:"name"`
		Empty   string `field:"empty,omitempty"`
		Zero    int    `field:"zero,omitempty"`
		Enabled bool   `field:"enabled,omitempty"`
	}

	v := Data{Name: "test"}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if !strings.Contains(data, `"name"`) {
		t.Error("expected JSON to contain 'name' key")
	}
	if strings.Contains(data, `"empty"`) {
		t.Error("expected JSON to not contain 'empty' key due to omitempty")
	}
	if strings.Contains(data, `"zero"`) {
		t.Error("expected JSON to not contain 'zero' key due to omitempty")
	}
	if strings.Contains(data, `"enabled"`) {
		t.Error("expected JSON to not contain 'enabled' key due to omitempty")
	}
}

func TestEncode_NonStructMap(t *testing.T) {
	v := map[string]any{
		"key": "value",
	}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if !strings.Contains(data, `"key"`) {
		t.Error("expected JSON to contain 'key'")
	}
	if !strings.Contains(data, `"value"`) {
		t.Error("expected JSON to contain 'value'")
	}
}

func TestEncode_EmptySlice(t *testing.T) {
	type Data struct {
		Items []string `field:"items,omitempty"`
	}

	v := Data{Items: []string{}}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if strings.Contains(data, `"items"`) {
		t.Error("expected JSON to not contain 'items' key due to omitempty")
	}
}

func TestEncode_Array(t *testing.T) {
	type Data struct {
		Values [3]int `field:"values"`
	}

	v := Data{Values: [3]int{1, 2, 3}}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if !strings.Contains(data, `"values"`) {
		t.Error("expected JSON to contain 'values' key")
	}
	if !strings.Contains(data, `[1,2,3]`) {
		t.Error("expected JSON to contain array values")
	}
}

func TestEncode_NestedMap(t *testing.T) {
	type Data struct {
		Config map[string]map[string]int `field:"config"`
	}

	v := Data{
		Config: map[string]map[string]int{
			"section": {"value": 42},
		},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if !strings.Contains(data, `"config"`) {
		t.Error("expected JSON to contain 'config' key")
	}
	if !strings.Contains(data, `"section"`) {
		t.Error("expected JSON to contain 'section' key")
	}
	if !strings.Contains(data, `42`) {
		t.Error("expected JSON to contain value 42")
	}
}

func TestEncode_PointerToStruct(t *testing.T) {
	type Inner struct {
		Value string `field:"value"`
	}
	type Outer struct {
		Data *Inner `field:"data"`
	}

	v := Outer{
		Data: &Inner{Value: "test"},
	}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if !strings.Contains(data, `"data"`) {
		t.Error("expected JSON to contain 'data' key")
	}
	if !strings.Contains(data, `"value"`) {
		t.Error("expected JSON to contain 'value' key")
	}
	if !strings.Contains(data, `"test"`) {
		t.Error("expected JSON to contain 'test' value")
	}
}

func TestEncode_FloatTypes(t *testing.T) {
	type Data struct {
		Float32 float32 `field:"float32,omitempty"`
		Float64 float64 `field:"float64,omitempty"`
	}

	v := Data{Float32: 0.0, Float64: 0.0}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if strings.Contains(data, `"float32"`) {
		t.Error("expected JSON to not contain 'float32' due to omitempty")
	}
	if strings.Contains(data, `"float64"`) {
		t.Error("expected JSON to not contain 'float64' due to omitempty")
	}
}

func TestEncode_IntTypes(t *testing.T) {
	type Data struct {
		Int8   int8   `field:"int8,omitempty"`
		Int16  int16  `field:"int16,omitempty"`
		Int32  int32  `field:"int32,omitempty"`
		Int64  int64  `field:"int64,omitempty"`
		Uint8  uint8  `field:"uint8,omitempty"`
		Uint16 uint16 `field:"uint16,omitempty"`
		Uint32 uint32 `field:"uint32,omitempty"`
		Uint64 uint64 `field:"uint64,omitempty"`
	}

	v := Data{}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if strings.Contains(data, `"int8"`) || strings.Contains(data, `"int16"`) {
		t.Error("expected JSON to not contain zero int fields due to omitempty")
	}
	if strings.Contains(data, `"uint8"`) || strings.Contains(data, `"uint16"`) {
		t.Error("expected JSON to not contain zero uint fields due to omitempty")
	}
}

func TestEncode_EmptyMapOmitted(t *testing.T) {
	type Data struct {
		Values map[string]string `field:"values,omitempty"`
	}

	v := Data{Values: map[string]string{}}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if strings.Contains(data, `"values"`) {
		t.Error("expected JSON to not contain 'values' due to omitempty with empty map")
	}
}

func TestEncode_InterfaceValue(t *testing.T) {
	type Data struct {
		Value any `field:"value,omitempty"`
	}

	v := Data{Value: nil}

	var buf bytes.Buffer
	if err := Encode(&buf, ContentTypeJSON, "field", false, v); err != nil {
		t.Fatal(err)
	}

	data := buf.String()
	if strings.Contains(data, `"value"`) {
		t.Error("expected JSON to not contain 'value' due to omitempty with nil interface")
	}
}
