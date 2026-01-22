package resource

// Type represents the type of a Crucible resource.
type Type string

const (
	TypeService  Type = "service"  // Service resource type.
	TypeTemplate Type = "template" // Template resource type.
	TypeWidget   Type = "widget"   // Widget resource type.
)
