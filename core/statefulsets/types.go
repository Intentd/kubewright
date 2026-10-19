package statefulsets

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/intentd/kubewright/core/types"
)

// VolumeClaimTemplate is a claim each of a StatefulSet's pods gets.
type VolumeClaimTemplate struct {
	// Name — the claim and volume name. Default: data.
	Name string
	// StorageClass — the storage class each pod's claim provisions from.
	StorageClass string
	// StorageSize — the size of each pod's claim, e.g. "10Gi". Required.
	StorageSize string
	// MountPath — where the volume is mounted in the first container.
	// Default: /data.
	MountPath string
	// AccessModes — ReadWriteOnce (the default), ReadOnlyMany, ReadWriteMany or
	// ReadWriteOncePod.
	AccessModes []string
}

// Spec describes a StatefulSet.
type Spec struct {
	// Name — the StatefulSet's name, and the prefix of its pods' names. Required.
	Name string
	// Namespace — where to create it. Empty means the owner's.
	Namespace string
	// Image — the container image, when the flat fields describe the one
	// container.
	Image string
	// Replicas — how many pods to run.
	Replicas int32
	// Port — the port the container listens on; probes without a port use it.
	Port int32
	// Protocol — for Port: TCP (the default), UDP or SCTP.
	Protocol corev1.Protocol
	// ServiceName — the headless Service that gives the pods stable network
	// names.
	ServiceName string

	// VolumeClaimTemplates — the claims each pod gets.
	VolumeClaimTemplates []VolumeClaimTemplate

	// Labels — set on the object and the pod template.
	Labels map[string]string
	// Annotations — set on the object.
	Annotations map[string]string

	// Env — environment variables for the container, or defaults for every
	// container when Containers is set.
	Env []types.EnvVar
	// EnvFrom — Secrets and ConfigMaps whose keys become environment variables.
	EnvFrom *types.EnvFrom
	// Resources — the container's requests and limits.
	Resources *types.ResourceRequirements

	// Selector is the label selector for the pods, set on the pod template too.
	// Empty means app.kubernetes.io/name: Name.
	Selector map[string]string

	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	// +mapType=atomic
	NodeSelector map[string]string

	// ServiceAccountName is the name of the ServiceAccount to use to run this pod.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	// +optional
	ServiceAccountName string

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use
	// for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use.
	ImagePullSecrets []string

	// VolumeClaimRetentionPolicy — what happens to the pods' claims when the
	// StatefulSet is deleted or scaled down. Default: kept.
	VolumeClaimRetentionPolicy VolumeClaimRetentionPolicy

	// Probes — startup, liveness, and readiness probe configuration.
	Probes *types.ProbesConfig

	// SecurityContext — container-level security settings.
	SecurityContext *types.ContainerSecurityContext

	// PodSecurity — pod-level security settings.
	PodSecurity *types.PodSecurityContext

	// RollingUpdate — how a rollout replaces pods. nil: pods are replaced only
	// when deleted (OnDelete).
	RollingUpdate *types.RollingUpdateBehavior

	// Volumes — the pod's volumes.
	Volumes []types.VolumeSource
	// VolumeMounts — where the container mounts them, or defaults for every
	// container when Containers is set.
	VolumeMounts []types.VolumeMount

	// Containers — the pod's containers, when it has more than the one the
	// flat fields describe. See types.Container for how the two combine.
	Containers []types.Container

	// InitContainers — run in order, each to completion, before Containers
	// start. They inherit the flat Env, EnvFrom, SecurityContext and
	// VolumeMounts.
	InitContainers []types.Container
}

// VolumeClaimRetentionPolicy describes the policy used for PVCs
// created from the StatefulSet VolumeClaimTemplates.
type VolumeClaimRetentionPolicy struct {
	// WhenDeleted specifies what happens to PVCs created from StatefulSet
	// VolumeClaimTemplates when the StatefulSet is deleted. The default policy
	// of `Retain` causes PVCs to not be affected by StatefulSet deletion. The
	// `Delete` policy causes those PVCs to be deleted.
	WhenDeleted string
	// WhenScaled specifies what happens to PVCs created from StatefulSet
	// VolumeClaimTemplates when the StatefulSet is scaled down. The default
	// policy of `Retain` causes PVCs to not be affected by a scaledown. The
	// `Delete` policy causes the associated PVCs for any excess pods above
	// the replica count to be deleted.
	WhenScaled string
}
