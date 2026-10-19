package secrets

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/intentd/kubewright/core/types"
	"github.com/intentd/kubewright/internal/testutil"
)

func spec() Spec {
	return Spec{Name: "web", StringData: map[string]string{"token": "abc"}}
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
	if o.StringData["token"] != "abc" && string(o.Data["token"]) != "abc" {
		t.Errorf("secret data = %v / %v, want token", o.StringData, o.Data)
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
			return cs.CoreV1().Secrets("team").Get(ctx, "web", metav1.GetOptions{})
		},
	}.Check(t)
}

func TestCopyToNamespaces(t *testing.T) {
	ctx := context.Background()
	cs := fake.NewClientset()
	copyAll := func() {
		t.Helper()
		if err := CopyToNamespaces(ctx, cs, testutil.Owner(), spec(), []string{"a", "b", ""}, types.Options{}); err != nil {
			t.Fatal(err)
		}
	}

	copyAll()
	for _, ns := range []string{"a", "b"} {
		if _, err := cs.CoreV1().Secrets(ns).Get(ctx, "web", metav1.GetOptions{}); err != nil {
			t.Errorf("no copy in %s: %v", ns, err)
		}
	}
	copyAll() // existing copies are left as they are
}

func TestCopyToNamespaces_FromSource(t *testing.T) {
	ctx := context.Background()
	cs := fake.NewClientset(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "src", Namespace: "shared"},
		Data:       map[string][]byte{"token": []byte("from-source")},
	})
	s := Spec{Name: "web", FromSecret: "src", FromNamespace: "shared"}
	if err := CopyToNamespaces(ctx, cs, testutil.Owner(), s, []string{"a"}, types.Options{}); err != nil {
		t.Fatal(err)
	}
	sec, err := cs.CoreV1().Secrets("a").Get(ctx, "web", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if string(sec.Data["token"]) != "from-source" {
		t.Errorf("data = %v, want the source's", sec.Data)
	}
}

// A copy into another namespace can't be owned by a namespaced owner: the
// garbage collector would delete it at once.
func TestBuild_OtherNamespace(t *testing.T) {
	s := spec()
	s.Namespace = "elsewhere"
	if o := Build(testutil.Owner(), s); len(o.OwnerReferences) != 0 {
		t.Errorf("owner references = %+v, want none in another namespace", o.OwnerReferences)
	}
	if o := Build(testutil.ClusterOwner(), s); !testutil.Owned(o) {
		t.Errorf("owner references = %+v, want one to a cluster-scoped owner", o.OwnerReferences)
	}
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
