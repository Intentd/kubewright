package clusterroles

import (
	"context"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/intentd/kubewright/core/apply"
	"github.com/intentd/kubewright/core/types"
)

// Spec describes a ClusterRole.
type Spec struct {
	// Name — the ClusterRole's name. Required.
	Name string
	// Labels — set on the object.
	Labels map[string]string
	// Annotations — set on the object.
	Annotations map[string]string
	// Rules — what the role allows, across the cluster.
	Rules []rbacv1.PolicyRule
}

// Create creates the ClusterRole spec describes, owned by owner, unless it already
// exists (Unchanged).
func Create(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("clusterrole.Create: invalid spec: %w", err)
	}

	_, err := kube.RbacV1().ClusterRoles().Get(ctx, spec.Name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return "", fmt.Errorf("clusterrole.Create: checking existence of %q: %w", spec.Name, err)
	}
	if err == nil {
		return types.Unchanged, nil
	}

	cr := Build(owner, spec)

	_, err = kube.RbacV1().ClusterRoles().Create(ctx, cr, metav1.CreateOptions{FieldManager: opts.FieldManager})
	if err != nil {
		return "", fmt.Errorf("clusterrole.Create: creating %q: %w", spec.Name, err)
	}

	return types.Created, nil
}

// Apply writes the ClusterRole spec describes with server-side apply.
func Apply(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {

	cr := Build(owner, spec)
	cr.TypeMeta = metav1.TypeMeta{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRole"}

	if _, err := apply.Typed(ctx, kube.RbacV1().ClusterRoles(), cr, opts); err != nil {
		return "", fmt.Errorf("clusterrole.Apply: %w", err)
	}

	return types.Applied, nil
}

// Delete deletes the ClusterRole (NotFound if it is already gone).
func Delete(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec) (types.Result, error) {

	err := kube.RbacV1().ClusterRoles().Delete(ctx, spec.Name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", fmt.Errorf("clusterrole.Delete: deleting %q: %w", spec.Name, err)
	}

	return types.Deleted, nil
}

// DeleteIfOwned deletes the ClusterRole only if owned reports true for it (NotOwned
// otherwise).
func DeleteIfOwned(ctx context.Context, kube kubernetes.Interface, name string, owned func(metav1.Object) bool) (types.Result, error) {
	existing, err := kube.RbacV1().ClusterRoles().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", err
	}
	if !owned(existing) {
		return types.NotOwned, nil
	}
	if err := kube.RbacV1().ClusterRoles().Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return "", err
	}
	return types.Deleted, nil
}

// ── Build and its helpers ─────────────────────────────────────────────────

// Build returns the ClusterRole spec describes, with an owner reference to owner
// where Kubernetes allows one (types.OwnerReferences). It reads nothing from
// the cluster.
func Build(owner types.Owner, spec Spec) *rbacv1.ClusterRole {
	cr := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:        spec.Name,
			Labels:      spec.Labels,
			Annotations: spec.Annotations,
		},
		Rules: spec.Rules,
	}
	cr.OwnerReferences = types.OwnerReferences(owner, "")
	return cr
}

func validateSpec(spec Spec) error {
	if spec.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
