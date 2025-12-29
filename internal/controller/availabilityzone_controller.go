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

	halcyonv1alpha1 "github.com/halcyonproj/halcyon/api/v1alpha1"
)

// AvailabilityZoneReconciler reconciles a AvailabilityZone object
type AvailabilityZoneReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=halcyonproj.dev,resources=availabilityzones,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=availabilityzones/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=availabilityzones/finalizers,verbs=update

func (r *AvailabilityZoneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Availability zone reconcile ran")

	// fetch the availability zone resource
	az, err := r.GetAvailabilityZone(ctx, req)
	if err != nil {
		log.Error(err, "Failed to get availability zone")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// fetch the region associated with the availability zone
	region, err := r.GetRegion(ctx, az)
	if err != nil {
		log.Error(err, "Failed to get associated region")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// update availability zone topology labels
	err = r.UpdateTopologyLabels(ctx, az, region)
	if err != nil {
		log.Error(err, "Unable to update availability zone")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AvailabilityZoneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&halcyonv1alpha1.AvailabilityZone{}).
		Complete(r)
}

// GetAvailabilityZone returns the AvailabilityZone resource being reconciled
func (r *AvailabilityZoneReconciler) GetAvailabilityZone(ctx context.Context, req ctrl.Request) (*halcyonv1alpha1.AvailabilityZone, error) {
	az := &halcyonv1alpha1.AvailabilityZone{}
	err := r.Get(ctx, req.NamespacedName, az)
	return az, err
}

// GetRegion returns the Region resource associated with the availability zone being reconciled
func (r *AvailabilityZoneReconciler) GetRegion(ctx context.Context, az *halcyonv1alpha1.AvailabilityZone) (*halcyonv1alpha1.Region, error) {
	region := &halcyonv1alpha1.Region{}
	err := r.Get(ctx, types.NamespacedName{Name: az.Spec.Location.RegionName, Namespace: az.Namespace}, region)
	return region, err
}

// UpdateTopologyLabels updates the AvailabilityZone topology labels
func (r *AvailabilityZoneReconciler) UpdateTopologyLabels(ctx context.Context, az *halcyonv1alpha1.AvailabilityZone, region *halcyonv1alpha1.Region) error {
	var regionTopologyKey = "topology.halcyonproj.dev/region"

	if az.Labels == nil {
		az.Labels = make(map[string]string)
	}

	az.Labels[regionTopologyKey] = region.Name

	return r.Update(ctx, az)
}
