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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	halcyonv1alpha1 "halcyonproj.dev/halcyon/api/v1alpha1"
)

// DataCenterReconciler reconciles a DataCenter object
type DataCenterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=halcyonproj.dev,resources=datacenters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=datacenters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=datacenters/finalizers,verbs=update

func (r *DataCenterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("DataCenter reconcile ran:", "Name", req.NamespacedName)

	// fetch the datacenter resource
	datacenter, err := r.GetDataCenter(ctx, req)
	if err != nil {
		log.Error(err, "Failed to get datacenter")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// fetch the availability zone associated with the datacenter
	az, err := r.GetAvailabilityZone(ctx, datacenter)
	if err != nil {
		log.Error(err, "Failed to get associated availability zone")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// update datacenter topology labels
	err = r.UpdateTopologyLabels(ctx, datacenter, az)
	if err != nil {
		log.Error(err, "Unable to update datacenter")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DataCenterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&halcyonv1alpha1.DataCenter{}).
		Complete(r)
}

// GetDataCenter returns the DataCenter resource being reconciled
func (r *DataCenterReconciler) GetDataCenter(ctx context.Context, req ctrl.Request) (*halcyonv1alpha1.DataCenter, error) {
	datacenter := &halcyonv1alpha1.DataCenter{}
	err := r.Get(ctx, req.NamespacedName, datacenter)
	return datacenter, err
}

// GetAvailabilityZone returns the AvailabilityZone resource associated with the datacenter being reconciled
func (r *DataCenterReconciler) GetAvailabilityZone(ctx context.Context, datacenter *halcyonv1alpha1.DataCenter) (*halcyonv1alpha1.AvailabilityZone, error) {
	az := &halcyonv1alpha1.AvailabilityZone{}
	err := r.Get(ctx, types.NamespacedName{Name: datacenter.Spec.Location.AvailabilityZoneName, Namespace: datacenter.Namespace}, az)
	return az, err
}

// UpdateTopologyLabels updates the DataCenter topology labels
func (r *DataCenterReconciler) UpdateTopologyLabels(ctx context.Context, datacenter *halcyonv1alpha1.DataCenter, az *halcyonv1alpha1.AvailabilityZone) error {
	var regionTopologyKey = "topology.halcyonproj.dev/region"
	var availabilityZoneTopologyKey = "topology.halcyonproj.dev/availability-zone"

	if datacenter.Labels == nil {
		datacenter.Labels = make(map[string]string)
	}

	datacenter.Labels[availabilityZoneTopologyKey] = az.Name
	datacenter.Labels[regionTopologyKey] = az.Labels[regionTopologyKey]

	return r.Update(ctx, datacenter)
}
