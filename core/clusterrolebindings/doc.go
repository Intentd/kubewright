// Package clusterrolebindings builds and writes ClusterRoleBindings with
// Build, Create, Apply, Delete and DeleteIfOwned.
//
// A ClusterRoleBinding is cluster-scoped, so DeleteIfOwned takes no namespace.
// A namespaced owner can't hold an owner reference on one: mark it with labels
// when you create it and delete it with types.OwnedByLabels.
package clusterrolebindings
