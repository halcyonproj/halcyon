package fabric

import (
	"fmt"
	"github.com/unpoller/unifi"
)

type UnifiClient struct {
	// hack: leaving as anonymous struct field so
	// clients can leverage {UnifiClient}.(*unifi.Unifi)
	// to reach into internal APIs if needed.
	//
	// For example:
	//
	// fabricClient, _ := r.FabricGetter.GetClient(&fabric.ClientConfig{...})
	// uni, _ := fabricClient.(*fabric.UnifiClient)
	// sites, err := uni.Unifi.GetSites()
	*unifi.Unifi
}

func (c *UnifiClient) GetFabricDevices() ([]*Device, error) {
	sites, err := c.GetSites()
	if err != nil {
		return nil, fmt.Errorf("getting unifi sites: %w", err)
	}

	uniDevices, err := c.GetDevices(sites)
	if err != nil {
		return nil, fmt.Errorf("getting unifi clients: %w", err)
	}

	var devices []*Device

	for _, dvc := range uniDevices.UAPs {
		devices = append(devices, &Device{
			MAC: dvc.Mac,
			IP:  dvc.IP,
		})
	}
	for _, dvc := range uniDevices.USGs {
		devices = append(devices, &Device{
			MAC: dvc.Mac,
			IP:  dvc.IP,
		})
	}
	for _, dvc := range uniDevices.USWs {
		devices = append(devices, &Device{
			MAC: dvc.Mac,
			IP:  dvc.IP,
		})
	}
	for _, dvc := range uniDevices.UDMs {
		devices = append(devices, &Device{
			MAC: dvc.Mac,
			IP:  dvc.IP,
		})
	}
	for _, dvc := range uniDevices.UXGs {
		devices = append(devices, &Device{
			MAC: dvc.Mac,
			IP:  dvc.IP,
		})
	}
	for _, dvc := range uniDevices.PDUs {
		devices = append(devices, &Device{
			MAC: dvc.Mac,
			IP:  dvc.IP,
		})
	}

	return devices, nil
}
