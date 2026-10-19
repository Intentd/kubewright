package shared

import (
	"fmt"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/intentd/kubewright/core/types"
)

// PodContainers returns a pod's containers and init containers. main is the
// container a pod kind's flat fields describe.
//
// Without containers, main is the pod's one container. With containers, main
// only lends its shared fields (Env, EnvFrom, SecurityContext, VolumeMounts)
// as defaults. Init containers inherit those defaults either way.
func PodContainers(main types.Container, containers, inits []types.Container) ([]corev1.Container, []corev1.Container) {
	var out []corev1.Container
	if len(containers) == 0 {
		out = []corev1.Container{BuildContainer(main)}
	} else {
		out = make([]corev1.Container, 0, len(containers))
		for _, c := range containers {
			out = append(out, BuildContainer(withDefaults(c, main)))
		}
	}
	var initOut []corev1.Container
	for _, c := range inits {
		initOut = append(initOut, BuildContainer(withDefaults(c, main)))
	}
	return out, initOut
}

// ValidateContainers checks main, containers and init containers together:
// with containers, main's container-only fields must be empty; every listed
// container needs a name and an image, names are unique in the pod, init
// containers have no probes, and no preset is left unexpanded (see
// presets.ExpandContainer).
func ValidateContainers(main types.Container, containers, inits []types.Container) error {
	var problems []string
	names := map[string]bool{}
	if len(containers) == 0 {
		if main.Image == "" {
			problems = append(problems, "image is required")
		}
		names[main.Name] = true
	} else {
		var set []string
		for field, isSet := range map[string]bool{
			"image":     main.Image != "",
			"command":   len(main.Command) > 0,
			"args":      len(main.Args) > 0,
			"port":      main.Port != 0,
			"probes":    main.Probes != nil,
			"resources": main.Resources != nil,
		} {
			if isSet {
				set = append(set, field)
			}
		}
		if len(set) > 0 {
			sort.Strings(set)
			problems = append(problems, fmt.Sprintf("%s set next to containers: they belong to one container, so set them in containers",
				strings.Join(set, ", ")))
		}
	}
	check := func(kind string, i int, c types.Container) {
		if c.Name == "" {
			problems = append(problems, fmt.Sprintf("%s[%d]: name is required", kind, i))
		} else if names[c.Name] {
			problems = append(problems, fmt.Sprintf("%s[%d]: name %q is already used in the pod", kind, i, c.Name))
		}
		names[c.Name] = true
		if c.Image == "" {
			problems = append(problems, fmt.Sprintf("%s[%d] %q: image is required", kind, i, c.Name))
		}
	}
	unexpanded := func(where string, c types.Container) {
		for _, p := range UnexpandedPresets(c) {
			problems = append(problems, fmt.Sprintf("%s: %s is not expanded: call presets.ExpandContainer first", where, p))
		}
	}
	unexpanded("flat fields", main)
	for i, c := range containers {
		check("containers", i, c)
		unexpanded(fmt.Sprintf("containers[%d] %q", i, c.Name), c)
	}
	for i, c := range inits {
		check("initContainers", i, c)
		if c.Probes != nil {
			problems = append(problems, fmt.Sprintf("initContainers[%d] %q: an init container has no probes", i, c.Name))
		}
		unexpanded(fmt.Sprintf("initContainers[%d] %q", i, c.Name), c)
	}
	if len(problems) > 0 {
		return fmt.Errorf("%s", strings.Join(problems, "; "))
	}
	return nil
}

// UnexpandedPresets names the presets c still refers to by name: Build
// can't resolve them, so they must be expanded first.
func UnexpandedPresets(c types.Container) []string {
	var out []string
	if c.Resources != nil && c.Resources.Preset != "" {
		out = append(out, fmt.Sprintf("resources preset %q", c.Resources.Preset))
	}
	if c.SecurityContext != nil && c.SecurityContext.Preset != "" {
		out = append(out, fmt.Sprintf("securityContext preset %q", c.SecurityContext.Preset))
	}
	if p := c.Probes; p != nil {
		for name, cfg := range map[string]*types.ProbeConfig{"startup": p.Startup, "liveness": p.Liveness, "readiness": p.Readiness} {
			if cfg != nil && cfg.Preset != "" {
				out = append(out, fmt.Sprintf("%s probe preset %q", name, cfg.Preset))
			}
		}
	}
	sort.Strings(out)
	return out
}

// BuildContainer returns the Kubernetes container c describes.
func BuildContainer(c types.Container) corev1.Container {
	out := corev1.Container{
		Name:            c.Name,
		Image:           c.Image,
		ImagePullPolicy: c.ImagePullPolicy,
	}
	if len(c.Command) > 0 {
		out.Command = c.Command
	}
	if len(c.Args) > 0 {
		out.Args = c.Args
	}
	if c.Port > 0 {
		out.Ports = []corev1.ContainerPort{{ContainerPort: c.Port, Protocol: c.Protocol}}
	}
	if c.Resources != nil {
		out.Resources = BuildResourceRequirements(c.Resources)
	}
	ApplyProbes(&out, c.Probes, c.Port)
	if c.SecurityContext != nil {
		out.SecurityContext = BuildContainerSecurityContext(c.SecurityContext)
	}
	out.Env = BuildEnv(c.Env)
	envFrom, extraEnv := ExpandEnvFrom(c.EnvFrom)
	out.EnvFrom = envFrom
	out.Env = append(out.Env, extraEnv...)
	if mounts := BuildVolumeMounts(c.VolumeMounts); len(mounts) > 0 {
		out.VolumeMounts = mounts
	}
	return out
}

// BuildEnv converts env vars to Kubernetes ones; nil for none.
func BuildEnv(env []types.EnvVar) []corev1.EnvVar {
	if len(env) == 0 {
		return nil
	}
	out := make([]corev1.EnvVar, 0, len(env))
	for _, ev := range env {
		kev := corev1.EnvVar{Name: ev.Name}
		if ev.ValueFrom == nil {
			kev.Value = ev.Value
			out = append(out, kev)
			continue
		}
		kev.ValueFrom = &corev1.EnvVarSource{}
		if r := ev.ValueFrom.SecretKeyRef; r != nil {
			kev.ValueFrom.SecretKeyRef = &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: r.Name},
				Key:                  r.Key,
				Optional:             r.Optional,
			}
		}
		if r := ev.ValueFrom.ConfigMapKeyRef; r != nil {
			kev.ValueFrom.ConfigMapKeyRef = &corev1.ConfigMapKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: r.Name},
				Key:                  r.Key,
				Optional:             r.Optional,
			}
		}
		out = append(out, kev)
	}
	return out
}

// withDefaults fills c's shared fields from defaults: Env and VolumeMounts
// merge (c wins by env name and by mount path), EnvFrom and SecurityContext
// apply when c has none.
func withDefaults(c, defaults types.Container) types.Container {
	if len(defaults.Env) > 0 {
		own := map[string]bool{}
		for _, e := range c.Env {
			own[e.Name] = true
		}
		var env []types.EnvVar
		for _, e := range defaults.Env {
			if !own[e.Name] {
				env = append(env, e)
			}
		}
		c.Env = append(env, c.Env...)
	}
	if c.EnvFrom == nil {
		c.EnvFrom = defaults.EnvFrom
	}
	if c.SecurityContext == nil {
		c.SecurityContext = defaults.SecurityContext
	}
	if len(defaults.VolumeMounts) > 0 {
		own := map[string]bool{}
		for _, m := range c.VolumeMounts {
			own[m.MountPath] = true
		}
		var mounts []types.VolumeMount
		for _, m := range defaults.VolumeMounts {
			if !own[m.MountPath] {
				mounts = append(mounts, m)
			}
		}
		c.VolumeMounts = append(mounts, c.VolumeMounts...)
	}
	return c
}
