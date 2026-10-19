package ingresses

// Spec describes an Ingress.
type Spec struct {
	// Name — Ingress resource name. Required.
	Name string

	// Namespace — target namespace.
	Namespace string

	// Host — virtual host for the Ingress rule.
	Host string

	// ServiceName — backend Service name.
	ServiceName string

	// ServicePort — backend Service port.
	ServicePort int32

	// Path — HTTP path prefix. Default: "/"
	Path string

	// PathType — Prefix | Exact | ImplementationSpecific. Default: Prefix.
	PathType string

	// IngressClass — Ingress class name. Empty means cluster default.
	IngressClass string

	// Labels applied to Ingress metadata.
	Labels map[string]string

	// Annotations applied to Ingress metadata.
	Annotations map[string]string

	// TLS — optional TLS configuration. nil means no TLS.
	TLS *TLS
}

// TLS is an Ingress's TLS configuration.
type TLS struct {
	// Create — whether to configure TLS on this Ingress.
	Create bool

	// SecretName — name of the kubernetes.io/tls Secret.
	SecretName string

	// Hosts — hostnames covered by the TLS certificate.
	Hosts []string

	// ValidFor — certificate validity duration string passed to GenerateTLSBundle.
	ValidFor string
}
