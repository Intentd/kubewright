package types

import corev1 "k8s.io/api/core/v1"

// Container is one container of a pod, for a pod kind's Containers and
// InitContainers.
//
// A pod kind's flat fields (Image, Port, Env, ...) describe its one container
// when Containers is empty. With Containers set, the flat Env, EnvFrom,
// SecurityContext and VolumeMounts become defaults every container and init
// container inherits (a container's own value wins: Env by name, VolumeMounts
// by mount path), and the flat Image, Command, Args, Port, Probes and
// Resources must be empty: they belong to one container, so say which.
type Container struct {
	// Name — must be unique in the pod. Required.
	Name string `yaml:"name" json:"name"`
	// Image — container image. Required.
	Image string `yaml:"image" json:"image"`
	// ImagePullPolicy — Always, IfNotPresent or Never; empty leaves the
	// Kubernetes default.
	ImagePullPolicy corev1.PullPolicy `yaml:"imagePullPolicy,omitempty" json:"imagePullPolicy,omitempty"`
	Command         []string          `yaml:"command,omitempty" json:"command,omitempty"`
	Args            []string          `yaml:"args,omitempty" json:"args,omitempty"`
	// Port — the port the container listens on; probes without a port use it.
	Port int32 `yaml:"port,omitempty" json:"port,omitempty"`
	// Protocol — for Port; defaults to TCP.
	Protocol  corev1.Protocol       `yaml:"protocol,omitempty" json:"protocol,omitempty"`
	Env       []EnvVar              `yaml:"env,omitempty" json:"env,omitempty"`
	EnvFrom   *EnvFrom              `yaml:"envFrom,omitempty" json:"envFrom,omitempty"`
	Resources *ResourceRequirements `yaml:"resources,omitempty" json:"resources,omitempty"`
	// Probes — not allowed on an init container.
	Probes          *ProbesConfig             `yaml:"probes,omitempty" json:"probes,omitempty"`
	SecurityContext *ContainerSecurityContext `yaml:"securityContext,omitempty" json:"securityContext,omitempty"`
	VolumeMounts    []VolumeMount             `yaml:"volumeMounts,omitempty" json:"volumeMounts,omitempty"`
}
