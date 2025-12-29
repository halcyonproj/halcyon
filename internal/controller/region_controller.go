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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	halcyonv1alpha1 "halcyonproj.dev/halcyon/api/v1alpha1"
)

// RegionReconciler reconciles a Region object
type RegionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=halcyonproj.dev,resources=regions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=regions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=regions/finalizers,verbs=update

func (r *RegionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Region reconcile ran")

	// fetch the region resource
	_, err := r.GetRegion(ctx, req)
	if err != nil {
		log.Error(err, "Failed to get region")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RegionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&halcyonv1alpha1.Region{}).
		Complete(r)
}

// GetRegion returns the Region resource being reconciled
func (r *RegionReconciler) GetRegion(ctx context.Context, req ctrl.Request) (*halcyonv1alpha1.Region, error) {
	region := &halcyonv1alpha1.Region{}
	err := r.Get(ctx, req.NamespacedName, region)
	return region, err
}
