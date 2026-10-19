package pods

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/intentd/kubewright/core/apply"
	"github.com/intentd/kubewright/core/types"
	"github.com/intentd/kubewright/internal/shared"
)

// Create creates the Pod spec describes, owned by owner, unless it already
// exists (Unchanged).
func Create(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("pod.Create: invalid spec: %w", err)
	}

	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	_, err := kube.CoreV1().Pods(namespace).Get(ctx, spec.Name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return "", fmt.Errorf("pod.Create: checking existence of %q: %w", spec.Name, err)
	}
	if err == nil {
		return types.Unchanged, nil
	}

	pod := Build(owner, spec)

	_, err = kube.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{FieldManager: opts.FieldManager})
	if err != nil {
		return "", fmt.Errorf("pod.Create: creating pod %q in %q: %w", spec.Name, namespace, err)
	}

	return types.Created, nil
}

// Apply writes the Pod spec describes with server-side apply.
// A Pod's spec is mostly immutable: when the change can't be applied, the Pod
// is replaced.
func Apply(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("pod.Apply: invalid spec: %w", err)
	}

	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	p := Build(owner, spec)
	p.TypeMeta = metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"}

	if _, err := apply.Typed(ctx, kube.CoreV1().Pods(namespace), p, opts); err != nil {
		if errors.IsInvalid(err) {
			if delErr := kube.CoreV1().Pods(namespace).Delete(ctx, spec.Name, metav1.DeleteOptions{}); delErr != nil && !errors.IsNotFound(delErr) {
				return "", fmt.Errorf("pod.Apply: deleting stale pod %q: %w", spec.Name, delErr)
			}
			if _, err := Create(ctx, kube, owner, spec, opts); err != nil {
				return "", err
			}
			return types.Replaced, nil
		}
		return "", fmt.Errorf("pod.Apply: %w", err)
	}

	return types.Applied, nil
}

// Delete deletes the Pod (NotFound if it is already gone).
func Delete(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec) (types.Result, error) {
	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	err := kube.CoreV1().Pods(namespace).Delete(ctx, spec.Name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", fmt.Errorf("pod.Delete: deleting pod %q in %q: %w", spec.Name, namespace, err)
	}

	return types.Deleted, nil
}

// DeleteIfOwned deletes the Pod only if owned reports true for it (NotOwned
// otherwise).
func DeleteIfOwned(ctx context.Context, kube kubernetes.Interface, name, namespace string, owned func(metav1.Object) bool) (types.Result, error) {
	existing, err := kube.CoreV1().Pods(namespace).
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
	if err := kube.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return "", err
	}
	return types.Deleted, nil
}

// ── Build and its helpers ─────────────────────────────────────────────────

// Build returns the Pod spec describes, in spec.Namespace or else the owner's namespace, with an owner reference to owner
// where Kubernetes allows one (types.OwnerReferences). It reads nothing from
// the cluster.
func Build(owner types.Owner, spec Spec) *corev1.Pod {
	namespace := shared.ResolveNamespace(owner, spec.Namespace)
	containers, inits := shared.PodContainers(mainContainer(spec), spec.Containers, spec.InitContainers)
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:            spec.Name,
			Namespace:       namespace,
			Labels:          spec.Labels,
			Annotations:     spec.Annotations,
			OwnerReferences: types.OwnerReferences(owner, namespace),
		},
		Spec: corev1.PodSpec{
			ImagePullSecrets:   shared.ToPullSecrets(spec.ImagePullSecrets),
			ServiceAccountName: spec.ServiceAccountName,
			NodeSelector:       spec.NodeSelector,
			RestartPolicy:      spec.RestartPolicy,
			InitContainers:     inits,
			Containers:         containers,
		},
	}

	// Pod security
	if spec.PodSecurity != nil {
		pod.Spec.SecurityContext = shared.BuildPodSecurityContext(spec.PodSecurity)
	}

	// Volumes
	if vols := shared.BuildVolumes(spec.Volumes); len(vols) > 0 {
		pod.Spec.Volumes = vols
	}

	return pod
}

// mainContainer is the container the flat fields describe.
func mainContainer(spec Spec) types.Container {
	return types.Container{
		Name:            spec.Name,
		Image:           spec.Image,
		Port:            int32(spec.Port),
		Protocol:        spec.Protocol,
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
	if err := shared.ValidateRestartPolicy(spec.RestartPolicy, corev1.RestartPolicyAlways, corev1.RestartPolicyOnFailure, corev1.RestartPolicyNever); err != nil {
		return err
	}
	return shared.ValidateContainers(mainContainer(spec), spec.Containers, spec.InitContainers)
}
