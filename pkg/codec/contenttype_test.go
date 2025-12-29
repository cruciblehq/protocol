package codec

import (
	"errors"
	"testing"
)

func TestParse_ApplicationJSON(t *testing.T) {
	ct, _, err := Parse("application/json")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON, got %v", ct)
	}
}

func TestParse_ApplicationJSONWithCharset(t *testing.T) {
	ct, _, err := Parse("application/json; charset=utf-8")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON, got %v", ct)
	}
}

func TestParse_TextJSON(t *testing.T) {
	ct, _, err := Parse("text/json")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON, got %v", ct)
	}
}

func TestParse_ApplicationYAML(t *testing.T) {
	ct, _, err := Parse("application/yaml")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeYAML {
		t.Errorf("expected ContentTypeYAML, got %v", ct)
	}
}

func TestParse_ApplicationXYAML(t *testing.T) {
	ct, _, err := Parse("application/x-yaml")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeYAML {
		t.Errorf("expected ContentTypeYAML, got %v", ct)
	}
}

func TestParse_TextYAML(t *testing.T) {
	ct, _, err := Parse("text/yaml")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeYAML {
		t.Errorf("expected ContentTypeYAML, got %v", ct)
	}
}

func TestParse_ApplicationTOML(t *testing.T) {
	ct, _, err := Parse("application/toml")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeTOML {
		t.Errorf("expected ContentTypeTOML, got %v", ct)
	}
}

func TestParse_ApplicationXTOML(t *testing.T) {
	ct, _, err := Parse("application/x-toml")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeTOML {
		t.Errorf("expected ContentTypeTOML, got %v", ct)
	}
}

func TestParse_TextTOML(t *testing.T) {
	ct, _, err := Parse("text/toml")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeTOML {
		t.Errorf("expected ContentTypeTOML, got %v", ct)
	}
}

func TestParse_Unsupported(t *testing.T) {
	_, _, err := Parse("application/xml")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrUnsupportedContentType) {
		t.Errorf("expected ErrUnsupportedContentType, got %v", err)
	}
}

func TestParse_Invalid(t *testing.T) {
	_, _, err := Parse("not a media type")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrInvalidContentType) {
		t.Errorf("expected ErrInvalidContentType, got %v", err)
	}
}

func TestParse_Empty(t *testing.T) {
	_, _, err := Parse("")
	if err == nil {
		t.Fatal("expected error")
	}
}
func TestParse_VendorJSON(t *testing.T) {
	ct, _, err := Parse("application/vnd.example+json")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON, got %v", ct)
	}
}

func TestParse_VendorYAML(t *testing.T) {
	ct, _, err := Parse("application/vnd.crucible.namespace.v0+yaml")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeYAML {
		t.Errorf("expected ContentTypeYAML, got %v", ct)
	}
}

func TestParse_VendorTOML(t *testing.T) {
	ct, _, err := Parse("application/vnd.example.resource.v1+toml")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeTOML {
		t.Errorf("expected ContentTypeTOML, got %v", ct)
	}
}

func TestParse_VendorJSONWithCharset(t *testing.T) {
	ct, _, err := Parse("application/vnd.example+json; charset=utf-8")
	if err != nil {
		t.Fatal(err)
	}
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON, got %v", ct)
	}
}

func TestParse_VendorUnsupportedSuffix(t *testing.T) {
	_, _, err := Parse("application/vnd.example+xml")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrUnsupportedContentType) {
		t.Errorf("expected ErrUnsupportedContentType, got %v", err)
	}
}

func TestNegotiate_Empty(t *testing.T) {
	ct := Negotiate("")
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON for empty Accept, got %v", ct)
	}
}

func TestNegotiate_Wildcard(t *testing.T) {
	ct := Negotiate("*/*")
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON for wildcard Accept, got %v", ct)
	}
}

func TestNegotiate_JSON(t *testing.T) {
	ct := Negotiate("application/json")
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON, got %v", ct)
	}
}

func TestNegotiate_YAML(t *testing.T) {
	ct := Negotiate("application/yaml")
	if ct != ContentTypeYAML {
		t.Errorf("expected ContentTypeYAML, got %v", ct)
	}
}

func TestNegotiate_TOML(t *testing.T) {
	ct := Negotiate("application/toml")
	if ct != ContentTypeTOML {
		t.Errorf("expected ContentTypeTOML, got %v", ct)
	}
}

func TestNegotiate_VendorJSON(t *testing.T) {
	ct := Negotiate("application/vnd.crucible.resource.v1+json")
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON, got %v", ct)
	}
}

func TestNegotiate_VendorYAML(t *testing.T) {
	ct := Negotiate("application/vnd.example+yaml")
	if ct != ContentTypeYAML {
		t.Errorf("expected ContentTypeYAML, got %v", ct)
	}
}

func TestNegotiate_Unsupported(t *testing.T) {
	ct := Negotiate("application/xml")
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON for unsupported type, got %v", ct)
	}
}

func TestNegotiate_Invalid(t *testing.T) {
	ct := Negotiate("not a valid media type")
	if ct != ContentTypeJSON {
		t.Errorf("expected ContentTypeJSON for invalid type, got %v", ct)
	}
}

