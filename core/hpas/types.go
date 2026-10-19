package hpas

import "github.com/intentd/kubewright/core/types"

// Spec describes a HorizontalPodAutoscaler.
type Spec struct {
	// Name — HPA resource name. Required.
	Name string

	// Namespace — target namespace.
	Namespace string

	// ScaleTargetRef — the target workload this HPA scales.
	// Supports Deployment, ReplicaSet, StatefulSet, or any scalable resource.
	ScaleTargetRef types.ScaleTargetRef

	// MinReplicas — minimum pod replica count. Default: 1.
	MinReplicas int32

	// MaxReplicas — maximum pod replica count. Required.
	MaxReplicas int32

	// TargetCPUUtilizationPercentage — CPU utilization target (0-100).
	// When 0, no CPU metric is configured.
	TargetCPUUtilizationPercentage int32

	// Labels applied to HPA metadata.
	Labels map[string]string

	// Annotations — set on the object.
	Annotations map[string]string

	// Behavior — how fast to scale; expand a preset with presets.ExpandHPA
	// first. nil uses the Kubernetes defaults.
	Behavior *types.HPABehavior
}
