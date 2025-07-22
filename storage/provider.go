package storage

import (
	"context"

	phantomdbv1alpha1 "github.com/phantom-db/phantom-db/apis/v1alpha1"
	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
)

// StorageProvider defines the interface for a storage provider.
type StorageProvider interface {
	// CreateSnapshot creates a snapshot of a volume.
	CreateSnapshot(ctx context.Context, dataProfile *phantomdbv1alpha1.DataProfile, pvcName string) (*snapshotv1.VolumeSnapshot, error)
	// GetSnapshotStatus returns the status of a snapshot.
	GetSnapshotStatus(ctx context.Context, snapshot *snapshotv1.VolumeSnapshot) (string, error)
	// CleanupSnapshots deletes old snapshots based on the retention policy.
	CleanupSnapshots(ctx context.Context, dataProfile *phantomdbv1alpha1.DataProfile) error
}
