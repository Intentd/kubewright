package replicasets

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/intentd/kubewright/core/types"
)

// Spec describes a ReplicaSet.
type Spec struct {
	// Name — the ReplicaSet's name. Required.
	Name string

	// Image — container image. Required.
	Image string

	// Replicas — number of pod replicas. Default: 1.
	Replicas int32

	// Port — container port. 0 means no port exposed.
	Port int32
	// Protocol — for Port: TCP (the default), UDP or SCTP.
	Protocol corev1.Protocol

	// Namespace — target namespace. Required.
	Namespace string

	// Labels — applied to the ReplicaSet and pod template.
	Labels map[string]string

	// Annotations — applied to the ReplicaSet.
	Annotations map[string]string

	// Env — environment variables.
	Env []types.EnvVar
	// EnvFrom — Secrets and ConfigMaps whose keys become environment variables.
	EnvFrom *types.EnvFrom

	// Resources — CPU and memory requests/limits. nil means no limits set.
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

	// Probes — startup, liveness, and readiness probe configuration.
	Probes *types.ProbesConfig

	// SecurityContext — container-level security settings.
	SecurityContext *types.ContainerSecurityContext

	// PodSecurity — pod-level security settings.
	PodSecurity *types.PodSecurityContext

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
