// Package statefulsets builds and writes StatefulSets with Build, Create,
// Apply, Delete and DeleteIfOwned.
//
// Its flat container fields describe one container. Containers and
// InitContainers describe more; types.Container says how the two combine.
//
// Volume claim templates are mounted on the first container.
package statefulsets
