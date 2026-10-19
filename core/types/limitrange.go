package types

// LimitRangeItem is the defaults and bounds for one type of object in a namespace.
type LimitRangeItem struct {
	// Type — Container, Pod or PersistentVolumeClaim.
	Type string `yaml:"type" json:"type"`

	// Max — the most each resource may be, e.g. {"cpu": "2", "memory": "4Gi"}.
	Max map[string]string `yaml:"max,omitempty" json:"max,omitempty"`

	// Min — the least each resource may be.
	Min map[string]string `yaml:"min,omitempty" json:"min,omitempty"`

	// Default — the limits given to a container that sets none.
	Default map[string]string `yaml:"default,omitempty" json:"default,omitempty"`

	// DefaultRequest — the requests given to a container that sets none.
	DefaultRequest map[string]string `yaml:"defaultRequest,omitempty" json:"defaultRequest,omitempty"`

	// MaxLimitRequestRatio — the most a limit may be, as a multiple of its request.
	MaxLimitRequestRatio map[string]string `yaml:"maxLimitRequestRatio,omitempty" json:"maxLimitRequestRatio,omitempty"`
}
