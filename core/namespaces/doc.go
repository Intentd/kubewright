// Package namespaces builds and writes Namespaces with Build, Create, Apply,
// Delete and DeleteIfOwned.
//
// A Namespace is cluster-scoped, so DeleteIfOwned takes no namespace. A
// namespaced owner can't hold an owner reference on one: mark it with labels
// when you create it and delete it with types.OwnedByLabels.
package namespaces
