package types

// CustomResourceMetadata is the metadata a custom resource spec sets.
type CustomResourceMetadata struct {
	// Name — the object's name. Required.
	Name string `json:"name" yaml:"name"`

	// Namespace — where to create it. Empty means the owner's.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// Namespaced — false for a cluster-scoped kind: the object gets no
	// namespace. nil or true: namespaced.
	Namespaced *bool `json:"namespaced,omitempty" yaml:"namespaced,omitempty"`

	// Labels — set on the object.
	Labels Labels `json:"labels,omitempty" yaml:"labels,omitempty"`

	// Annotations — set on the object.
	Annotations Labels `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}
