// Package limitranges builds and writes LimitRanges with Build, Create, Apply,
// Delete and DeleteIfOwned.
//
// CopyToNamespaces copies a LimitRange into each of a list of namespaces. It
// reads the source once and skips namespaces that already have the copy.
package limitranges
