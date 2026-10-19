// Package configmaps builds and writes ConfigMaps with Build, Create, Apply,
// Delete and DeleteIfOwned.
//
// CopyToNamespaces copies a ConfigMap into each of a list of namespaces. It
// reads the source once and skips namespaces that already have the copy.
package configmaps
