// Package cronjobs builds and writes CronJobs with Build, Create, Apply,
// Delete and DeleteIfOwned.
//
// Its flat container fields describe one container. Containers and
// InitContainers describe more; types.Container says how the two combine.
//
// RestartPolicy is OnFailure, the default, or Never, which starts a new pod
// for each retry and keeps the failed ones.
package cronjobs
