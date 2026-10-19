// Package types holds the field types a resource spec is made of: env and
// envFrom, resource requirements, probes, pod and container security,
// volumes, rolling update, HPA behavior, network policy rules, limit range
// items, labels and custom resource metadata. Each kind's Spec uses them;
// their yaml and json tags let a program decode them from a file.
//
// It also holds what every kind shares: Owner, the object resources are built
// for; Options, how Apply writes; and Result, what a verb did.
package types
