package pvs

// Spec describes a PersistentVolume. PersistentVolumes are cluster-scoped.
type Spec struct {
	// Name — the volume's name. Required.
	Name string
	// StorageClassName — the class claims must request to bind to it.
	StorageClassName string
	// Capacity — the volume's size, e.g. "100Gi". Required.
	Capacity string
	// AccessModes — ReadWriteOnce, ReadOnlyMany, ReadWriteMany or
	// ReadWriteOncePod.
	AccessModes []string
	// ReclaimPolicy — Retain (the default), Delete or Recycle: what happens to
	// the volume when its claim is released.
	ReclaimPolicy string
	// HostPath — a directory on the node. Set this or CSIDriver.
	HostPath string
	// CSIDriver — the CSI driver that provides the volume. Set this or HostPath.
	CSIDriver string
	// CSIVolumeHandle — the volume's ID within the CSI driver.
	CSIVolumeHandle string
	// Labels — set on the object.
	Labels map[string]string
	// Annotations — set on the object.
	Annotations map[string]string
}
