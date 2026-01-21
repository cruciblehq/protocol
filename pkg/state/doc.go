// Package state provides structures for tracking deployment state.
//
// State records what resources have been deployed, their runtime identifiers,
// and their current operational status. It is used for incremental deployments
// and resource lifecycle management.
//
// State enables incremental deployments by comparing desired state (from a
// plan) with current state to determine what changes need to be applied.
//
// Example usage:
//
//	// Read existing state
//	st, err := state.Read("state.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Save updated state
//	st.DeployedAt = time.Now()
//	err = st.Write("state.json")
//	if err != nil {
//		log.Fatal(err)
//	}
package state
