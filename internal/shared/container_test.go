package shared

import (
	"strings"
	"testing"

	"github.com/intentd/kubewright/core/types"
)

func flatContainer() types.Container {
	return types.Container{
		Name:         "web",
		Image:        "nginx",
		Port:         8080,
		Env:          []types.EnvVar{{Name: "LOG", Value: "info"}, {Name: "REGION", Value: "eu"}},
		VolumeMounts: []types.VolumeMount{{Name: "cfg", MountPath: "/etc/cfg"}},
	}
}

func TestPodContainers_FlatIsTheOneContainer(t *testing.T) {
	cs, inits := PodContainers(flatContainer(), nil, nil)
	if len(cs) != 1 || cs[0].Name != "web" || cs[0].Image != "nginx" || len(inits) != 0 {
		t.Fatalf("containers = %+v, inits = %+v; want the one flat container", cs, inits)
	}
}

func TestPodContainers_SharedFieldsAreDefaults(t *testing.T) {
	flat := flatContainer()
	flat.Image, flat.Port = "", 0 // main-only fields go in containers
	containers := []types.Container{
		{Name: "app", Image: "app:1", Env: []types.EnvVar{{Name: "LOG", Value: "debug"}}},
		{Name: "proxy", Image: "envoy", VolumeMounts: []types.VolumeMount{{Name: "own", MountPath: "/etc/cfg"}}},
	}
	inits := []types.Container{{Name: "migrate", Image: "app:1", Command: []string{"migrate"}}}

	cs, is := PodContainers(flat, containers, inits)
	if len(cs) != 2 || len(is) != 1 {
		t.Fatalf("got %d containers and %d init containers, want 2 and 1", len(cs), len(is))
	}
	env := func(i int) map[string]string {
		m := map[string]string{}
		for _, e := range cs[i].Env {
			m[e.Name] = e.Value
		}
		return m
	}
	if got := env(0); got["LOG"] != "debug" || got["REGION"] != "eu" || len(got) != 2 {
		t.Errorf("app env = %v, want its own LOG and the shared REGION", got)
	}
	if got := env(1); got["LOG"] != "info" {
		t.Errorf("proxy env = %v, want the shared LOG", got)
	}
	if m := cs[1].VolumeMounts; len(m) != 1 || m[0].Name != "own" {
		t.Errorf("proxy mounts = %+v, want its own mount at /etc/cfg to win", m)
	}
	if m := is[0].VolumeMounts; len(m) != 1 || m[0].Name != "cfg" {
		t.Errorf("init mounts = %+v, want the shared mount", m)
	}
	if len(is[0].Env) != 2 {
		t.Errorf("init env = %+v, want the shared env", is[0].Env)
	}
	if cs[0].Ports != nil || cs[0].ReadinessProbe != nil {
		t.Errorf("app = %+v, want no port or probe inherited", cs[0])
	}
}

func TestValidateContainers(t *testing.T) {
	ok := []types.Container{{Name: "app", Image: "app:1"}}
	cases := []struct {
		name       string
		flat       types.Container
		containers []types.Container
		inits      []types.Container
		want       string
	}{
		{"flat needs an image", types.Container{Name: "web"}, nil, nil, "image is required"},
		{"flat with an init", types.Container{Name: "web", Image: "nginx"}, nil, []types.Container{{Name: "init", Image: "busybox"}}, ""},
		{"init named like the flat container", types.Container{Name: "web", Image: "nginx"}, nil, []types.Container{{Name: "web", Image: "busybox"}}, "already used"},
		{"main-only fields next to containers", types.Container{Name: "web", Image: "nginx", Port: 80}, ok, nil, "image, port set next to containers"},
		{"shared fields next to containers", types.Container{Name: "web", Env: []types.EnvVar{{Name: "A"}}}, ok, nil, ""},
		{"container needs name and image", types.Container{Name: "web"}, []types.Container{{}}, nil, "name is required"},
		{"duplicate names", types.Container{Name: "web"}, []types.Container{{Name: "a", Image: "x"}, {Name: "a", Image: "y"}}, nil, "already used"},
		{"init container probes", types.Container{Name: "web", Image: "nginx"}, nil, []types.Container{{Name: "i", Image: "x", Probes: &types.ProbesConfig{}}}, "no probes"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateContainers(c.flat, c.containers, c.inits)
			if c.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("error = %v, want one containing %q", err, c.want)
			}
		})
	}
}

func TestValidateContainers_UnexpandedPresets(t *testing.T) {
	flat := types.Container{
		Name:      "web",
		Image:     "nginx",
		Resources: &types.ResourceRequirements{Preset: "small"},
		Probes:    &types.ProbesConfig{Readiness: &types.ProbeConfig{Path: "/", Preset: "fast"}},
	}
	err := ValidateContainers(flat, nil, nil)
	if err == nil || !strings.Contains(err.Error(), `resources preset "small" is not expanded`) ||
		!strings.Contains(err.Error(), `readiness probe preset "fast" is not expanded`) {
		t.Fatalf("error = %v, want both unexpanded presets named", err)
	}
}
