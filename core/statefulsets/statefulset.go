package statefulsets

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/intentd/kubewright/core/apply"
	"github.com/intentd/kubewright/core/types"
	"github.com/intentd/kubewright/internal/shared"
)

// Create creates a StatefulSet owned by owner if it does not already exist.
func Create(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("statefulset.Create: invalid spec: %w", err)
	}
	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	_, err := kube.AppsV1().StatefulSets(namespace).Get(ctx, spec.Name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return "", fmt.Errorf("statefulset.Create: checking existence of %q: %w", spec.Name, err)
	}
	if err == nil {
		return types.Unchanged, nil
	}

	sts := Build(owner, spec)
	_, err = kube.AppsV1().StatefulSets(namespace).Create(ctx, sts, metav1.CreateOptions{FieldManager: opts.FieldManager})
	if err != nil {
		return "", fmt.Errorf("statefulset.Create: creating %q in %q: %w", spec.Name, namespace, err)
	}

	return types.Created, nil
}

// Apply creates or updates a StatefulSet using Server-Side Apply.
// Sends only the fields spec sets; defaults the API server fills in are left alone.
func Apply(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec, opts types.Options) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("statefulset.Apply: invalid spec: %w", err)
	}
	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	sts := Build(owner, spec)
	sts.TypeMeta = metav1.TypeMeta{APIVersion: "apps/v1", Kind: "StatefulSet"}

	if _, err := apply.Typed(ctx, kube.AppsV1().StatefulSets(namespace), sts, opts); err != nil {
		return "", fmt.Errorf("statefulset.Apply: %w", err)
	}

	return types.Applied, nil
}

// Delete deletes the StatefulSet if it exists.
func Delete(ctx context.Context, kube kubernetes.Interface, owner types.Owner, spec Spec) (types.Result, error) {
	namespace := shared.ResolveNamespace(owner, spec.Namespace)

	err := kube.AppsV1().StatefulSets(namespace).Delete(ctx, spec.Name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", fmt.Errorf("statefulset.Delete: %w", err)
	}
	return types.Deleted, nil
}

// DeleteIfOwned deletes the StatefulSet if it exists and owned reports true for it.
func DeleteIfOwned(ctx context.Context, kube kubernetes.Interface, name, namespace string, owned func(metav1.Object) bool) (types.Result, error) {
	existing, err := kube.AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", err
	}
	if !owned(existing) {
		return types.NotOwned, nil
	}
	if err := kube.AppsV1().StatefulSets(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return "", err
	}
	return types.Deleted, nil
}

// ── Build and its helpers ─────────────────────────────────────────────────

func resolveAccessModes(modes []string) []corev1.PersistentVolumeAccessMode {
	if len(modes) == 0 {
		return []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
	}
	out := make([]corev1.PersistentVolumeAccessMode, 0, len(modes))
	for _, m := range modes {
		switch m {
		case "ReadWriteMany":
			out = append(out, corev1.ReadWriteMany)
		case "ReadOnlyMany":
			out = append(out, corev1.ReadOnlyMany)
		case "ReadWriteOncePod":
			out = append(out, corev1.ReadWriteOncePod)
		default:
			out = append(out, corev1.ReadWriteOnce)
		}
	}
	return out
}

// Build returns the StatefulSet spec describes, owned by owner, in spec.Namespace or else the owner's namespace.
// It reads nothing from the cluster.
func Build(owner types.Owner, spec Spec) *appsv1.StatefulSet {
	ns := shared.ResolveNamespace(owner, spec.Namespace)

	replicas := spec.Replicas
	containers, inits := shared.PodContainers(mainContainer(spec), spec.Containers, spec.InitContainers)

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            spec.Name,
			Namespace:       ns,
			Labels:          spec.Labels,
			Annotations:     spec.Annotations,
			OwnerReferences: types.OwnerReferences(owner, ns),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &replicas,
			ServiceName: spec.ServiceName,
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
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{},
			PersistentVolumeClaimRetentionPolicy: &appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy{
				WhenDeleted: appsv1.PersistentVolumeClaimRetentionPolicyType(spec.VolumeClaimRetentionPolicy.WhenDeleted),
				WhenScaled:  appsv1.PersistentVolumeClaimRetentionPolicyType(spec.VolumeClaimRetentionPolicy.WhenScaled),
			},
			UpdateStrategy: func() appsv1.StatefulSetUpdateStrategy {
				if spec.RollingUpdate != nil {
					return shared.BuildStatefulSetUpdateStrategy(spec.RollingUpdate)
				}
				return appsv1.StatefulSetUpdateStrategy{Type: appsv1.OnDeleteStatefulSetStrategyType}
			}(),
			PodManagementPolicy: appsv1.ParallelPodManagement,
		},
	}

	// Pod security
	if spec.PodSecurity != nil {
		sts.Spec.Template.Spec.SecurityContext = shared.BuildPodSecurityContext(spec.PodSecurity)
	}

	// Volume claim templates mount on the first container, before its own mounts.
	var claimMounts []corev1.VolumeMount

	for _, vct := range spec.VolumeClaimTemplates {
		storageQty := resource.MustParse(vct.StorageSize)
		name := vct.Name
		if name == "" {
			name = "data"
		}
		mountPath := vct.MountPath
		if mountPath == "" {
			mountPath = "/data"
		}
		accessModes := resolveAccessModes(vct.AccessModes)
		claimMounts = append(claimMounts, corev1.VolumeMount{Name: name, MountPath: mountPath})
		sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes:      accessModes,
				StorageClassName: &vct.StorageClass,
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: storageQty,
					},
				},
			},
		})
	}

	if len(claimMounts) > 0 {
		first := &sts.Spec.Template.Spec.Containers[0]
		first.VolumeMounts = append(claimMounts, first.VolumeMounts...)
	}

	// Volumes
	if vols := shared.BuildVolumes(spec.Volumes); len(vols) > 0 {
		sts.Spec.Template.Spec.Volumes = vols
	}

	return sts
}

// mainContainer is the container the flat fields describe.
func mainContainer(spec Spec) types.Container {
	return types.Container{
		Name:            spec.Name,
		Image:           spec.Image,
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
