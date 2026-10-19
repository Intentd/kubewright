// Package secrets builds and writes Secrets with Build, Create, Apply, Delete
// and DeleteIfOwned.
//
// CopyToNamespaces copies a Secret into each of a list of namespaces. It reads
// the source once and skips namespaces that already have the copy.
package secrets
