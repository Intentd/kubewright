package shared

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	"github.com/intentd/kubewright/core/types"
)

// BuildProbe constructs a Kubernetes Probe from a ProbeConfig whose preset,
// if any, presets.ExpandProbes already resolved. It sets only the timings cfg
// sets; the rest are left to Kubernetes' defaults. containerPort is used as
// the probe target when cfg.Port is 0. Returns nil if cfg is nil or no probe
// type can be inferred.
func BuildProbe(cfg *types.ProbeConfig, containerPort int32) *corev1.Probe {
	if cfg == nil {
		return nil
	}

	isHTTP := cfg.Type == "http" || (cfg.Type == "" && cfg.Path != "")
	isTCP := cfg.Type == "tcp"

	if !isHTTP && !isTCP {
		return nil
	}

	port := containerPort
	if cfg.Port > 0 {
		port = cfg.Port
	}

	probe := &corev1.Probe{
		InitialDelaySeconds: ptr.Deref(cfg.InitialDelaySeconds, 0),
		PeriodSeconds:       ptr.Deref(cfg.PeriodSeconds, 0),
		FailureThreshold:    ptr.Deref(cfg.FailureThreshold, 0),
		SuccessThreshold:    ptr.Deref(cfg.SuccessThreshold, 0),
		TimeoutSeconds:      ptr.Deref(cfg.TimeoutSeconds, 0),
	}

	if isHTTP {
		probe.ProbeHandler = corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   cfg.Path,
				Port:   intstr.FromInt32(port),
				Scheme: corev1.URISchemeHTTP,
			},
		}
	} else {
		probe.ProbeHandler = corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt32(port),
			},
		}
	}

	return probe
}

// ApplyProbes sets startup, liveness, and readiness probes on c.
// containerPort is used for probes that do not declare their own port.
// No-op when probes is nil.
func ApplyProbes(c *corev1.Container, probes *types.ProbesConfig, containerPort int32) {
	if probes == nil {
		return
	}
	c.StartupProbe = BuildProbe(probes.Startup, containerPort)
	c.LivenessProbe = BuildProbe(probes.Liveness, containerPort)
	c.ReadinessProbe = BuildProbe(probes.Readiness, containerPort)
}
