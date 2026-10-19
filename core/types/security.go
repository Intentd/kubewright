package types

// PodSecurityContext is a pod's security settings: a preset, or the fields.
type PodSecurityContext struct {
	// Preset — baseline, restricted, hardened, or one from the registry.
	// presets.ExpandPodSecurity expands it.
	Preset string `yaml:"preset,omitempty" json:"preset,omitempty"`

	// RunAsNonRoot — true refuses to run any container as root.
	RunAsNonRoot *bool `yaml:"runAsNonRoot,omitempty" json:"runAsNonRoot,omitempty"`

	// RunAsUser — the user ID every container runs as.
	RunAsUser *int64 `yaml:"runAsUser,omitempty" json:"runAsUser,omitempty"`

	// RunAsGroup — the group ID every container runs as.
	RunAsGroup *int64 `yaml:"runAsGroup,omitempty" json:"runAsGroup,omitempty"`

	// FSGroup — the group that owns the pod's volumes.
	FSGroup *int64 `yaml:"fsGroup,omitempty" json:"fsGroup,omitempty"`
}

// HasMixedFields reports whether a field is set alongside the preset.
func (p *PodSecurityContext) HasMixedFields() bool {
	if p == nil || p.Preset == "" {
		return false
	}
	return p.RunAsNonRoot != nil || p.RunAsUser != nil || p.RunAsGroup != nil || p.FSGroup != nil
}

// ContainerSecurityContext is a container's security settings: a preset, or
// the fields.
type ContainerSecurityContext struct {
	// Preset — baseline, restricted, hardened, or one from the registry.
	// presets.ExpandContainerSecurity expands it.
	Preset string `yaml:"preset,omitempty" json:"preset,omitempty"`

	// AllowPrivilegeEscalation — false stops a process gaining more privileges
	// than its parent. Kubernetes default: true.
	AllowPrivilegeEscalation *bool `yaml:"allowPrivilegeEscalation,omitempty" json:"allowPrivilegeEscalation,omitempty"`

	// ReadOnlyRootFilesystem — true mounts the root filesystem read-only.
	ReadOnlyRootFilesystem *bool `yaml:"readOnlyRootFilesystem,omitempty" json:"readOnlyRootFilesystem,omitempty"`

	// RunAsNonRoot — true refuses to run the container as root.
	RunAsNonRoot *bool `yaml:"runAsNonRoot,omitempty" json:"runAsNonRoot,omitempty"`

	// RunAsUser — the user ID the container runs as.
	RunAsUser *int64 `yaml:"runAsUser,omitempty" json:"runAsUser,omitempty"`

	// RunAsGroup — the group ID the container runs as.
	RunAsGroup *int64 `yaml:"runAsGroup,omitempty" json:"runAsGroup,omitempty"`

	// Capabilities — Linux capabilities to add and drop.
	Capabilities *CapabilitiesConfig `yaml:"capabilities,omitempty" json:"capabilities,omitempty"`
}

// HasMixedFields reports whether any explicit field is set alongside a preset.
func (c *ContainerSecurityContext) HasMixedFields() bool {
	if c == nil || c.Preset == "" {
		return false
	}
	return c.AllowPrivilegeEscalation != nil ||
		c.ReadOnlyRootFilesystem != nil ||
		c.RunAsNonRoot != nil ||
		c.RunAsUser != nil ||
		c.RunAsGroup != nil ||
		c.Capabilities != nil
}

// CapabilitiesConfig is the Linux capabilities to add and drop.
type CapabilitiesConfig struct {
	// Add — capabilities to add to the default set.
	Add []string `yaml:"add,omitempty" json:"add,omitempty"`

	// Drop — capabilities to remove; ["ALL"] removes every one.
	Drop []string `yaml:"drop,omitempty" json:"drop,omitempty"`
}
