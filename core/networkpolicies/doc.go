// Package networkpolicies builds and writes NetworkPolicies with Build,
// Create, Apply, Delete and DeleteIfOwned.
//
// CopyToNamespaces copies a NetworkPolicy into each of a list of namespaces.
// It reads the source once and skips namespaces that already have the copy.
package networkpolicies
