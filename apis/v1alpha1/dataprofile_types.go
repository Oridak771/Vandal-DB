package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// DataProfilePhasePending is the initial phase of a DataProfile.
	DataProfilePhasePending = "Pending"
	// DataProfilePhaseCreatingSnapshot is the phase when a snapshot is being created.
	DataProfilePhaseCreatingSnapshot = "CreatingSnapshot"
	// DataProfilePhaseSnapshotReady is the phase when the snapshot is ready.
	DataProfilePhaseSnapshotReady = "SnapshotReady"
	// DataProfilePhaseFailed is the phase when the snapshot has failed.
	DataProfilePhaseFailed = "Failed"
)

// DataProfileSpec defines the desired state of DataProfile
type DataProfileSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Schedule for automated snapshots, in cron format.
	// +optional
	Schedule string `json:"schedule,omitempty"`

	// RetentionPolicy defines how many snapshots to keep.
	// +optional
	RetentionPolicy *RetentionPolicy `json:"retentionPolicy,omitempty"`

	// Target defines the database to be profiled.
	Target DatabaseTarget `json:"target"`

	// Masking defines the data masking rules.
	// +optional
	Masking MaskingSpec `json:"masking,omitempty"`
}

// DatabaseTarget defines the database connection information.
type DatabaseTarget struct {
	// SecretName is the name of the secret containing the database credentials.
	SecretName string `json:"secretName"`
	// PVCName is the name of the PersistentVolumeClaim to be snapshotted.
	PVCName string `json:"pvcName"`
}

// RetentionPolicy defines the policy for retaining snapshots.
type RetentionPolicy struct {
	// Number of snapshots to keep.
	// +optional
	Count int32 `json:"count,omitempty"`
}

// MaskingSpec defines the data masking configuration.
type MaskingSpec struct {
	// Rules is a list of masking rules to apply.
	// +optional
	Rules []MaskingRule `json:"rules,omitempty"`
}

// MaskingRule defines a single data masking rule.
type MaskingRule struct {
	// Table to apply the rule to.
	Table string `json:"table"`
	// Column to apply the rule to.
	Column string `json:"column"`
	// Transformation to apply.
	Transformation string `json:"transformation"`
}

// DataProfileStatus defines the observed state of DataProfile
type DataProfileStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase is the current lifecycle phase of the data profile.
	// +optional
	Phase string `json:"phase,omitempty"`

	// LastSnapshotTime is the time the last snapshot was taken.
	// +optional
	LastSnapshotTime *metav1.Time `json:"lastSnapshotTime,omitempty"`

	// Conditions represent the latest available observations of the DataProfile's state.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DataProfile is the Schema for the dataprofiles API
type DataProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataProfileSpec   `json:"spec,omitempty"`
	Status DataProfileStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataProfileList contains a list of DataProfile
type DataProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataProfile `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataProfile{}, &DataProfileList{})
}
