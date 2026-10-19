package deployments

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/intentd/kubewright/core/apply"
	"github.com/intentd/kubewright/core/types"
	"github.com/intentd/kubewright/internal/shared"
)

// Create creates the Deployment spec describes, owned by owner, unless it already
// exists (Unchanged).
func Create(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("deployment.Create: invalid spec: %w", err)
	}

	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	_, err := kube.AppsV1().Deployments(namespace).Get(ctx, spec.Name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return "", fmt.Errorf("deployment.Create: checking existence of %q: %w", spec.Name, err)
	}
	if err == nil {
		return types.Unchanged, nil
	}

	deployment := Build(owner, spec)

	_, err = kube.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{FieldManager: opts.FieldManager})
	if err != nil {
		return "", fmt.Errorf("deployment.Create: creating deployment %q in %q: %w", spec.Name, namespace, err)
	}

	return types.Created, nil
}

// Apply writes the Deployment spec describes with server-side apply.
func Apply(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("deployment.Apply: invalid spec: %w", err)
	}

	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	d := Build(owner, spec)
	d.TypeMeta = metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"}

	if _, err := apply.Typed(ctx, kube.AppsV1().Deployments(namespace), d, opts); err != nil {
		return "", fmt.Errorf("deployment.Apply: %w", err)
	}

	return types.Applied, nil
}

// Delete deletes the Deployment (NotFound if it is already gone).
func Delete(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec) (types.Result, error) {
	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	err := kube.AppsV1().Deployments(namespace).Delete(ctx, spec.Name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", fmt.Errorf("deployment.Delete: deleting deployment %q in %q: %w", spec.Name, namespace, err)
	}

	return types.Deleted, nil
}

// DeleteIfOwned deletes the Deployment only if owned reports true for it (NotOwned
// otherwise).
func DeleteIfOwned(ctx context.Context, kube kubernetes.Interface, name, namespace string, owned func(metav1.Object) bool) (types.Result, error) {
	existing, err := kube.AppsV1().Deployments(namespace).
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
	if err := kube.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return "", err
	}
	return types.Deleted, nil
}

// ── Build and its helpers ─────────────────────────────────────────────────

// Build returns the Deployment spec describes, in spec.Namespace or else the owner's namespace, with an owner reference to owner
// where Kubernetes allows one (types.OwnerReferences). It reads nothing from
// the cluster.
func Build(owner types.Owner, spec Spec) *appsv1.Deployment {
	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	replicas := spec.Replicas
	containers, inits := shared.PodContainers(mainContainer(spec), spec.Containers, spec.InitContainers)

	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            spec.Name,
			Namespace:       namespace,
			Labels:          spec.Labels,
			Annotations:     spec.Annotations,
			OwnerReferences: types.OwnerReferences(owner, namespace),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: shared.PodSelector(spec.Selector, spec.Name),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: shared.PodLabels(spec.Labels, shared.PodSelector(spec.Selector, spec.Name)),
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets:   shared.ToPullSecrets(spec.ImagePullSecrets),
					ServiceAccountName: spec.ServiceAccountName,
					NodeSelector:       spec.NodeSelector,
					InitContainers:     inits,
					Containers:         containers,
				},
			},
		},
	}

	// Rolling update strategy
	if spec.RollingUpdate != nil {
		d.Spec.Strategy = shared.BuildDeploymentRollingUpdateStrategy(spec.RollingUpdate)
	}

	// Pod security
	if spec.PodSecurity != nil {
		d.Spec.Template.Spec.SecurityContext = shared.BuildPodSecurityContext(spec.PodSecurity)
	}

	// Volumes
	if vols := shared.BuildVolumes(spec.Volumes); len(vols) > 0 {
		d.Spec.Template.Spec.Volumes = vols
	}

	return d
}

// mainContainer is the container the flat fields describe.
func mainContainer(spec Spec) types.Container {
	return types.Container{
		Name:            spec.Name,
		Image:           spec.Image,
		Command:         spec.Command,
		Port:            spec.Port,
		Protocol:        spec.Protocol,
		Env:             spec.Env,
		EnvFrom:         spec.EnvFrom,
		Resources:       spec.Resources,
		Probes:          spec.Probes,
		SecurityContext: spec.SecurityContext,
		VolumeMounts:    spec.VolumeMounts,
	}
}

func validateSpec(spec Spec) error {
	if spec.Name == "" {
		return fmt.Errorf("missing required fields: [name]")
	}
	return shared.ValidateContainers(mainContainer(spec), spec.Containers, spec.InitContainers)
}
