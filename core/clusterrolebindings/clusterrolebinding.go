package clusterrolebindings

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

// Spec describes a ClusterRoleBinding.
type Spec struct {
	// Name — the ClusterRoleBinding's name. Required.
	Name string
	// Labels — set on the object.
	Labels map[string]string
	// Annotations — set on the object.
	Annotations map[string]string
	// RoleRef — the ClusterRole to bind. It can't change after creation; Apply
	// replaces the binding when it does.
	RoleRef rbacv1.RoleRef
	// Subjects — the users, groups and service accounts granted the role.
	Subjects []rbacv1.Subject
}

// Create creates the ClusterRoleBinding spec describes, owned by owner, unless it already
// exists (Unchanged).
func Create(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("clusterrolebinding.Create: invalid spec: %w", err)
	}

	_, err := kube.RbacV1().ClusterRoleBindings().Get(ctx, spec.Name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return "", fmt.Errorf("clusterrolebinding.Create: checking existence of %q: %w", spec.Name, err)
	}
	if err == nil {
		return types.Unchanged, nil
	}

	crb := Build(owner, spec)

	_, err = kube.RbacV1().ClusterRoleBindings().Create(ctx, crb, metav1.CreateOptions{FieldManager: opts.FieldManager})
	if err != nil {
		return "", fmt.Errorf("clusterrolebinding.Create: creating %q: %w", spec.Name, err)
	}

	return types.Created, nil
}

// Apply writes the ClusterRoleBinding spec describes with server-side apply.
// A changed roleRef can't be applied, so the binding is replaced.
func Apply(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {

	crb := Build(owner, spec)
	crb.TypeMeta = metav1.TypeMeta{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRoleBinding"}

	if _, err := apply.Typed(ctx, kube.RbacV1().ClusterRoleBindings(), crb, opts); err != nil {
		if errors.IsInvalid(err) {
			// roleRef is immutable — delete and recreate.
			if delErr := kube.RbacV1().ClusterRoleBindings().Delete(ctx, spec.Name, metav1.DeleteOptions{}); delErr != nil && !errors.IsNotFound(delErr) {
				return "", fmt.Errorf("clusterrolebinding.Apply: deleting stale binding %q: %w", spec.Name, delErr)
			}
			if _, err := Create(ctx, kube, owner, spec, opts); err != nil {
				return "", err
			}
			return types.Replaced, nil
		}
		return "", fmt.Errorf("clusterrolebinding.Apply: %w", err)
	}

	return types.Applied, nil
}

// Delete deletes the ClusterRoleBinding (NotFound if it is already gone).
func Delete(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec) (types.Result, error) {

	err := kube.RbacV1().ClusterRoleBindings().Delete(ctx, spec.Name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", fmt.Errorf("clusterrolebinding.Delete: deleting %q: %w", spec.Name, err)
	}

	return types.Deleted, nil
}

// DeleteIfOwned deletes the ClusterRoleBinding only if owned reports true for it (NotOwned
// otherwise).
func DeleteIfOwned(ctx context.Context, kube kubernetes.Interface, name string, owned func(metav1.Object) bool) (types.Result, error) {
	existing, err := kube.RbacV1().ClusterRoleBindings().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", err
	}
	if !owned(existing) {
		return types.NotOwned, nil
	}
	if err := kube.RbacV1().ClusterRoleBindings().Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return "", err
	}
	return types.Deleted, nil
}

// ── Build and its helpers ─────────────────────────────────────────────────

// Build returns the ClusterRoleBinding spec describes, with an owner reference to owner
// where Kubernetes allows one (types.OwnerReferences). It reads nothing from
// the cluster.
func Build(owner types.Owner, spec Spec) *rbacv1.ClusterRoleBinding {
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        spec.Name,
			Labels:      spec.Labels,
			Annotations: spec.Annotations,
		},
		RoleRef:  spec.RoleRef,
		Subjects: spec.Subjects,
	}
	crb.OwnerReferences = types.OwnerReferences(owner, "")
	return crb
}

func validateSpec(spec Spec) error {
	if spec.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
