package serviceaccounts

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/intentd/kubewright/core/types"
	"github.com/intentd/kubewright/internal/shared"
)

// Spec describes a ServiceAccount.
type Spec struct {
	// Name — ServiceAccount name. Required.
	Name string

	// Namespace — target namespace. Required.
	Namespace string

	// Labels — applied to ServiceAccount metadata.
	Labels map[string]string

	// Annotations — set on the object.
	Annotations map[string]string
}

// Create creates the ServiceAccount spec describes, owned by owner, unless it already
// exists (Unchanged).
// An existing ServiceAccount is left as it is: it has nothing that drifts.
func Create(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("serviceaccount.Create: invalid spec: %w", err)
	}

	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	_, err := kube.CoreV1().ServiceAccounts(namespace).Get(ctx, spec.Name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return "", fmt.Errorf("serviceaccount.Create: checking existence of %q: %w", spec.Name, err)
	}
	if err == nil {
		return types.Unchanged, nil
	}

	sa := Build(owner, spec)

	_, err = kube.CoreV1().ServiceAccounts(namespace).Create(ctx, sa, metav1.CreateOptions{FieldManager: opts.FieldManager})
	if err != nil {
		return "", fmt.Errorf("serviceaccount.Create: creating serviceaccount %q in %q: %w", spec.Name, namespace, err)
	}

	return types.Created, nil
}

// Delete deletes the ServiceAccount (NotFound if it is already gone).
func Delete(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec) (types.Result, error) {
	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	err := kube.CoreV1().ServiceAccounts(namespace).Delete(ctx, spec.Name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", fmt.Errorf("serviceaccount.Delete: deleting serviceaccount %q in %q: %w", spec.Name, namespace, err)
	}

	return types.Deleted, nil
}

// DeleteIfOwned deletes the ServiceAccount only if owned reports true for it (NotOwned
// otherwise).
func DeleteIfOwned(ctx context.Context, kube kubernetes.Interface, name, namespace string, owned func(metav1.Object) bool) (types.Result, error) {
	existing, err := kube.CoreV1().ServiceAccounts(namespace).
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
	if err := kube.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return "", err
	}
	return types.Deleted, nil
}

// ── Build and its helpers ─────────────────────────────────────────────────

// Build returns the ServiceAccount spec describes, in spec.Namespace or else the owner's namespace, with an owner reference to owner
// where Kubernetes allows one (types.OwnerReferences). It reads nothing from
// the cluster.
func Build(owner types.Owner, spec Spec) *corev1.ServiceAccount {
	namespace := shared.ResolveNamespace(owner, spec.Namespace)
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:            spec.Name,
			Namespace:       namespace,
			Labels:          spec.Labels,
			Annotations:     spec.Annotations,
			OwnerReferences: types.OwnerReferences(owner, namespace),
		},
	}
}

func validateSpec(spec Spec) error {
	if spec.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
