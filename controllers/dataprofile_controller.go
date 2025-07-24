/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vandalv1alpha1 "github.com/Oridak771/Vandal/apis/v1alpha1"
	"github.com/Oridak771/Vandal/storage"
)

// DataProfileReconciler reconciles a DataProfile object
type DataProfileReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	Cron            *cron.Cron
	StorageProvider storage.StorageProvider
}

//+kubebuilder:rbac:groups=vandal.db.io,resources=dataprofiles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vandal.db.io,resources=dataprofiles/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vandal.db.io,resources=dataprofiles/finalizers,verbs=update
//+kubebuilder:rbac:groups=snapshot.storage.k8s.io,resources=volumesnapshots,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=snapshot.storage.k8s.io,resources=volumesnapshotclasses,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DataProfile object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *DataProfileReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the DataProfile instance
	var dataProfile vandalv1alpha1.DataProfile
	if err := r.Get(ctx, req.NamespacedName, &dataProfile); err != nil {
		log.Error(err, "unable to fetch DataProfile")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling DataProfile", "Name", dataProfile.Name)

	// Handle cron job scheduling
	if dataProfile.Spec.Schedule != "" {
		entryID := cron.EntryID(dataProfile.UID)
		entry := r.Cron.Entry(entryID)

		if entry.ID == 0 {
			// No entry found, create a new one
			id, err := r.Cron.AddFunc(dataProfile.Spec.Schedule, func() {
				log.Info("Triggering snapshot for DataProfile", "Name", dataProfile.Name)
				if err := r.createVolumeSnapshot(ctx, &dataProfile); err != nil {
					log.Error(err, "failed to create volume snapshot", "DataProfile", dataProfile.Name)
				}
			})
			if err != nil {
				log.Error(err, "unable to add cron job", "DataProfile", dataProfile.Name)
				return ctrl.Result{}, err
			}
			log.Info("Added new cron job", "EntryID", id)
		} else {
			// Entry exists, check if schedule needs update
			if entry.Schedule.String() != dataProfile.Spec.Schedule {
				r.Cron.Remove(entryID)
				id, err := r.Cron.AddFunc(dataProfile.Spec.Schedule, func() {
					log.Info("Triggering snapshot for DataProfile", "Name", dataProfile.Name)
				if err := r.createVolumeSnapshot(ctx, &dataProfile); err != nil {
					log.Error(err, "failed to create volume snapshot", "DataProfile", dataProfile.Name)
				}
				})
				if err != nil {
					log.Error(err, "unable to update cron job", "DataProfile", dataProfile.Name)
					return ctrl.Result{}, err
				}
				log.Info("Updated cron job", "EntryID", id)
			}
		}
	}

	// Cleanup old snapshots
	if err := r.StorageProvider.CleanupSnapshots(ctx, &dataProfile); err != nil {
		log.Error(err, "unable to cleanup snapshots", "DataProfile", dataProfile.Name)
		// Update status with failure condition
		meta.SetStatusCondition(&dataProfile.Status.Conditions, metav1.Condition{
			Type:    "SnapshotCleanup",
			Status:  metav1.ConditionFalse,
			Reason:  "Error",
			Message: err.Error(),
		})
		if err := r.Status().Update(ctx, &dataProfile); err != nil {
			log.Error(err, "unable to update DataProfile status", "DataProfile", dataProfile.Name)
		}
		return ctrl.Result{}, err
	}

	// Update status with success condition
	meta.SetStatusCondition(&dataProfile.Status.Conditions, metav1.Condition{
		Type:   "SnapshotCleanup",
		Status: metav1.ConditionTrue,
		Reason: "Success",
	})
	dataProfile.Status.LastSnapshotTime = &metav1.Time{Time: time.Now()}
	if err := r.Status().Update(ctx, &dataProfile); err != nil {
		log.Error(err, "unable to update DataProfile status", "DataProfile", dataProfile.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *DataProfileReconciler) createVolumeSnapshot(ctx context.Context, dataProfile *vandalv1alpha1.DataProfile) error {
	log := log.FromContext(ctx)

	// 1. Set the phase to CreatingSnapshot
	dataProfile.Status.Phase = vandalv1alpha1.DataProfilePhaseCreatingSnapshot
	if err := r.Status().Update(ctx, dataProfile); err != nil {
		log.Error(err, "unable to update DataProfile status")
		return err
	}

	// 2. Get the PVC name from the spec
	pvcName := dataProfile.Spec.Target.PVCName

	// 3. Create the VolumeSnapshot object
	snapshot, err := r.StorageProvider.CreateSnapshot(ctx, dataProfile, pvcName)
	if err != nil {
		log.Error(err, "unable to create VolumeSnapshot")
		dataProfile.Status.Phase = vandalv1alpha1.DataProfilePhaseFailed
		if err := r.Status().Update(ctx, dataProfile); err != nil {
			log.Error(err, "unable to update DataProfile status")
		}
		return err
	}

	// 4. Set the phase to SnapshotReady
	dataProfile.Status.Phase = vandalv1alpha1.DataProfilePhaseSnapshotReady
	if err := r.Status().Update(ctx, dataProfile); err != nil {
		log.Error(err, "unable to update DataProfile status")
		return err
	}

	log.Info("Created VolumeSnapshot", "Name", snapshot.Name)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DataProfileReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vandalv1alpha1.DataProfile{}).
		Complete(r)
}
