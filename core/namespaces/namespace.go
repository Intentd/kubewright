package namespaces

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/intentd/kubewright/core/apply"
	"github.com/intentd/kubewright/core/types"
)

// Spec describes a Namespace.
type Spec struct {
	// Name — name. Required.
	Name string

	// Labels — applied to Namespace metadata.
	Labels map[string]string

	// Annotations — applied to Namespace metadata.
	Annotations map[string]string

	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// optional
	Finalizers []string
}

// Create creates the Namespace spec describes, owned by owner, unless it already
// exists (Unchanged).
// An existing Namespace is left as it is.
func Create(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("namespace.Create: invalid spec: %w", err)
	}

	_, err := kube.CoreV1().Namespaces().Get(ctx, spec.Name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return "", fmt.Errorf("namespace.Create: checking existence of %q: %w", spec.Name, err)
	}
	if err == nil {
		return types.Unchanged, nil
	}

	sa := Build(owner, spec)

	_, err = kube.CoreV1().Namespaces().Create(ctx, sa, metav1.CreateOptions{FieldManager: opts.FieldManager})
	if err != nil {
		return "", fmt.Errorf("namespace.Create: creating namespace %q: %w", spec.Name, err)
	}

	return types.Created, nil
}

// Apply writes the Namespace spec describes with server-side apply.
func Apply(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("namespace.Apply: invalid spec: %w", err)
	}

	ns := Build(owner, spec)
	ns.TypeMeta = metav1.TypeMeta{APIVersion: "v1", Kind: "Namespace"}

	if _, err := apply.Typed(ctx, kube.CoreV1().Namespaces(), ns, opts); err != nil {
		return "", fmt.Errorf("namespace.Apply: %w", err)
	}

	return types.Applied, nil
}

// Delete deletes the Namespace (NotFound if it is already gone).
func Delete(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec) (types.Result, error) {

	err := kube.CoreV1().Namespaces().Delete(ctx, spec.Name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", fmt.Errorf("namespace.Delete: deleting namespace %q: %w", spec.Name, err)
	}

	return types.Deleted, nil
}

// DeleteIfOwned deletes the Namespace only if owned reports true for it (NotOwned
// otherwise).
func DeleteIfOwned(ctx context.Context, kube kubernetes.Interface, name string, owned func(metav1.Object) bool) (types.Result, error) {
	existing, err := kube.CoreV1().Namespaces().
		Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", err
	}
	// Only delete if we own it
	if !owned(existing) {
		return types.NotOwned, nil
	}
	if err := kube.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return "", err
	}
	return types.Deleted, nil
}

// ── Build and its helpers ─────────────────────────────────────────────────

// Build returns the Namespace spec describes, with an owner reference to owner
// where Kubernetes allows one (types.OwnerReferences). It reads nothing from
// the cluster.
func Build(owner types.Owner, spec Spec) *corev1.Namespace {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        spec.Name,
			Labels:      spec.Labels,
			Annotations: spec.Annotations,
		},
	}
	// Only a cluster-scoped owner can own a Namespace; a namespaced owner leaves
	// it without a reference, to be found by labels.
	ns.OwnerReferences = types.OwnerReferences(owner, "")
	return ns
}

func validateSpec(spec Spec) error {
	if spec.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
