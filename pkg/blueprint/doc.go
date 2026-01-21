// Package blueprint provides structures for defining system compositions.
//
// A blueprint orchestrates how resources are deployed, declaring which
// services and widgets should be included in a system deployment. It
// serves as the input to the planning phase, where references are resolved
// and a concrete deployment plan is generated.
//
// Load a blueprint:
//
//	bp, err := blueprint.Read("blueprint.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
package blueprint
