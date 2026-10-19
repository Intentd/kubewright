package types

// ResourceRequirements is a container's requests and limits as quantity
// strings, or a preset that supplies both.
type ResourceRequirements struct {
	// Preset — a named size (tiny, small, medium, large, …). When set,
	// Requests and Limits are ignored: the preset supplies both.
	Preset string `yaml:"preset,omitempty" json:"preset,omitempty" validate:"omitempty"`
	// Requests — the resources guaranteed, e.g. {"cpu": "100m", "memory": "128Mi"}.
	Requests map[string]string `yaml:"requests,omitempty" json:"requests,omitempty" validate:"omitempty"`
	// Limits — the most the container may use, e.g. {"cpu": "500m", "memory": "512Mi"}.
	Limits map[string]string `yaml:"limits,omitempty" json:"limits,omitempty" validate:"omitempty"`
}

// EnvVar is one environment variable.
type EnvVar struct {
	// Name — the variable's name.
	Name string `yaml:"name" json:"name"`
	// Value — a literal value. Set this or ValueFrom.
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
	// ValueFrom — read the value from a Secret or ConfigMap key.
	ValueFrom *ValueFrom `yaml:"valueFrom,omitempty" json:"valueFrom,omitempty"`
}

// ValueFrom is where an environment variable's value comes from. Set one.
type ValueFrom struct {
	// SecretKeyRef — a key of a Secret in the pod's namespace.
	SecretKeyRef *SecretKeyRef `yaml:"secretKeyRef,omitempty" json:"secretKeyRef,omitempty"`
	// ConfigMapKeyRef — a key of a ConfigMap in the pod's namespace.
	ConfigMapKeyRef *ConfigMapKeyRef `yaml:"configMapKeyRef,omitempty" json:"configMapKeyRef,omitempty"`
}

// SecretKeyRef selects a key of a Secret.
type SecretKeyRef struct {
	// Name — the Secret's name.
	Name string `yaml:"name" json:"name"`
	// Key — the key within the Secret.
	Key string `yaml:"key" json:"key"`
	// Optional — true: the pod starts even if the Secret or key is missing.
	Optional *bool `yaml:"optional,omitempty" json:"optional,omitempty"`
}

// ConfigMapKeyRef selects a key of a ConfigMap.
type ConfigMapKeyRef struct {
	// Name — the ConfigMap's name.
	Name string `yaml:"name" json:"name"`
	// Key — the key within the ConfigMap.
	Key string `yaml:"key" json:"key"`
	// Optional — true: the pod starts even if the ConfigMap or key is missing.
	Optional *bool `yaml:"optional,omitempty" json:"optional,omitempty"`
}

// EnvFromRef imports the keys of one Secret or ConfigMap as environment
// variables. With Keys set, only those keys are imported, each as its own
// variable, which is how Suffix can rename them.
type EnvFromRef struct {
	// Name — the Secret or ConfigMap to import.
	Name string `yaml:"name" json:"name" validate:"required"`
	// Prefix — prepended to each imported variable's name.
	Prefix string `yaml:"prefix,omitempty" json:"prefix,omitempty"`
	// Suffix — appended to each imported variable's name. Needs Keys.
	Suffix string `yaml:"suffix,omitempty" json:"suffix,omitempty"`
	// Optional — true: the pod starts even if the source is missing.
	Optional *bool `yaml:"optional,omitempty" json:"optional,omitempty"`
	// Keys — import only these keys, each as its own variable named
	// Prefix + key + Suffix. Empty imports every key.
	Keys []string `yaml:"keys,omitempty" json:"keys,omitempty"`
}

// EnvFrom is the Secrets and ConfigMaps whose keys become environment variables.
type EnvFrom struct {
	// SecretRef — Secrets whose keys become environment variables.
	SecretRef []EnvFromRef `yaml:"secretRef,omitempty" json:"secretRef,omitempty"`
	// ConfigMapRef — ConfigMaps whose keys become environment variables.
	ConfigMapRef []EnvFromRef `yaml:"configMapRef,omitempty" json:"configMapRef,omitempty"`
}
