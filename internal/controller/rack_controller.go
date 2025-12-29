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

// RackReconciler reconciles a Rack object
type RackReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=halcyonproj.dev,resources=racks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=racks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=racks/finalizers,verbs=update

func (r *RackReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Rack reconcile ran")

	// fetch the rack resource
	rack, err := r.GetRack(ctx, req)
	if err != nil {
		log.Error(err, "Failed to get rack")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// fetch the datacenter associated with the rack
	datacenter, err := r.GetDataCenter(ctx, rack)
	if err != nil {
		log.Error(err, "Failed to get associated datacenter")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// TODO validate rack type

	// update rack topology labels
	err = r.UpdateTopologyLabels(ctx, rack, datacenter)
	if err != nil {
		log.Error(err, "Unable to update rack:", req.NamespacedName)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&halcyonv1alpha1.Rack{}).
		Complete(r)
}

// GetRack returns the Rack resource being reconciled
func (r *RackReconciler) GetRack(ctx context.Context, req ctrl.Request) (*halcyonv1alpha1.Rack, error) {
	rack := &halcyonv1alpha1.Rack{}
	err := r.Get(ctx, req.NamespacedName, rack)
	return rack, err
}

// GetDataCenter returns the DataCenter resource associated with the rack being reconciled
func (r *RackReconciler) GetDataCenter(ctx context.Context, rack *halcyonv1alpha1.Rack) (*halcyonv1alpha1.DataCenter, error) {
	datacenter := &halcyonv1alpha1.DataCenter{}
	err := r.Get(ctx, types.NamespacedName{Name: rack.Spec.Location.DataCenterName, Namespace: rack.Namespace}, datacenter)
	return datacenter, err
}

// UpdateTopologyLabels updates the Rack topology labels
func (r *RackReconciler) UpdateTopologyLabels(ctx context.Context, rack *halcyonv1alpha1.Rack, datacenter *halcyonv1alpha1.DataCenter) error {
	var regionTopologyKey = "topology.halcyonproj.dev/region"
	var availabilityZoneTopologyKey = "topology.halcyonproj.dev/availability-zone"
	var datacenterTopologyKey = "topology.halcyonproj.dev/datacenter"

	if rack.Labels == nil {
		rack.Labels = make(map[string]string)
	}

	rack.Labels[datacenterTopologyKey] = datacenter.Name
	rack.Labels[availabilityZoneTopologyKey] = datacenter.Labels[availabilityZoneTopologyKey]
	rack.Labels[regionTopologyKey] = datacenter.Labels[regionTopologyKey]

	return r.Update(ctx, rack)
}
