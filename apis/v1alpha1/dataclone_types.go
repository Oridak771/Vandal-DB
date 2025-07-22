package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
}

// DataCloneStatus defines the observed state of DataClone
type DataCloneStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase is the current lifecycle phase of the clone.
	// +optional
	Phase string `json:"phase,omitempty"`

	// ConnectionInfoSecret is the name of the secret containing the connection information for the cloned database.
	// +optional
	ConnectionInfoSecret string `json:"connectionInfoSecret,omitempty"`

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
