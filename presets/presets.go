package presets

import (
	"strconv"

	"github.com/intentd/kubewright/core/types"
)

// Resource is a named set of CPU and memory requests and limits.
type Resource struct {
	// Name — what a resources.preset names to use this preset; names ignore case.
	Name string `yaml:"name" json:"name"`
	// Description — what the preset is for.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// Requests — the resources guaranteed, e.g. {"cpu": "100m", "memory": "128Mi"}.
	Requests map[string]string `yaml:"requests,omitempty" json:"requests,omitempty"`
	// Limits — the most the container may use.
	Limits map[string]string `yaml:"limits,omitempty" json:"limits,omitempty"`
}

// Requirements returns the preset as a container's resource requirements.
func (p Resource) Requirements() *types.ResourceRequirements {
	return &types.ResourceRequirements{Requests: p.Requests, Limits: p.Limits}
}

// Probe is a named set of probe timings.
type Probe struct {
	// Name — what a probe's preset names to use this preset; names ignore case.
	Name string `yaml:"name" json:"name"`
	// Description — what the preset is for.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// InitialDelaySeconds — wait this long after start before the first probe.
	InitialDelaySeconds int32 `yaml:"initialDelaySeconds,omitempty" json:"initialDelaySeconds,omitempty"`
	// PeriodSeconds — how often to probe.
	PeriodSeconds int32 `yaml:"periodSeconds,omitempty" json:"periodSeconds,omitempty"`
	// FailureThreshold — consecutive failures before the probe counts as failed.
	FailureThreshold int32 `yaml:"failureThreshold,omitempty" json:"failureThreshold,omitempty"`
	// SuccessThreshold — consecutive successes before it counts as passed again.
	SuccessThreshold int32 `yaml:"successThreshold,omitempty" json:"successThreshold,omitempty"`
	// TimeoutSeconds — how long one probe may take.
	TimeoutSeconds int32 `yaml:"timeoutSeconds,omitempty" json:"timeoutSeconds,omitempty"`
}

// ContainerSecurity is a named container security context.
type ContainerSecurity struct {
	// Name — what a securityContext.preset names to use this preset; names ignore case.
	Name string `yaml:"name" json:"name"`
	// Description — what the preset is for.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// AllowPrivilegeEscalation — false stops a process gaining more privileges
	// than its parent.
	AllowPrivilegeEscalation *bool `yaml:"allowPrivilegeEscalation,omitempty" json:"allowPrivilegeEscalation,omitempty"`
	// ReadOnlyRootFilesystem — true mounts the container's root filesystem read-only.
	ReadOnlyRootFilesystem *bool `yaml:"readOnlyRootFilesystem,omitempty" json:"readOnlyRootFilesystem,omitempty"`
	// RunAsNonRoot — true refuses to start the container as root.
	RunAsNonRoot *bool `yaml:"runAsNonRoot,omitempty" json:"runAsNonRoot,omitempty"`
	// RunAsUser — the user ID the container runs as.
	RunAsUser *int64 `yaml:"runAsUser,omitempty" json:"runAsUser,omitempty"`
	// RunAsGroup — the group ID the container runs as.
	RunAsGroup *int64 `yaml:"runAsGroup,omitempty" json:"runAsGroup,omitempty"`
	// Capabilities — Linux capabilities to add and drop.
	Capabilities *types.CapabilitiesConfig `yaml:"capabilities,omitempty" json:"capabilities,omitempty"`
}

// Context returns the preset as a container security context.
func (p ContainerSecurity) Context() *types.ContainerSecurityContext {
	return &types.ContainerSecurityContext{
		AllowPrivilegeEscalation: p.AllowPrivilegeEscalation,
		ReadOnlyRootFilesystem:   p.ReadOnlyRootFilesystem,
		RunAsNonRoot:             p.RunAsNonRoot,
		RunAsUser:                p.RunAsUser,
		RunAsGroup:               p.RunAsGroup,
		Capabilities:             p.Capabilities,
	}
}

// PodSecurity is a named pod security context.
type PodSecurity struct {
	// Name — what a podSecurity.preset names to use this preset; names ignore case.
	Name string `yaml:"name" json:"name"`
	// Description — what the preset is for.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// RunAsNonRoot — true refuses to start any container as root.
	RunAsNonRoot *bool `yaml:"runAsNonRoot,omitempty" json:"runAsNonRoot,omitempty"`
	// RunAsUser — the user ID every container runs as.
	RunAsUser *int64 `yaml:"runAsUser,omitempty" json:"runAsUser,omitempty"`
	// RunAsGroup — the group ID every container runs as.
	RunAsGroup *int64 `yaml:"runAsGroup,omitempty" json:"runAsGroup,omitempty"`
	// FSGroup — the group that owns mounted volumes.
	FSGroup *int64 `yaml:"fsGroup,omitempty" json:"fsGroup,omitempty"`
}

// Context returns the preset as a pod security context.
func (p PodSecurity) Context() *types.PodSecurityContext {
	return &types.PodSecurityContext{
		RunAsNonRoot: p.RunAsNonRoot,
		RunAsUser:    p.RunAsUser,
		RunAsGroup:   p.RunAsGroup,
		FSGroup:      p.FSGroup,
	}
}

// RollingUpdate is a named Deployment rolling update strategy. Each value is
// a count ("1") or a percentage ("25%").
type RollingUpdate struct {
	// Name — what a rollingUpdate.preset names to use this preset; names ignore case.
	Name string `yaml:"name" json:"name"`
	// Description — what the preset is for.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// MaxSurge — pods allowed above the desired count during a rollout.
	MaxSurge string `yaml:"maxSurge,omitempty" json:"maxSurge,omitempty"`
	// MaxUnavailable — pods allowed to be unavailable during a rollout.
	MaxUnavailable string `yaml:"maxUnavailable,omitempty" json:"maxUnavailable,omitempty"`
}

// Behavior returns the preset as a rolling update behavior.
func (p RollingUpdate) Behavior() *types.RollingUpdateBehavior {
	return &types.RollingUpdateBehavior{MaxSurge: p.MaxSurge, MaxUnavailable: p.MaxUnavailable}
}

// HPA is a named HorizontalPodAutoscaler setup: replica bounds, a CPU target
// and scaling behavior. The numbers are strings, so a program may resolve them
// before use.
type HPA struct {
	// Name — what a HorizontalPodAutoscaler's behavior.preset names to use this preset; names ignore case.
	Name string `yaml:"name" json:"name"`
	// Description — what the preset is for.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// MinReplicas — the fewest pods, used when the HPA sets none.
	MinReplicas string `yaml:"minReplicas,omitempty" json:"minReplicas,omitempty"`
	// MaxReplicas — the most pods, used when the HPA sets none.
	MaxReplicas string `yaml:"maxReplicas,omitempty" json:"maxReplicas,omitempty"`
	// TargetCPUUtilizationPercentage — the average CPU use to scale toward,
	// used when the HPA sets none.
	TargetCPUUtilizationPercentage string `yaml:"targetCPUUtilizationPercentage,omitempty" json:"targetCPUUtilizationPercentage,omitempty"`
	// Behavior — how fast to scale up and down. Replaces the HPA's behavior.
	Behavior *types.HPABehavior `yaml:"behavior,omitempty" json:"behavior,omitempty"`
}

// CPUTarget returns the CPU utilization target as a number, or 0 when it is
// unset or not a whole number.
func (p HPA) CPUTarget() int32 { return wholeNumber(p.TargetCPUUtilizationPercentage) }

// Replicas returns the preset's replica bounds as numbers; each is 0 when it
// is unset or not a whole number.
func (p HPA) Replicas() (minReplicas, maxReplicas int32) {
	return wholeNumber(p.MinReplicas), wholeNumber(p.MaxReplicas)
}

func wholeNumber(s string) int32 {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int32(n)
}

// PDB is a named PodDisruptionBudget limit: MinAvailable or MaxUnavailable,
// each a count or a percentage.
type PDB struct {
	// Name — what a PodDisruptionBudget names to use this preset; names ignore case.
	Name string `yaml:"name" json:"name"`
	// Description — what the preset is for.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// MinAvailable — pods that must stay up, used when the budget doesn't set it.
	MinAvailable string `yaml:"minAvailable,omitempty" json:"minAvailable,omitempty"`
	// MaxUnavailable — pods that may be down, used when the budget doesn't set it.
	MaxUnavailable string `yaml:"maxUnavailable,omitempty" json:"maxUnavailable,omitempty"`
}

// NetworkPolicy is a named set of NetworkPolicy rules.
type NetworkPolicy struct {
	// Name — what a NetworkPolicy names to use this preset; names ignore case.
	Name string `yaml:"name" json:"name"`
	// Description — what the preset is for.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// PodSelector — the labels of the pods the policy applies to, used when the
	// policy sets none. An empty map selects every pod.
	PodSelector map[string]string `yaml:"podSelector,omitempty" json:"podSelector,omitempty"`
	// Ingress — the rules for traffic in. Replaces the policy's.
	Ingress []types.NetworkPolicyIngressRule `yaml:"ingress,omitempty" json:"ingress,omitempty"`
	// Egress — the rules for traffic out. Replaces the policy's.
	Egress []types.NetworkPolicyEgressRule `yaml:"egress,omitempty" json:"egress,omitempty"`
	// PolicyTypes — Ingress, Egress or both. Replaces the policy's.
	PolicyTypes []string `yaml:"policyTypes,omitempty" json:"policyTypes,omitempty"`
}

// ResourceQuota is a named set of a namespace's hard limits.
type ResourceQuota struct {
	// Name — what a ResourceQuota names to use this preset; names ignore case.
	Name string `yaml:"name" json:"name"`
	// Description — what the preset is for.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// Hard — the namespace's limits, e.g. {"pods": "10", "requests.cpu": "4"}.
	// Replaces the quota's.
	Hard map[string]string `yaml:"hard" json:"hard"`
}

// LimitRange is a named set of LimitRange items. There are no built-in ones.
type LimitRange struct {
	// Name — what a LimitRange names to use this preset; names ignore case.
	Name string `yaml:"name" json:"name"`
	// Description — what the preset is for.
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// Limits — per-container and per-pod defaults and bounds, used when the
	// LimitRange sets none.
	Limits []types.LimitRangeItem `yaml:"limits" json:"limits"`
}

func (p Resource) presetName() string          { return p.Name }
func (p Probe) presetName() string             { return p.Name }
func (p ContainerSecurity) presetName() string { return p.Name }
func (p PodSecurity) presetName() string       { return p.Name }
func (p RollingUpdate) presetName() string     { return p.Name }
func (p HPA) presetName() string               { return p.Name }
func (p PDB) presetName() string               { return p.Name }
func (p NetworkPolicy) presetName() string     { return p.Name }
func (p ResourceQuota) presetName() string     { return p.Name }
func (p LimitRange) presetName() string        { return p.Name }
