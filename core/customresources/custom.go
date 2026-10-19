package customresources

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"

	"github.com/intentd/kubewright/core/types"
	"github.com/intentd/kubewright/internal/shared"
)

// Spec describes a custom resource.
type Spec struct {
	// APIVersion is required and must be a group/version string (e.g. "foo.io/v1").
	// This field is used to derive the GroupVersionKind for REST mapping.
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`

	// Kind is required and must be a valid Kubernetes Kind (e.g. "Bar").
	// Used together with APIVersion to resolve the GVR for dynamic client calls.
	Kind string `json:"kind" yaml:"kind"`

	// Metadata is the subset of metav1.ObjectMeta a spec sets.
	// metadata.Name is required.
	// Namespace is required for namespaced CRDs; for cluster-scoped CRDs the
	// namespace field should be empty. Whether a CRD is namespaced is determined
	// by discovery/validation and not by this struct alone.
	Metadata types.CustomResourceMetadata `json:"metadata" yaml:"metadata"`

	// Spec — the object's spec. The API server validates it against the CRD.
	Spec map[string]any `json:"spec,omitempty" yaml:"spec,omitempty"`

	// Other — top-level fields besides metadata and spec, for kinds that keep
	// their configuration at the top level. A status is never written: it
	// belongs to the object's own controller.
	Other map[string]any `json:"-" yaml:",inline"`
}

// Create creates the custom resource spec describes, owned by owner, unless it already
// exists (Unchanged).
func Create(ctx context.Context, dyn dynamic.Interface, mapper meta.RESTMapper, owner types.Owner, spec Spec, opts types.Options, mutate func(*unstructured.Unstructured)) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("custom.Create: %w", err)
	}

	name := spec.Metadata.Name
	// Resolve namespace (owner may provide defaulting)
	namespace := shared.ResolveNamespace(owner, spec.Metadata.Namespace)

	// Build GVR from APIVersion/Kind
	gvk, err := shared.GVKFromFields(spec.APIVersion, spec.Kind)
	if err != nil {
		return "", fmt.Errorf("custom.Create: invalid GVK: %w", err)
	}

	resourceIfc, err := shared.ResourceFor(dyn, mapper, gvk, namespace)
	if err != nil {
		return "", fmt.Errorf("custom.Create: resolving GVR for %s: %w", gvk.String(), err)
	}

	// Check existence
	_, err = resourceIfc.Get(ctx, name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return "", fmt.Errorf("custom.Create: checking existence of %q: %w", name, err)
	}
	if err == nil {
		return types.Unchanged, nil
	}

	// Build unstructured object
	obj := Build(owner, spec)
	if mutate != nil {
		mutate(obj)
	}

	// Create
	_, err = resourceIfc.Create(ctx, obj, metav1.CreateOptions{FieldManager: opts.FieldManager})
	if err != nil {
		return "", fmt.Errorf("custom.Create: creating %q in %q: %w", name, namespace, err)
	}

	return types.Created, nil
}

// Apply writes the custom resource spec describes with server-side apply.
// It creates the object when it does not exist.
func Apply(ctx context.Context, dyn dynamic.Interface, mapper meta.RESTMapper, owner types.Owner, spec Spec, opts types.Options, mutate func(*unstructured.Unstructured)) (types.Result, error) {
	if err := validateSpec(spec); err != nil {
		return "", fmt.Errorf("custom.Apply: %w", err)
	}
	name := spec.Metadata.Name

	namespace := shared.ResolveNamespace(owner, spec.Metadata.Namespace)

	gvk, err := shared.GVKFromFields(spec.APIVersion, spec.Kind)
	if err != nil {
		return "", fmt.Errorf("custom.Apply: invalid GVK: %w", err)
	}

	resourceIfc, err := shared.ResourceFor(dyn, mapper, gvk, namespace)
	if err != nil {
		return "", fmt.Errorf("custom.Apply: resolving GVR for %s: %w", gvk.String(), err)
	}

	existing, err := resourceIfc.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return Create(ctx, dyn, mapper, owner, spec, opts, mutate)
		}
		return "", fmt.Errorf("custom.Apply: getting %q: %w", name, err)
	}

	// Build desired object
	desired := Build(owner, spec)

	// Compare spec and top-level Other fields for drift.
	existingSpec, _, _ := unstructured.NestedFieldNoCopy(existing.Object, "spec")
	desiredSpec, _, _ := unstructured.NestedFieldNoCopy(desired.Object, "spec")

	needUpdate := !reflect.DeepEqual(existingSpec, desiredSpec)

	// Also check other top-level fields (non-spec/status/metadata) for drift.
	if !needUpdate {
		for k, v := range desired.Object {
			if k == "apiVersion" || k == "kind" || k == "metadata" || k == "spec" || k == "status" {
				continue
			}
			if !reflect.DeepEqual(existing.Object[k], v) {
				needUpdate = true
				break
			}
		}
	}

	if !needUpdate {
		return types.Unchanged, nil
	}

	// Prepare updated object: merge desired spec into existing object.
	// Use direct map assignment (not SetNestedField) to avoid DeepCopyJSONValue
	// panicking on non-JSON-safe types (int, int64) produced by YAML unmarshalling.
	updated := existing.DeepCopy()
	if mutate != nil {
		mutate(updated)
	}
	if desired.Object["spec"] != nil {
		updated.Object["spec"] = shared.ToJSONSafe(desired.Object["spec"])
	}
	// Merge other top-level fields from desired (conservative: overwrite)
	for k, v := range desired.Object {
		if k == "apiVersion" || k == "kind" || k == "metadata" || k == "spec" || k == "status" {
			continue
		}
		updated.Object[k] = shared.ToJSONSafe(v)
	}

	_, err = resourceIfc.Update(ctx, updated, metav1.UpdateOptions{FieldManager: opts.FieldManager})
	if errors.IsConflict(err) {
		// The object changed between Get and Update: re-read it and retry once.
		latest, getErr := resourceIfc.Get(ctx, name, metav1.GetOptions{})
		if getErr != nil {
			return "", fmt.Errorf("custom.Apply: refreshing %q after conflict: %w", name, getErr)
		}
		updated = latest.DeepCopy()
		if desired.Object["spec"] != nil {
			updated.Object["spec"] = shared.ToJSONSafe(desired.Object["spec"])
		}
		for k, v := range desired.Object {
			if k == "apiVersion" || k == "kind" || k == "metadata" || k == "spec" || k == "status" {
				continue
			}
			updated.Object[k] = shared.ToJSONSafe(v)
		}
		_, err = resourceIfc.Update(ctx, updated, metav1.UpdateOptions{FieldManager: opts.FieldManager})
	}
	if err != nil {
		return "", fmt.Errorf("custom.Apply: updating %q: %w", name, err)
	}

	return types.Applied, nil
}

// DeleteIfOwned deletes the custom resource only if owned reports true for it (NotOwned
// otherwise).
func DeleteIfOwned(ctx context.Context, dyn dynamic.Interface, mapper meta.RESTMapper, name, namespace, apiVersion, kind string, owned func(metav1.Object) bool) (types.Result, error) {
	// Build GVK and resolve GVR
	gvk, err := shared.GVKFromFields(apiVersion, kind)
	if err != nil {
		return "", fmt.Errorf("custom.DeleteIfOwned: invalid GVK: %w", err)
	}

	resourceIfc, err := shared.ResourceFor(dyn, mapper, gvk, namespace)
	if err != nil {
		return "", fmt.Errorf("custom.DeleteIfOwned: resolving GVR for %s: %w", gvk.String(), err)
	}

	existing, err := resourceIfc.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return types.NotFound, nil
		}
		return "", fmt.Errorf("custom.DeleteIfOwned: getting %q: %w", name, err)
	}

	// Only delete if we own it
	if !owned(existing) {
		return types.NotOwned, nil
	}

	if err := resourceIfc.Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return "", fmt.Errorf("custom.DeleteIfOwned: deleting %q: %w", name, err)
	}

	return types.Deleted, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// Build returns the custom resource spec describes, owned by owner, in
// spec.Metadata.Namespace or else the owner's namespace. It reads nothing
// from the cluster.
func Build(owner types.Owner, spec Spec) *unstructured.Unstructured {
	namespace := shared.ResolveNamespace(owner, spec.Metadata.Namespace)
	if n := spec.Metadata.Namespaced; n != nil && !*n {
		namespace = "" // cluster-scoped
	}
	u := &unstructured.Unstructured{}
	u.SetAPIVersion(spec.APIVersion)
	u.SetKind(spec.Kind)
	u.SetName(spec.Metadata.Name)
	if namespace != "" {
		u.SetNamespace(namespace)
	}

	// Labels: copy the declared labels rather than
	// mutate spec.Metadata.Labels directly — that map may be shared/reused
	// (e.g. across forEach expansions of the same source).
	lbls := make(map[string]string, len(spec.Metadata.Labels)+3)
	for k, v := range spec.Metadata.Labels {
		lbls[k] = v
	}
	u.SetLabels(lbls)

	// Annotations: copy for the same reason as Labels above.
	if len(spec.Metadata.Annotations) > 0 {
		annotations := make(map[string]string, len(spec.Metadata.Annotations))
		for k, v := range spec.Metadata.Annotations {
			annotations[k] = v
		}
		u.SetAnnotations(annotations)
	}

	// Owner reference, so Kubernetes deletes it with owner, when it can hold
	// one (see types.OwnerReferences).
	if refs := types.OwnerReferences(owner, namespace); refs != nil {
		u.SetOwnerReferences(refs)
	}

	// Other top-level fields (non-core) from the spec declaration.
	for k, v := range spec.Other {
		if k == "apiVersion" || k == "kind" || k == "metadata" || k == "spec" || k == "status" {
			continue
		}
		u.Object[k] = v
	}

	if spec.Spec != nil {
		u.Object["spec"] = spec.Spec
	}

	// yaml.v3 decodes integer literals as Go int, but k8s fake client's DeepCopy
	// only handles JSON-native types (float64 for numbers, not int). A JSON
	// round-trip normalises the entire object map before it reaches the client.
	if raw, err := json.Marshal(u.Object); err == nil {
		var normalised map[string]any
		if json.Unmarshal(raw, &normalised) == nil {
			u.Object = normalised
		}
	}

	return u
}

// ── Internal helpers ──────────────────────────────────────────────────────────

func validateSpec(spec Spec) error {
	if spec.Metadata.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
