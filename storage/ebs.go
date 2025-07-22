package storage

import (
	"context"
	"fmt"
	"sort"
	"time"

	phantomdbv1alpha1 "github.com/phantom-db/phantom-db/apis/v1alpha1"
	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ebsProvider is an implementation of the StorageProvider interface for AWS EBS.
type ebsProvider struct {
	client.Client
}

// NewEBSProvider creates a new EBS storage provider.
func NewEBSProvider(client client.Client) StorageProvider {
	return &ebsProvider{Client: client}
}

// CreateSnapshot creates a new VolumeSnapshot for an EBS volume.
func (p *ebsProvider) CreateSnapshot(ctx context.Context, dataProfile *phantomdbv1alpha1.DataProfile, pvcName string) (*snapshotv1.VolumeSnapshot, error) {
	snapshot := &snapshotv1.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%d", dataProfile.Name, time.Now().Unix()),
			Namespace: dataProfile.Namespace,
		},
		Spec: snapshotv1.VolumeSnapshotSpec{
			Source: snapshotv1.VolumeSnapshotSource{
				PersistentVolumeClaimName: &pvcName,
			},
			// TODO: Make VolumeSnapshotClass configurable
		},
	}

	if err := p.Create(ctx, snapshot); err != nil {
		return nil, err
	}

	return snapshot, nil
}

// GetSnapshotStatus returns the status of a VolumeSnapshot.
func (p *ebsProvider) GetSnapshotStatus(ctx context.Context, snapshot *snapshotv1.VolumeSnapshot) (string, error) {
	// Implementation to get snapshot status from AWS
	return "Ready", nil // Placeholder
}

// CleanupSnapshots deletes old snapshots based on the retention policy.
func (p *ebsProvider) CleanupSnapshots(ctx context.Context, dataProfile *phantomdbv1alpha1.DataProfile) error {
	// 1. List all snapshots for the DataProfile
	var snapshotList snapshotv1.VolumeSnapshotList
	if err := p.List(ctx, &snapshotList, client.InNamespace(dataProfile.Namespace)); err != nil {
		return err
	}

	// 2. Filter snapshots belonging to this DataProfile
	var snapshots []snapshotv1.VolumeSnapshot
	for _, snapshot := range snapshotList.Items {
		if metav1.IsControlledBy(&snapshot, dataProfile) {
			snapshots = append(snapshots, snapshot)
		}
	}

	// 3. Sort snapshots by creation time
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].CreationTimestamp.Before(&snapshots[j].CreationTimestamp)
	})

	// 4. Delete old snapshots
	if dataProfile.Spec.RetentionPolicy != nil {
		retentionCount := int(dataProfile.Spec.RetentionPolicy.Count)
		if len(snapshots) > retentionCount {
			for _, snapshotToDelete := range snapshots[retentionCount:] {
				if err := p.Delete(ctx, &snapshotToDelete); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
