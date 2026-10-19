package types

// RollingUpdateBehavior is how a rollout replaces pods: a preset, or MaxSurge
// and MaxUnavailable.
type RollingUpdateBehavior struct {
	// Preset — safe, fast, blue-green, or one from the registry.
	// presets.ExpandRollingUpdate expands it.
	Preset string `yaml:"preset,omitempty" json:"preset,omitempty"`

	// MaxSurge — pods allowed above the desired count, as a count ("1") or a
	// percentage ("25%").
	MaxSurge string `yaml:"maxSurge,omitempty" json:"maxSurge,omitempty"`

	// MaxUnavailable — pods allowed to be unavailable, as a count or a
	// percentage.
	MaxUnavailable string `yaml:"maxUnavailable,omitempty" json:"maxUnavailable,omitempty"`
}
