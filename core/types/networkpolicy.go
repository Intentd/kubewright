package types

// NetworkPolicyIngressRule describes one ingress rule.
type NetworkPolicyIngressRule struct {
	// From — the sources allowed in. Empty allows every source.
	From []NetworkPolicyPeer `yaml:"from,omitempty" json:"from,omitempty"`
	// Ports — the ports allowed. Empty allows every port.
	Ports []NetworkPolicyPort `yaml:"ports,omitempty" json:"ports,omitempty"`
}

// NetworkPolicyEgressRule describes one egress rule.
type NetworkPolicyEgressRule struct {
	// To — the destinations allowed out. Empty allows every destination.
	To []NetworkPolicyPeer `yaml:"to,omitempty" json:"to,omitempty"`
	// Ports — the ports allowed. Empty allows every port.
	Ports []NetworkPolicyPort `yaml:"ports,omitempty" json:"ports,omitempty"`
}

// NetworkPolicyPeer selects a set of pods or namespaces.
type NetworkPolicyPeer struct {
	// PodSelector — pods in the policy's namespace, by label. Empty selects
	// every pod.
	PodSelector map[string]string `yaml:"podSelector,omitempty" json:"podSelector,omitempty"`
	// NamespaceSelector — namespaces, by label. Empty selects every namespace.
	NamespaceSelector map[string]string `yaml:"namespaceSelector,omitempty" json:"namespaceSelector,omitempty"`
	// IPBlock — a range of IP addresses.
	IPBlock *NetworkPolicyIPBlock `yaml:"ipBlock,omitempty" json:"ipBlock,omitempty"`
}

// NetworkPolicyIPBlock describes an IP CIDR range.
type NetworkPolicyIPBlock struct {
	// CIDR — the range allowed, e.g. "10.0.0.0/16".
	CIDR string `yaml:"cidr" json:"cidr"`
	// Except — ranges inside CIDR to leave out.
	Except []string `yaml:"except,omitempty" json:"except,omitempty"`
}

// NetworkPolicyPort describes a port allowed by a rule.
type NetworkPolicyPort struct {
	// Protocol — TCP (default), UDP, or SCTP.
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
	// Port — a port number or name.
	Port string `yaml:"port,omitempty" json:"port,omitempty"`
}
