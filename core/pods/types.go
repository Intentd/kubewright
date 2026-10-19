package pods

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/intentd/kubewright/core/types"
)

// Spec describes a Pod.
type Spec struct {
	// Name — the Pod's name. Required.
	Name string

	// Image — container image. Required.
	Image string

	// Namespace — target namespace. Required.
	Namespace string

	// Port — container port. 0 means no port exposed.
	Port int
	// Protocol — for Port: TCP (the default), UDP or SCTP.
	Protocol corev1.Protocol

	// Labels — merged labels from both sources.
	Labels map[string]string

	// Annotations — merged annotations from both sources.
	Annotations map[string]string

	// Resources — CPU and memory requests/limits. nil means no limits set.
	Resources *types.ResourceRequirements

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

	// RestartPolicy — Always, OnFailure or Never. Empty leaves the Kubernetes
	// default, Always.
	RestartPolicy corev1.RestartPolicy

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
