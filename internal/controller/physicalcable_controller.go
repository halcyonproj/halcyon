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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	halcyonv1alpha1 "halcyonproj.dev/halcyon/api/v1alpha1"
	"halcyonproj.dev/halcyon/internal/fabric"
)

// PhysicalCableReconciler reconciles a PhysicalCable object
type PhysicalCableReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	FabricGetter  fabric.ClientGetter
	UnifiUsername string
	UnifiPassword string
}

// +kubebuilder:rbac:groups=halcyonproj.dev,resources=physicalcables,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=physicalcables/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=physicalcables/finalizers,verbs=update

func (r *PhysicalCableReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// TODO: move out to FabricControllerManager
	log := log.FromContext(ctx)

	log.Info("Cable reconcile ran")

	// fetch the cable resource
	cable, err := r.GetCable(ctx, req)
	if err != nil {
		log.Error(err, "Failed to get cable")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// fetch associated resources
	port, err := r.GetPort(ctx, cable)
	iface, err := r.GetInterface(ctx, cable)

	// update cable topology labels
	err = r.UpdateTopologyLabels(ctx, cable, port, iface)
	if err != nil {
		log.Error(err, "Unable to update cable")
		return ctrl.Result{}, err
	}

	// derive and update the state of the cable
	cableState := r.DeriveCableState(port, iface)
	err = r.UpdateCableState(ctx, cable, cableState)
	if err != nil {
		log.Error(err, "Unable to update cable status")
		return ctrl.Result{}, err
	}

	// update interface status with IP address
	err = r.UpdateInterfaceStatus(ctx, iface)
	if err != nil {
		log.Error(err, "Failed to update PhysicalInterface Status", "PhysicalInterface.Namespace", iface.Namespace, "PhysicalInterface.Name", iface.Name)
	}

	// TODO delete the IP if the cable is "unplugged"

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PhysicalCableReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&halcyonv1alpha1.PhysicalCable{}).
		Complete(r)
}

// GetCable returns the PhysicalCable resource being reconciled
func (r *PhysicalCableReconciler) GetCable(ctx context.Context, req ctrl.Request) (*halcyonv1alpha1.PhysicalCable, error) {
	cable := &halcyonv1alpha1.PhysicalCable{}
	err := r.Get(ctx, req.NamespacedName, cable)
	return cable, err
}

// GetPort returns the PhysicalPort resource associated with the cable being reconciled
func (r *PhysicalCableReconciler) GetPort(ctx context.Context, cable *halcyonv1alpha1.PhysicalCable) (*halcyonv1alpha1.PhysicalPort, error) {
	port := &halcyonv1alpha1.PhysicalPort{}
	err := r.Get(ctx, types.NamespacedName{Name: cable.Spec.TargetPort, Namespace: cable.Namespace}, port)
	// log.Error(err, "unable to fetch port")
	//cableState = "disconnected"
	//if port.Status.State != "connected" {
	//		log.Error(nil, "Port status is not 'connected'")
	//		cableState = "disconnected"
	//	}
	return port, err
}

// GetInterface returns the PhysicalInterface resource associated with the cable being reconciled
func (r *PhysicalCableReconciler) GetInterface(ctx context.Context, cable *halcyonv1alpha1.PhysicalCable) (*halcyonv1alpha1.PhysicalInterface, error) {
	iface := &halcyonv1alpha1.PhysicalInterface{}
	err := r.Get(ctx, types.NamespacedName{Name: cable.Spec.SourceInterface, Namespace: cable.Namespace}, iface)
	// log.Error(err, "unable to fetch interface")
	// cableState = "disconnected"

	return iface, err
}

// UpdateCableState updates the PhysicalCable resource Status object to indicate whether or not the cable is connected
func (r *PhysicalCableReconciler) UpdateCableState(ctx context.Context, cable *halcyonv1alpha1.PhysicalCable, cableState string) error {
	cable.Status = halcyonv1alpha1.PhysicalCableStatus{
		State: cableState,
	}

	return r.Status().Update(ctx, cable)
}

// DeriveCableState determines whether the cable is connected or disconnected
func (r *PhysicalCableReconciler) DeriveCableState(port *halcyonv1alpha1.PhysicalPort, iface *halcyonv1alpha1.PhysicalInterface) string {
	cableState := "connected"

	// TODO ensure port and iface even exist first

	if port.Status.ObservedMAC != iface.Status.MACAddress {
		cableState = "disconnected"
	}

	return cableState
}

// UpdateTopologyLabels updates the PhysicalCable topology labels
func (r *PhysicalCableReconciler) UpdateTopologyLabels(ctx context.Context, cable *halcyonv1alpha1.PhysicalCable, port *halcyonv1alpha1.PhysicalPort, iface *halcyonv1alpha1.PhysicalInterface) error {
	var regionTopologyKey = "topology.halcyonproj.dev/region"
	var availabilityZoneTopologyKey = "topology.halcyonproj.dev/availability-zone"
	var datacenterTopologyKey = "topology.halcyonproj.dev/datacenter"
	var switchTopologyKey = "topology.halcyonproj.dev/switch"
	var computeTopologyKey = "topology.halcyonproj.dev/compute"

	if cable.Labels == nil {
		cable.Labels = make(map[string]string)
	}

	cable.Labels[computeTopologyKey] = iface.Labels[computeTopologyKey]
	cable.Labels[switchTopologyKey] = port.Labels[switchTopologyKey]
	cable.Labels[datacenterTopologyKey] = port.Labels[datacenterTopologyKey]
	cable.Labels[availabilityZoneTopologyKey] = port.Labels[availabilityZoneTopologyKey]
	cable.Labels[regionTopologyKey] = port.Labels[regionTopologyKey]

	return r.Update(ctx, cable)
}

// UpdateInterfaceStatus updates the associated PhysicalInterface Status object to include the IP address of the interface, if there is a connected cable
func (r *PhysicalCableReconciler) UpdateInterfaceStatus(ctx context.Context, iface *halcyonv1alpha1.PhysicalInterface) error {
	// TODO: get credentails from provider configuration / secret
	fabricClient, err := r.FabricGetter.GetClient(&fabric.ClientConfig{
		Fabric:   fabric.Unifi,
		Username: r.UnifiUsername,
		Password: r.UnifiPassword,
		Url:      "https://192.168.2.1/",
	})
	if err != nil {
		return err
	}

	/*
		 * Need to test; seems like Port type from switch did not contain
		 * MAC or IP in response. Leaving commented out for now. The above
		 * Updates still leverage a cached version of the client.

			devices, err := fabricClient.GetFabricDevices()
			clientIP := ""
			for _, device := range devices {
				if device.MAC == iface.Status.MACAddress {
					clientIP = device.IP
				}
			}
	*/

	uni, _ := fabricClient.(*fabric.UnifiClient)
	sites, err := uni.GetSites()
	if err != nil {
		return err
	}
	clients, err := uni.GetClients(sites)
	if err != nil {
		return err
	}

	clientIP := ""
	for _, client := range clients {
		if client.Mac == iface.Status.MACAddress {
			clientIP = client.IP
		}
	}

	iface.Status.IPAddress = clientIP
	return r.Status().Update(ctx, iface)
}
