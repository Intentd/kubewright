// Package resourcequotas builds and writes ResourceQuotas with Build, Create,
// Apply, Delete and DeleteIfOwned.
//
// CopyToNamespaces copies a ResourceQuota into each of a list of namespaces.
// It reads the source once and skips namespaces that already have the copy.
package resourcequotas
