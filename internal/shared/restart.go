package shared

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// ValidateRestartPolicy reports an error when policy is set and is not one of
// allowed. An empty policy is valid: each kind applies its own default.
func ValidateRestartPolicy(policy corev1.RestartPolicy, allowed ...corev1.RestartPolicy) error {
	if policy == "" {
		return nil
	}
	for _, a := range allowed {
		if policy == a {
			return nil
		}
	}
	return fmt.Errorf("restartPolicy %q: must be one of %v", policy, allowed)
}

// RestartPolicyOr returns policy, or def when policy is empty.
func RestartPolicyOr(policy, def corev1.RestartPolicy) corev1.RestartPolicy {
	if policy == "" {
		return def
	}
	return policy
}
