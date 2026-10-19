// Package jobs builds and writes Jobs with Build, Create, Delete and
// CreateForDelete.
//
// There is no Apply: a Job's pod template can't change after it is created,
// so a changed spec needs a new Job. Delete removes the Job's pods with it
// (foreground propagation).
//
// CreateForDelete creates a Job with no owner reference, for cleanup work
// that runs while the owner is being deleted and must outlive it.
//
// Its flat container fields describe one container. Containers and
// InitContainers describe more; types.Container says how the two combine.
//
// RestartPolicy is OnFailure, the default, or Never, which starts a new pod
// for each retry and keeps the failed ones.
package jobs
