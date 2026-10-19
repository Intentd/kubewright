package types

import corev1 "k8s.io/api/core/v1"

// Container is one container of a pod kind's Containers or InitContainers.
//
// Without Containers, a pod kind's flat fields (Image, Port, Env, ...)
// describe its one container. With Containers set, the flat Env, EnvFrom,
// SecurityContext and VolumeMounts are defaults every container inherits, and
// the flat Image, Command, Args, Port, Probes and Resources must be empty.
type Container struct {
	// Name — must be unique in the pod. Required.
	Name string `yaml:"name" json:"name"`
	// Image — the container image. Required.
	Image string `yaml:"image" json:"image"`
	// ImagePullPolicy — Always, IfNotPresent or Never. Empty: the Kubernetes default.
	ImagePullPolicy corev1.PullPolicy `yaml:"imagePullPolicy,omitempty" json:"imagePullPolicy,omitempty"`
	// Command — the entrypoint, replacing the image's.
	Command []string `yaml:"command,omitempty" json:"command,omitempty"`
	// Args — arguments to the entrypoint, replacing the image's.
	Args []string `yaml:"args,omitempty" json:"args,omitempty"`
	// Port — the port the container listens on; probes without a port use it.
	Port int32 `yaml:"port,omitempty" json:"port,omitempty"`
	// Protocol — for Port: TCP (the default), UDP or SCTP.
	Protocol corev1.Protocol `yaml:"protocol,omitempty" json:"protocol,omitempty"`
	// Env — environment variables, merged by name over the spec's flat Env.
	Env []EnvVar `yaml:"env,omitempty" json:"env,omitempty"`
	// EnvFrom — Secrets and ConfigMaps whose keys become environment
	// variables. Set, it replaces the spec's flat EnvFrom; unset, inherits it.
	EnvFrom *EnvFrom `yaml:"envFrom,omitempty" json:"envFrom,omitempty"`
	// Resources — requests and limits, or a preset to expand first.
	Resources *ResourceRequirements `yaml:"resources,omitempty" json:"resources,omitempty"`
	// Probes — the container's probes. Not allowed on an init container.
	Probes *ProbesConfig `yaml:"probes,omitempty" json:"probes,omitempty"`
	// SecurityContext — the container's security settings. Set, it replaces
	// the spec's flat SecurityContext; unset, inherits it.
	SecurityContext *ContainerSecurityContext `yaml:"securityContext,omitempty" json:"securityContext,omitempty"`
	// VolumeMounts — the pod's volumes to mount, merged by mount path over the
	// spec's flat VolumeMounts.
	VolumeMounts []VolumeMount `yaml:"volumeMounts,omitempty" json:"volumeMounts,omitempty"`
}
