package shared

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
)

// SyncContainerSpec copies the fields desired sets onto existing and reports
// whether anything changed. A zero or nil field in desired is left alone:
// Kubernetes may have defaulted it, and overwriting a default with an empty
// value would show as a change on every Apply.
func SyncContainerSpec(existing *corev1.Container, desired corev1.Container) bool {
	drifted := false

	// Image — always sync; an image is always set.
	if existing.Image != desired.Image {
		existing.Image = desired.Image
		drifted = true
	}

	// Ports — only sync when desired sets ports.
	// Kubernetes defaults Protocol to "TCP"; buildDeployment sets it explicitly
	// so desired and existing match on the first reconcile after create.
	if len(desired.Ports) > 0 && !reflect.DeepEqual(existing.Ports, desired.Ports) {
		existing.Ports = desired.Ports
		drifted = true
	}

	// Env — sync when desired declares env vars.
	// Also sync when existing has env but desired does not, so explicitly-cleared
	// env vars are removed (e.g., after removing an env block from the spec).
	if !reflect.DeepEqual(existing.Env, desired.Env) && (len(desired.Env) > 0 || len(existing.Env) > 0) {
		existing.Env = desired.Env
		drifted = true
	}

	// EnvFrom — only sync when desired sets envFrom sources.
	if len(desired.EnvFrom) > 0 && !reflect.DeepEqual(existing.EnvFrom, desired.EnvFrom) {
		existing.EnvFrom = desired.EnvFrom
		drifted = true
	}

	// VolumeMounts — only sync when desired sets mounts.
	if len(desired.VolumeMounts) > 0 && !reflect.DeepEqual(existing.VolumeMounts, desired.VolumeMounts) {
		existing.VolumeMounts = desired.VolumeMounts
		drifted = true
	}

	// Probes — only sync when desired sets them.
	// Kubernetes does not add default probes — nil in desired means not set.
	if desired.LivenessProbe != nil && !reflect.DeepEqual(existing.LivenessProbe, desired.LivenessProbe) {
		existing.LivenessProbe = desired.LivenessProbe
		drifted = true
	}
	if desired.ReadinessProbe != nil && !reflect.DeepEqual(existing.ReadinessProbe, desired.ReadinessProbe) {
		existing.ReadinessProbe = desired.ReadinessProbe
		drifted = true
	}
	if desired.StartupProbe != nil && !reflect.DeepEqual(existing.StartupProbe, desired.StartupProbe) {
		existing.StartupProbe = desired.StartupProbe
		drifted = true
	}

	// SecurityContext — only sync when desired sets it.
	if desired.SecurityContext != nil && !reflect.DeepEqual(existing.SecurityContext, desired.SecurityContext) {
		existing.SecurityContext = desired.SecurityContext
		drifted = true
	}

	// WorkingDir — only sync when desired sets it.
	if desired.WorkingDir != "" && existing.WorkingDir != desired.WorkingDir {
		existing.WorkingDir = desired.WorkingDir
		drifted = true
	}

	return drifted
}

// SyncPodSpec is SyncContainerSpec for pod-level fields.
func SyncPodSpec(existing *corev1.PodSpec, desired corev1.PodSpec) bool {
	drifted := false

	// Volumes — only sync when desired sets volumes.
	if len(desired.Volumes) > 0 && !reflect.DeepEqual(existing.Volumes, desired.Volumes) {
		existing.Volumes = desired.Volumes
		drifted = true
	}

	// SecurityContext — only sync when desired sets it.
	// Kubernetes does not default this field.
	if desired.SecurityContext != nil && !reflect.DeepEqual(existing.SecurityContext, desired.SecurityContext) {
		existing.SecurityContext = desired.SecurityContext
		drifted = true
	}

	// ServiceAccountName — only sync when desired sets one.
	// Kubernetes always defaults this to "default". Comparing an empty desired
	// against "default" would trigger an update on every Apply.
	if desired.ServiceAccountName != "" && existing.ServiceAccountName != desired.ServiceAccountName {
		existing.ServiceAccountName = desired.ServiceAccountName
		drifted = true
	}

	// NodeSelector — only sync when desired sets selectors.
	if len(desired.NodeSelector) > 0 && !reflect.DeepEqual(existing.NodeSelector, desired.NodeSelector) {
		existing.NodeSelector = desired.NodeSelector
		drifted = true
	}

	return drifted
}
