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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

phantomdbv1alpha1 "github.com/vandal-db/vandal-db/apis/v1alpha1"
)

const dataCloneFinalizer = "vandal.db.io/finalizer"

// DataCloneReconciler reconciles a DataClone object
type DataCloneReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=vandal.db.io,resources=dataclones,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vandal.db.io,resources=dataclones/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vandal.db.io,resources=dataclones/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DataClone object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *DataCloneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the DataClone instance
	var dataClone phantomdbv1alpha1.DataClone
	if err := r.Get(ctx, req.NamespacedName, &dataClone); err != nil {
		log.Error(err, "unable to fetch DataClone")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the DataClone is being deleted
	if !dataClone.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&dataClone, dataCloneFinalizer) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.cleanupResources(ctx, &dataClone); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&dataClone, dataCloneFinalizer)
			if err := r.Update(ctx, &dataClone); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// The object is not being deleted, so if it does not have our finalizer,
	// then lets add the finalizer and update the object. This is equivalent
	// to registering our finalizer.
	if !controllerutil.ContainsFinalizer(&dataClone, dataCloneFinalizer) {
		controllerutil.AddFinalizer(&dataClone, dataCloneFinalizer)
		if err := r.Update(ctx, &dataClone); err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info("Reconciling DataClone", "Name", dataClone.Name)

	// 1. Create PVC from snapshot
	pvc, err := r.createPVCFromSnapshot(ctx, &dataClone)
	if err != nil {
		log.Error(err, "unable to create PVC from snapshot", "DataClone", dataClone.Name)
		return ctrl.Result{}, err
	}

	// 2. Create the database pod
	pod, err := r.createDatabasePod(ctx, &dataClone, pvc)
	if err != nil {
		log.Error(err, "unable to create database pod", "DataClone", dataClone.Name)
		return ctrl.Result{}, err
	}

	// 3. Create the connection secret
	secret, err := r.createConnectionSecret(ctx, &dataClone)
	if err != nil {
		log.Error(err, "unable to create connection secret", "DataClone", dataClone.Name)
		return ctrl.Result{}, err
	}

	// 4. Create the service
	service, err := r.createService(ctx, &dataClone)
	if err != nil {
		log.Error(err, "unable to create service", "DataClone", dataClone.Name)
		return ctrl.Result{}, err
	}

	// 5. Update status
	dataClone.Status.Phase = "Ready"
	dataClone.Status.ConnectionInfoSecret = secret.Name
	if err := r.Status().Update(ctx, &dataClone); err != nil {
		log.Error(err, "unable to update DataClone status", "DataClone", dataClone.Name)
		return ctrl.Result{}, err
	}

	// 6. Handle TTL
	if dataClone.Spec.TTL != nil {
		ttl := dataClone.Spec.TTL.Duration
		if ttl > 0 {
			return ctrl.Result{RequeueAfter: ttl}, nil
		}
	}

	return ctrl.Result{}, nil
}

func (r *DataCloneReconciler) createPVCFromSnapshot(ctx context.Context, dataClone *phantomdbv1alpha1.DataClone) (*corev1.PersistentVolumeClaim, error) {
	log := log.FromContext(ctx)

	// Define the PVC
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataClone.Name,
			Namespace: dataClone.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			DataSource: &corev1.TypedLocalObjectReference{
				APIGroup: &[]string{"snapshot.storage.k8s.io"}[0],
				Kind:     "VolumeSnapshot",
				Name:     dataClone.Spec.SnapshotName,
			},
			// TODO: Make storage class and resources configurable
		},
	}

	// Create the PVC
	if err := r.Create(ctx, pvc); err != nil {
		log.Error(err, "unable to create PVC")
		return nil, err
	}

	log.Info("Created PVC from snapshot", "PVC", pvc.Name, "Snapshot", dataClone.Spec.SnapshotName)
	return pvc, nil
}

func (r *DataCloneReconciler) createDatabasePod(ctx context.Context, dataClone *phantomdbv1alpha1.DataClone, pvc *corev1.PersistentVolumeClaim) (*corev1.Pod, error) {
	log := log.FromContext(ctx)

	// Define the Pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataClone.Name,
			Namespace: dataClone.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "postgres",
					Image: "postgres:13", // TODO: Make image configurable
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "data",
							MountPath: "/var/lib/postgresql/data",
						},
					},
					// TODO: Add env vars for password, etc.
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvc.Name,
						},
					},
				},
			},
		},
	}

	// Create the Pod
	if err := r.Create(ctx, pod); err != nil {
		log.Error(err, "unable to create pod")
		return nil, err
	}

	log.Info("Created database pod", "Pod", pod.Name)
	return pod, nil
}

func (r *DataCloneReconciler) createConnectionSecret(ctx context.Context, dataClone *phantomdbv1alpha1.DataClone) (*corev1.Secret, error) {
	log := log.FromContext(ctx)

	// Define the Secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataClone.Name,
			Namespace: dataClone.Namespace,
		},
		StringData: map[string]string{
			"host":     dataClone.Name,
			"port":     "5432",
			"user":     "postgres",
			"password": "password", // TODO: Generate random password
			"dbname":   "postgres",
		},
	}

	// Create the Secret
	if err := r.Create(ctx, secret); err != nil {
		log.Error(err, "unable to create secret")
		return nil, err
	}

	log.Info("Created connection secret", "Secret", secret.Name)
	return secret, nil
}

func (r *DataCloneReconciler) createService(ctx context.Context, dataClone *phantomdbv1alpha1.DataClone) (*corev1.Service, error) {
	log := log.FromContext(ctx)

	// Define the Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataClone.Name,
			Namespace: dataClone.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": dataClone.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Port: 5432,
				},
			},
		},
	}

	// Create the Service
	if err := r.Create(ctx, service); err != nil {
		log.Error(err, "unable to create service")
		return nil, err
	}

	log.Info("Created service", "Service", service.Name)
	return service, nil
}

func (r *DataCloneReconciler) cleanupResources(ctx context.Context, dataClone *phantomdbv1alpha1.DataClone) error {
	log := log.FromContext(ctx)

	// Delete the pod, pvc, secret, and service
	resources := []client.Object{
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: dataClone.Name, Namespace: dataClone.Namespace}},
		&corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: dataClone.Name, Namespace: dataClone.Namespace}},
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: dataClone.Name, Namespace: dataClone.Namespace}},
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: dataClone.Name, Namespace: dataClone.Namespace}},
	}

	for _, resource := range resources {
		if err := r.Delete(ctx, resource); err != nil && !client.IgnoreNotFound(err) {
			log.Error(err, "unable to delete resource", "resource", resource.GetName())
			return err
		}
	}

	log.Info("Cleaned up resources for DataClone", "Name", dataClone.Name)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DataCloneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&phantomdbv1alpha1.DataClone{}).
		Complete(r)
}
