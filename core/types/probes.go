package types

// ProbeConfig is one probe: an HTTP GET or a TCP connection, with timings
// from a preset, explicit fields, or both.
type ProbeConfig struct {
	// Type — http or tcp. Default: http when Path is set.
	Type string `yaml:"type,omitempty" json:"type,omitempty"`

	// Path — HTTP GET path. Required when type is "http".
	Path string `yaml:"path,omitempty" json:"path,omitempty"`

	// Port — the port to probe. Default: the container's port.
	Port int32 `yaml:"port,omitempty" json:"port,omitempty"`

	// Preset — a named set of timings: fast, standard, patient, slow-start,
	// or one from the registry. Explicit timings below override it.
	Preset string `yaml:"preset,omitempty" json:"preset,omitempty"`

	// InitialDelaySeconds — delay before the first probe.
	InitialDelaySeconds *int32 `yaml:"initialDelaySeconds,omitempty" json:"initialDelaySeconds,omitempty"`
	// PeriodSeconds — how often to probe.
	PeriodSeconds *int32 `yaml:"periodSeconds,omitempty" json:"periodSeconds,omitempty"`
	// FailureThreshold — consecutive failures before the probe counts as failed.
	FailureThreshold *int32 `yaml:"failureThreshold,omitempty" json:"failureThreshold,omitempty"`
	// SuccessThreshold — consecutive successes before it counts as passed
	// again. Must be 1 for liveness and startup.
	SuccessThreshold *int32 `yaml:"successThreshold,omitempty" json:"successThreshold,omitempty"`
	// TimeoutSeconds — how long one probe may take.
	TimeoutSeconds *int32 `yaml:"timeoutSeconds,omitempty" json:"timeoutSeconds,omitempty"`
}

// ProbesConfig is a container's startup, liveness and readiness probes.
type ProbesConfig struct {
	// Startup — holds off the other probes until the container has started;
	// failing it restarts the container.
	Startup *ProbeConfig `yaml:"startup,omitempty" json:"startup,omitempty"`
	// Liveness — failing it restarts the container.
	Liveness *ProbeConfig `yaml:"liveness,omitempty" json:"liveness,omitempty"`
	// Readiness — failing it takes the pod out of its Services' endpoints.
	Readiness *ProbeConfig `yaml:"readiness,omitempty" json:"readiness,omitempty"`
}
