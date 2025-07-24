package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// DataClonePhasePending is the initial phase of a DataClone.
	DataClonePhasePending = "Pending"
	// DataClonePhaseCreatingPVC is the phase when the PVC is being created.
	DataClonePhaseCreatingPVC = "CreatingPVC"
	// DataClonePhasePodInitializing is the phase when the pod is being initialized.
	DataClonePhasePodInitializing = "PodInitializing"
	// DataClonePhaseMasking is the phase when the data is being masked.
	DataClonePhaseMasking = "MaskingInProgress"
	// DataClonePhaseReady is the phase when the clone is ready.
	DataClonePhaseReady = "Ready"
	// DataClonePhaseFailed is the phase when the clone has failed.
	DataClonePhaseFailed = "Failed"
	// DataClonePhaseDeleting is the phase when the clone is being deleted.
	DataClonePhaseDeleting = "Deleting"
)

// DataCloneSpec defines the desired state of DataClone
type DataCloneSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// SourceProfile is the name of the DataProfile to clone from.
	SourceProfile string `json:"sourceProfile"`

	// SnapshotName is the name of the specific snapshot to use.
	// If not specified, the latest snapshot will be used.
	// +optional
	SnapshotName string `json:"snapshotName,omitempty"`

	// TTL is the time-to-live for the clone. After this duration, the clone will be deleted.
	// +optional
	TTL *metav1.Duration `json:"ttl,omitempty"`

	// Database defines the database configuration for the clone.
	// +optional
	Database *DatabaseSpec `json:"database,omitempty"`

	// Pod defines the pod configuration for the clone.
	// +optional
	Pod *PodSpec `json:"pod,omitempty"`
}

// DatabaseSpec defines the database configuration for a clone.
type DatabaseSpec struct {
	// Image is the PostgreSQL image to use for the clone.
	// +optional
	Image string `json:"image,omitempty"`
	// User is the database user to create.
	// +optional
	User string `json:"user,omitempty"`
	// PasswordSecretRef is a reference to the secret containing the database password.
	// +optional
	PasswordSecretRef *corev1.SecretKeySelector `json:"passwordSecretRef,omitempty"`
	// DBName is the name of the database to create.
	// +optional
	DBName string `json:"dbname,omitempty"`
}

// PodSpec defines the pod configuration for a clone.
type PodSpec struct {
	// Resources defines the compute resources for the pod.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// DatabaseConnection defines the connection information for a database.
type DatabaseConnection struct {
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// DataCloneStatus defines the observed state of DataClone
type DataCloneStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase is the current lifecycle phase of the clone.
	// +optional
	Phase string `json:"phase,omitempty"`

	// DatabaseConnection contains the connection information for the cloned database.
	// +optional
	DatabaseConnection *DatabaseConnection `json:"databaseConnection,omitempty"`

	// Conditions represent the latest available observations of the DataClone's state.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DataClone is the Schema for the dataclones API
type DataClone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataCloneSpec   `json:"spec,omitempty"`
	Status DataCloneStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataCloneList contains a list of DataClone
type DataCloneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataClone `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataClone{}, &DataCloneList{})
}
