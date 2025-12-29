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
	"fmt"
	"github.com/halcyonproj/halcyon/internal/fabric"
	"time"

	"github.com/unpoller/unifi"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	halcyonv1alpha1 "github.com/halcyonproj/halcyon/api/v1alpha1"
)

// PhysicalSwitchReconciler reconciles a PhysicalSwitch object
type PhysicalSwitchReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	FabricGetter  fabric.ClientGetter
	UnifiUsername string
	UnifiPassword string
}

// +kubebuilder:rbac:groups=halcyonproj.dev,resources=physicalswitches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=physicalports,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=physicalswitches/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=physicalports/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=physicalswitches/finalizers,verbs=update
// +kubebuilder:rbac:groups=halcyonproj.dev,resources=physicalports/finalizers,verbs=update

func (r *PhysicalSwitchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Switch reconcile ran")

	// fetch the switch resource
	sw, err := r.GetSwitch(ctx, req)
	if err != nil {
		log.Error(err, "Failed to get switch")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	rack, err := r.GetRack(ctx, sw)
	if err != nil {
		log.Error(err, "Failed to get associated rack")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	err = r.UpdateTopologyLabels(ctx, sw, rack)
	if err != nil {
		log.Error(err, "Failed to update topology labels")
	}

	// TODO: move out to FabricControllerManager
	uni, err := r.CreateUnifiClient(r.UnifiUsername, r.UnifiPassword, "https://192.168.2.1/")
	if err != nil {
		log.Error(err, "Failed to create unifi client")
	}

	sites, err := r.GetSites(uni)
	if err != nil {
		log.Error(err, "Failed to get unifi sites")
	}

	swDevice, err := r.GetSwitchDevice(uni, sites, sw.Spec.Identifier)
	if err != nil {
		log.Error(err, "Failed to get unifi switch")
	}

	err = r.CreateOrUpdatePorts(ctx, uni, sites, sw, swDevice)
	if err != nil {
		log.Error(err, "Failed to create or update ports")
	}

	err = r.UpdateSwitchStatus(ctx, sw, swDevice)
	if err != nil {
		log.Error(err, "Failed to update switch status")
	}

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PhysicalSwitchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// TODO: remove hardcoded fabric name
	fabricClassName := string(fabric.Unifi)

	// Predicate to trigger reconciliation only on unifi fabric devices
	fabricPredicate := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			newObj := e.ObjectNew.(*halcyonv1alpha1.PhysicalSwitch)
			return newObj.Spec.FabricClass == fabricClassName
		},

		// Allow create events
		CreateFunc: func(e event.CreateEvent) bool {
			obj := e.Object.(*halcyonv1alpha1.PhysicalSwitch)
			return obj.Spec.FabricClass == fabricClassName
		},

		// Allow delete events
		DeleteFunc: func(e event.DeleteEvent) bool {
			obj := e.Object.(*halcyonv1alpha1.PhysicalSwitch)
			return obj.Spec.FabricClass == fabricClassName
		},

		// Allow generic events (e.g., external triggers)
		GenericFunc: func(e event.GenericEvent) bool {
			obj := e.Object.(*halcyonv1alpha1.PhysicalSwitch)
			return obj.Spec.FabricClass == fabricClassName
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&halcyonv1alpha1.PhysicalSwitch{}).
		WithEventFilter(fabricPredicate).
		Complete(r)
}

// GetSwitch returns the PhysicalSwitch resource being reconciled
func (r *PhysicalSwitchReconciler) GetSwitch(ctx context.Context, req ctrl.Request) (*halcyonv1alpha1.PhysicalSwitch, error) {
	sw := &halcyonv1alpha1.PhysicalSwitch{}
	err := r.Get(ctx, req.NamespacedName, sw)
	return sw, err
}

// GetRack returns the Rack resource being reconciled
func (r *PhysicalSwitchReconciler) GetRack(ctx context.Context, sw *halcyonv1alpha1.PhysicalSwitch) (*halcyonv1alpha1.Rack, error) {
	rack := &halcyonv1alpha1.Rack{}
	err := r.Get(ctx, types.NamespacedName{Name: sw.Spec.Location.RackName, Namespace: sw.Namespace}, rack)
	return rack, err
}

// GetPort returns the PhysicalPort resource being reconciled
func (r *PhysicalSwitchReconciler) GetPort(ctx context.Context, name string, namespace string) (*halcyonv1alpha1.PhysicalPort, error) {
	port := &halcyonv1alpha1.PhysicalPort{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, port)
	return port, err
}

// UpdateTopologyLabels updates the PhysicalSwitch topology labels
func (r *PhysicalSwitchReconciler) UpdateTopologyLabels(ctx context.Context, sw *halcyonv1alpha1.PhysicalSwitch, rack *halcyonv1alpha1.Rack) error {
	var regionTopologyKey = "topology.halcyonproj.dev/region"
	var availabilityZoneTopologyKey = "topology.halcyonproj.dev/availability-zone"
	var datacenterTopologyKey = "topology.halcyonproj.dev/datacenter"
	var rackTopologyKey = "topology.halcyonproj.dev/rack"

	if sw.Labels == nil {
		sw.Labels = make(map[string]string)
	}

	sw.Labels[rackTopologyKey] = rack.Name
	sw.Labels[datacenterTopologyKey] = rack.Labels[datacenterTopologyKey]
	sw.Labels[availabilityZoneTopologyKey] = rack.Labels[availabilityZoneTopologyKey]
	sw.Labels[regionTopologyKey] = rack.Labels[regionTopologyKey]

	return r.Update(ctx, sw)
}

func (r *PhysicalSwitchReconciler) CreateUnifiClient(username string, password string, endpoint string) (*fabric.UnifiClient, error) {
	fabricClient, err := r.FabricGetter.GetClient(&fabric.ClientConfig{
		Fabric:   fabric.Unifi,
		Username: r.UnifiUsername,
		Password: r.UnifiPassword,
		Url:      endpoint,
	})
	if err != nil {
		return nil, err
	}

	uni, _ := fabricClient.(*fabric.UnifiClient)
	return uni, nil
}

func (r *PhysicalSwitchReconciler) UpdateSwitchStatus(ctx context.Context, sw *halcyonv1alpha1.PhysicalSwitch, swDevice *unifi.USW) error {
	sw.Status.IPAddress = swDevice.IP
	sw.Status.MACAddress = swDevice.Mac
	sw.Status.DeviceID = swDevice.DeviceID
	sw.Status.Model = swDevice.Model
	sw.Status.Serial = swDevice.Serial
	sw.Status.Version = swDevice.Version

	// TODO: to be set dynamically in future when abstracted to Fabric Controller Manager
	sw.Status.Manufacturer = sw.Spec.FabricClass

	return r.Status().Update(ctx, sw)
}

func (r *PhysicalSwitchReconciler) GetClientMACAddresses(uni *fabric.UnifiClient, sites []*unifi.Site, swMac string) (map[string]string, error) {
	clientMACs := make(map[string]string)
	clients, err := uni.GetClients(sites)

	if err != nil {
		return clientMACs, err
	}

	for _, client := range clients {
		if client.SwPort.Val != 0 && client.SwMac == swMac {
			clientMACs[client.SwPort.Txt] = client.Mac
		}
	}

	return clientMACs, err
}

func (r *PhysicalSwitchReconciler) GetSites(uni *fabric.UnifiClient) ([]*unifi.Site, error) {
	sites, err := uni.GetSites()
	if err != nil {
		return nil, err
	}

	return sites, err
}

func (r *PhysicalSwitchReconciler) GetSwitchDevice(uni *fabric.UnifiClient, sites []*unifi.Site, deviceID string) (*unifi.USW, error) {
	devices, err := uni.GetDevices(sites)
	if err != nil {
		return &unifi.USW{}, err
	}

	var swDevice *unifi.USW

	for _, device := range devices.USWs {
		if device.Serial == deviceID {
			swDevice = device
			break
		}
	}

	if swDevice.Serial != "" {
		return swDevice, nil
	} else {
		return &unifi.USW{}, fmt.Errorf("Could not find switch with associated device ID: %s", deviceID)
	}
}

func (r *PhysicalSwitchReconciler) CreateOrUpdatePorts(ctx context.Context, uni *fabric.UnifiClient, sites []*unifi.Site, sw *halcyonv1alpha1.PhysicalSwitch, swDevice *unifi.USW) error {
	// TODO remove ports that are no longer found
	clientMACs, err := r.GetClientMACAddresses(uni, sites, swDevice.Mac)
	if err != nil {
		return err
	}

	for _, port := range swDevice.PortTable {
		portName := sw.Name + "-" + port.PortIdx.Txt

		physicalPort, err := r.GetPort(ctx, portName, sw.Namespace)
		if err != nil && apierrors.IsNotFound(err) {
			err := r.CreatePort(ctx, sw, portName, sw.Namespace)
			if err != nil {
				// TODO handle error
			}
		} else if err != nil {
			// TODO handle error
		}

		// TODO error handling
		err = r.UpdatePortTopologyLabels(ctx, physicalPort, sw)
		err = r.UpdatePortStatus(ctx, physicalPort, port, clientMACs)
	}

	return nil
}

func (r *PhysicalSwitchReconciler) UpdatePortTopologyLabels(ctx context.Context, physicalPort *halcyonv1alpha1.PhysicalPort, sw *halcyonv1alpha1.PhysicalSwitch) error {
	var regionTopologyKey = "topology.halcyonproj.dev/region"
	var availabilityZoneTopologyKey = "topology.halcyonproj.dev/availability-zone"
	var datacenterTopologyKey = "topology.halcyonproj.dev/datacenter"
	var rackTopologyKey = "topology.halcyonproj.dev/rack"
	var switchTopologyKey = "topology.halcyonproj.dev/switch"

	if physicalPort.Labels == nil {
		physicalPort.Labels = make(map[string]string)
	}

	physicalPort.Labels[switchTopologyKey] = sw.Name
	physicalPort.Labels[rackTopologyKey] = sw.Labels[rackTopologyKey]
	physicalPort.Labels[datacenterTopologyKey] = sw.Labels[datacenterTopologyKey]
	physicalPort.Labels[availabilityZoneTopologyKey] = sw.Labels[availabilityZoneTopologyKey]
	physicalPort.Labels[regionTopologyKey] = sw.Labels[regionTopologyKey]

	return r.Update(ctx, physicalPort)
}

func (r *PhysicalSwitchReconciler) CreatePort(ctx context.Context, sw *halcyonv1alpha1.PhysicalSwitch, name string, namespace string) error {
	// Define a new PhysicalPort
	newPhysicalPort := &halcyonv1alpha1.PhysicalPort{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	controllerutil.SetControllerReference(sw, newPhysicalPort, r.Scheme)

	return r.Create(ctx, newPhysicalPort)
}

func (r *PhysicalSwitchReconciler) UpdatePortStatus(ctx context.Context, physicalPort *halcyonv1alpha1.PhysicalPort, port unifi.Port, clientMACs map[string]string) error {
	if port.Enable.Val == false {
		physicalPort.Status.State = "disabled"
	} else if port.Up.Val == true {
		physicalPort.Status.State = "connected"
	} else {
		physicalPort.Status.State = "disconnected"
	}

	if port.FullDuplex.Val == true {
		physicalPort.Status.Duplex = "full"
	} else {
		physicalPort.Status.Duplex = "half"
	}

	if value, ok := clientMACs[port.PortIdx.Txt]; ok {
		physicalPort.Status.ObservedMAC = value
	} else {
		physicalPort.Status.ObservedMAC = ""
	}

	if port.OpMode == "aggregate" || port.AggregatedBy.Txt != "false" {
		physicalPort.Status.AggregationProtocol = "LACP"
		physicalPort.Status.AggregationEnabled = true
	} else {
		physicalPort.Status.AggregationProtocol = ""
		physicalPort.Status.AggregationEnabled = false
	}

	physicalPort.Status.IsUplink = port.IsUplink.Val == true
	physicalPort.Status.MaxSpeed = int(port.SpeedCaps.Val)
	physicalPort.Status.NegotiatedSpeed = int(port.Speed.Val)
	physicalPort.Status.AutoNegotiated = port.Autoneg.Val || false

	return r.Status().Update(ctx, physicalPort)
}
