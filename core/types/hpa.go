package types

// HPABehavior is how fast a HorizontalPodAutoscaler scales up and down: a
// preset, or ScaleUp and ScaleDown.
type HPABehavior struct {
	// Preset — web, api, latency-sensitive, batch, cost-optimized, or one
	// from the registry. presets.ExpandHPA expands it.
	Preset string `yaml:"preset,omitempty" json:"preset,omitempty"`

	// ScaleUp — how fast replicas are added.
	ScaleUp *HPAScalingRules `yaml:"scaleUp,omitempty" json:"scaleUp,omitempty"`

	// ScaleDown — how fast replicas are removed.
	ScaleDown *HPAScalingRules `yaml:"scaleDown,omitempty" json:"scaleDown,omitempty"`
}

// HPAScalingRules is one scaling direction, up or down.
type HPAScalingRules struct {
	// StabilizationWindowSeconds — how far back to look before acting.
	// Kubernetes defaults: 300 down, 0 up.
	StabilizationWindowSeconds int32 `yaml:"stabilizationWindowSeconds,omitempty" json:"stabilizationWindowSeconds,omitempty"`

	// Policies — the limits on each step; SelectPolicy picks among them.
	Policies []HPAScalingPolicy `yaml:"policies,omitempty" json:"policies,omitempty"`

	// SelectPolicy — Max, Min or Disabled. Default: Max.
	SelectPolicy string `yaml:"selectPolicy,omitempty" json:"selectPolicy,omitempty"`
}

// HPAScalingPolicy limits one scaling step per period.
type HPAScalingPolicy struct {
	// Type — Pods (a count) or Percent (of current replicas).
	Type string `yaml:"type" json:"type"`

	// Value — the pods or percent allowed per period.
	Value int32 `yaml:"value" json:"value"`

	// PeriodSeconds — the period, 1 to 1800.
	PeriodSeconds int32 `yaml:"periodSeconds" json:"periodSeconds"`
}

type ScaleTargetRef struct {
	// APIVersion — the target's group and version, e.g. "apps/v1".
	APIVersion string `yaml:"apiVersion" json:"apiVersion"`
	// Kind — the target's kind, e.g. Deployment or StatefulSet.
	Kind string `yaml:"kind" json:"kind"`
	// Name — the target's name, in the HorizontalPodAutoscaler's namespace.
	Name string `yaml:"name" json:"name"`
}
