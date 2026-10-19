// Package pods builds and writes Pods with Build, Create, Apply, Delete and
// DeleteIfOwned.
//
// Its flat container fields describe one container. Containers and
// InitContainers describe more; types.Container says how the two combine.
//
// RestartPolicy is Always, OnFailure or Never; empty leaves the Kubernetes
// default, Always.
package pods
