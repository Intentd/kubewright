package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8stypes "k8s.io/apimachinery/pkg/types"
)

func TestOwnerReferences(t *testing.T) {
	owner := &unstructured.Unstructured{}
	owner.SetAPIVersion("demo.example.com/v1alpha1")
	owner.SetKind("Application")
	owner.SetName("my-app")
	owner.SetUID("test-uid")

	refs := OwnerReferences(owner, "")

	require.Len(t, refs, 1)
	assert.Equal(t, "demo.example.com/v1alpha1", refs[0].APIVersion)
	assert.Equal(t, "Application", refs[0].Kind)
	assert.Equal(t, "my-app", refs[0].Name)
	assert.Equal(t, k8stypes.UID("test-uid"), refs[0].UID)
	assert.True(t, *refs[0].Controller)
	assert.True(t, *refs[0].BlockOwnerDeletion)
}

// A namespaced owner can own objects only in its own namespace; anything else
// would be garbage-collected at once, so it gets no reference.
func TestOwnerReferences_Namespaces(t *testing.T) {
	owner := func(ns string) *unstructured.Unstructured {
		u := &unstructured.Unstructured{}
		u.SetAPIVersion("demo.example.com/v1")
		u.SetKind("App")
		u.SetName("app")
		u.SetNamespace(ns)
		u.SetUID("uid")
		return u
	}
	cases := []struct {
		name           string
		ownerNS, objNS string
		wantRef        bool
	}{
		{"namespaced owner, same namespace", "team", "team", true},
		{"namespaced owner, other namespace", "team", "other", false},
		{"namespaced owner, cluster-scoped object", "team", "", false},
		{"cluster-scoped owner, namespaced object", "", "other", true},
		{"cluster-scoped owner, cluster-scoped object", "", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			refs := OwnerReferences(owner(c.ownerNS), c.objNS)
			assert.Equal(t, c.wantRef, refs != nil)
			assert.Equal(t, c.wantRef, CanOwn(owner(c.ownerNS), c.objNS))
		})
	}

	noKind := &unstructured.Unstructured{}
	noKind.SetUID("uid")
	assert.Nil(t, OwnerReferences(noKind, ""), "an owner without apiVersion and kind")
}

func TestOwnedBy(t *testing.T) {
	owner := &unstructured.Unstructured{}
	owner.SetAPIVersion("demo.example.com/v1")
	owner.SetKind("App")
	owner.SetUID("owner-uid")

	obj := &unstructured.Unstructured{}
	obj.SetOwnerReferences(OwnerReferences(owner, ""))
	obj.SetLabels(map[string]string{"app.example.com/owner": "shop", "tier": "web"})

	other := &unstructured.Unstructured{}
	other.SetUID("other-uid")

	assert.True(t, OwnedByReference(owner)(obj), "reference to owner")
	assert.False(t, OwnedByReference(other)(obj), "reference to another owner")

	assert.True(t, OwnedByLabels(map[string]string{"app.example.com/owner": "shop"})(obj), "matching label")
	assert.False(t, OwnedByLabels(map[string]string{"app.example.com/owner": "billing"})(obj), "label with another value")
	assert.False(t, OwnedByLabels(map[string]string{"app.example.com/owner": "shop", "team": "x"})(obj), "one label missing")
	assert.False(t, OwnedByLabels(nil)(obj), "empty want matches nothing")

	obj.SetAnnotations(map[string]string{"app.example.com/owned-by": "team/shop"})
	assert.True(t, OwnedByAnnotations(map[string]string{"app.example.com/owned-by": "team/shop"})(obj), "matching annotation")
	assert.False(t, OwnedByAnnotations(map[string]string{"app.example.com/owned-by": "team/billing"})(obj), "annotation with another value")
	assert.False(t, OwnedByAnnotations(map[string]string{"app.example.com/owner": "shop"})(obj), "a label is not an annotation")
	assert.False(t, OwnedByAnnotations(nil)(obj), "empty want matches nothing")
}
