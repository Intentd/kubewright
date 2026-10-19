package cronjobs

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/intentd/kubewright/core/types"
	"github.com/intentd/kubewright/internal/testutil"
)

func spec() Spec {
	return Spec{Name: "web", Image: "busybox:1.36", Schedule: "*/5 * * * *"}
}

func TestBuild(t *testing.T) {
	o := Build(testutil.Owner(), spec())

	if o.Name != "web" {
		t.Errorf("name = %q, want web", o.Name)
	}
	if o.Namespace != "team" {
		t.Errorf("namespace = %q, want the owner's", o.Namespace)
	}
	if !testutil.Owned(o) {
		t.Errorf("owner references = %+v, want one to the owner", o.OwnerReferences)
	}
	if o.Spec.Schedule != "*/5 * * * *" || o.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image != "busybox:1.36" {
		t.Errorf("schedule %q image %q, want the spec's", o.Spec.Schedule, o.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image)
	}
}

func TestLifecycle(t *testing.T) {
	owner, s := testutil.Owner(), spec()
	testutil.Lifecycle{
		Create: func(ctx context.Context, cs kubernetes.Interface) (types.Result, error) {
			return Create(ctx, cs, owner, s, types.Options{})
		},
		Apply: func(ctx context.Context, cs kubernetes.Interface) (types.Result, error) {
			return Apply(ctx, cs, owner, s, types.Options{})
		},
		Delete: func(ctx context.Context, cs kubernetes.Interface) (types.Result, error) {
			return Delete(ctx, cs, owner, s)
		},
		DeleteIfOwned: func(ctx context.Context, cs kubernetes.Interface, owned func(metav1.Object) bool) (types.Result, error) {
			return DeleteIfOwned(ctx, cs, "web", "team", owned)
		},
		Get: func(ctx context.Context, cs kubernetes.Interface) (metav1.Object, error) {
			return cs.BatchV1().CronJobs("team").Get(ctx, "web", metav1.GetOptions{})
		},
	}.Check(t)
}

func TestBuild_Metadata(t *testing.T) {
	s := spec()
	s.Labels = map[string]string{"app.example.com/owner": "team.app"}
	s.Annotations = map[string]string{"app.example.com/owned-by": "team/app"}
	o := Build(testutil.Owner(), s)

	if o.Labels["app.example.com/owner"] != "team.app" {
		t.Errorf("labels = %v, want the spec's", o.Labels)
	}
	if o.Annotations["app.example.com/owned-by"] != "team/app" {
		t.Errorf("annotations = %v, want the spec's", o.Annotations)
	}
}

func TestBuild_PodPlacementAndRestartPolicy(t *testing.T) {
	s := spec()
	s.NodeSelector = map[string]string{"pool": "batch"}
	s.ServiceAccountName = "runner"
	pod := Build(testutil.Owner(), s).Spec.JobTemplate.Spec.Template.Spec
	if pod.NodeSelector["pool"] != "batch" || pod.ServiceAccountName != "runner" {
		t.Errorf("nodeSelector = %v, serviceAccountName = %q; want pool=batch and runner", pod.NodeSelector, pod.ServiceAccountName)
	}
	if pod.RestartPolicy != corev1.RestartPolicyOnFailure {
		t.Errorf("default restartPolicy = %q, want OnFailure", pod.RestartPolicy)
	}

	s.RestartPolicy = corev1.RestartPolicyNever
	if got := Build(testutil.Owner(), s).Spec.JobTemplate.Spec.Template.Spec.RestartPolicy; got != corev1.RestartPolicyNever {
		t.Errorf("restartPolicy = %q, want Never", got)
	}
	s.RestartPolicy = corev1.RestartPolicyAlways
	if err := validateSpec(s); err == nil {
		t.Error("restartPolicy Always was accepted; Kubernetes allows only OnFailure or Never here")
	}
}
