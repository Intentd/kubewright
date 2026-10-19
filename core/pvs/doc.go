// Package pvs builds and writes PersistentVolumes with Build, Create, Apply,
// Delete and DeleteIfOwned.
//
// A PersistentVolume is cluster-scoped, so DeleteIfOwned takes no namespace. A
// namespaced owner can't hold an owner reference on one: mark it with labels
// when you create it and delete it with types.OwnedByLabels.
package pvs
