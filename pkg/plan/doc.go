// Package plan provides structures for representing resolved deployment plans.
//
// A plan is the result of resolving all references in a blueprint against
// available resources and their versions. It contains concrete deployment
// configuration ready for execution by a deployment provider.
//
// Plans support incremental deployments by comparing against previous
// deployment state to determine what resources need to be added, updated,
// or removed.
//
// Example usage:
//
//	// Read an existing plan
//	p, err := plan.Read("plan.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Save a plan
//	err = p.Write("output.json")
//	if err != nil {
//		log.Fatal(err)
//	}
package plan
