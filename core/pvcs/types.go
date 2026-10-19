package pvcs

// Spec describes a PersistentVolumeClaim.
type Spec struct {
	// Name — the claim's name. Required.
	Name string
	// Namespace — where to create it. Empty means the owner's.
	Namespace string
	// StorageClassName — the storage class to provision from. Empty uses the
	// cluster's default class.
	StorageClassName string
	// AccessModes — ReadWriteOnce, ReadOnlyMany, ReadWriteMany or
	// ReadWriteOncePod.
	AccessModes []string
	// Storage — the size requested, e.g. "10Gi". Required.
	Storage string
	// VolumeMode — Filesystem (the default) or Block.
	VolumeMode string
	// VolumeName — bind to this existing PersistentVolume instead of
	// provisioning one.
	VolumeName string
	// Labels — set on the object.
	Labels map[string]string

	// Annotations — set on the object.
	Annotations map[string]string
}
