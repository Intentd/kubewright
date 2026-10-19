package types

// VolumeSource is one pod volume. Set exactly one source.
type VolumeSource struct {
	// Name — the volume's name, which a VolumeMount refers to.
	Name string `yaml:"name" json:"name"`

	// ConfigMap — a ConfigMap's keys as files.
	ConfigMap *ConfigMapVolumeSource `yaml:"configMap,omitempty" json:"configMap,omitempty"`

	// Secret — a Secret's keys as files.
	Secret *SecretVolumeSource `yaml:"secret,omitempty" json:"secret,omitempty"`

	// EmptyDir — an empty directory that lives as long as the pod.
	EmptyDir *EmptyDirVolumeSource `yaml:"emptyDir,omitempty" json:"emptyDir,omitempty"`

	// PersistentVolumeClaim — an existing claim.
	PersistentVolumeClaim *PVCVolumeSource `yaml:"persistentVolumeClaim,omitempty" json:"persistentVolumeClaim,omitempty"`

	// HostPath — a file or directory on the node.
	HostPath *HostPathVolumeSource `yaml:"hostPath,omitempty" json:"hostPath,omitempty"`
}

// ConfigMapVolumeSource is a ConfigMap mounted as a volume.
type ConfigMapVolumeSource struct {
	// Name — the ConfigMap's name.
	Name string `yaml:"name" json:"name"`
}

// SecretVolumeSource is a Secret mounted as a volume.
type SecretVolumeSource struct {
	// Name — the Secret's name.
	Name string `yaml:"name" json:"name"`
}

// EmptyDirVolumeSource is an empty directory that lives as long as the pod.
type EmptyDirVolumeSource struct{}

// PVCVolumeSource is an existing PersistentVolumeClaim mounted as a volume.
type PVCVolumeSource struct {
	// ClaimName — the claim's name.
	ClaimName string `yaml:"claimName" json:"claimName"`

	// ReadOnly — mount it read-only.
	ReadOnly bool `yaml:"readOnly,omitempty" json:"readOnly,omitempty"`
}

// HostPathVolumeSource is a path on the node mounted as a volume.
type HostPathVolumeSource struct {
	// Path — the path on the node.
	Path string `yaml:"path" json:"path"`

	// Type — what must exist at Path: DirectoryOrCreate, Directory,
	// FileOrCreate, File, Socket, CharDevice or BlockDevice. Empty checks nothing.
	Type string `yaml:"type,omitempty" json:"type,omitempty"`
}

// VolumeMount mounts a pod volume into a container.
type VolumeMount struct {
	// Name — the volume's name.
	Name string `yaml:"name" json:"name"`

	// MountPath — where the volume appears in the container.
	MountPath string `yaml:"mountPath" json:"mountPath"`

	// SubPath — mount only this path within the volume. Empty mounts its root.
	SubPath string `yaml:"subPath,omitempty" json:"subPath,omitempty"`

	// ReadOnly — mount it read-only.
	ReadOnly bool `yaml:"readOnly,omitempty" json:"readOnly,omitempty"`
}
