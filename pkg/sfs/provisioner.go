/*
Copyright 2018 The Kubernetes Authors.

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

package sfs

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/huaweicloud/external-sfs/pkg/config"
	"github.com/huaweicloud/external-sfs/pkg/sfs/backends"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// Defines constants
const (
	ShareStatusAvailable = "available"
)

// Provisioner implements controller.Provisioner interface
type Provisioner struct {
	clientset    clientset.Interface
	cloudconfig  config.CloudCredentials
	sharetimeout int
}

// NewProvisioner creates a new instance of sfs provisioner
func NewProvisioner(c clientset.Interface, cc config.CloudCredentials, timeout int) *Provisioner {

	// init backends for provisioner
	InitBackends()

	// return provisioner instance
	return &Provisioner{
		clientset:    c,
		cloudconfig:  cc,
		sharetimeout: timeout,
	}
}

// Provision a share in sfs
func (p *Provisioner) Provision(volOptions controller.VolumeOptions) (*v1.PersistentVolume, error) {

	// selector check
	if volOptions.PVC.Spec.Selector != nil {
		return nil, fmt.Errorf("claim Selector is not supported")
	}

	// init sfs client
	client, err := p.cloudconfig.SFSV2Client()
	if err != nil {
		return nil, fmt.Errorf("failed to create SFS v2 client: %v", err)
	}

	// create share
	share, err := CreateShare(client, &volOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create share: %v", err)
	}

	// wait fo share available
	err = WaitForShareStatus(client, share.ID, ShareStatusAvailable, p.sharetimeout)
	if err != nil {
		return nil, fmt.Errorf("waiting for share %s to become created failed: %v", share.ID, err)
	}

	// get new share
	share, err = GetShare(client, share.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get share: %v", err)
	}

	// get location
	location := share.ExportLocation
	if (len(location) == 0) && (len(share.ExportLocations) > 0) {
		location = share.ExportLocations[0]
	}
	if len(location) == 0 {
		return nil, fmt.Errorf("failed to get share %s location", share.ID)
	} else {
		glog.Infof("get share %s location: %s", share.ID, location)
	}

	// get backend
	b, err := GetBackend(share.ShareProto)
	if err != nil {
		return nil, fmt.Errorf("failed to get backend: %v", err)
	}

	// get persistent volume source
	pvsource, err := b.BuildSource(&backends.BuildSourceArgs{Location: location})
	if err != nil {
		return nil, fmt.Errorf("failed to build source from backend: %v", err)
	}

	return &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:        volOptions.PVName,
			Annotations: map[string]string{},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: volOptions.PersistentVolumeReclaimPolicy,
			AccessModes:                   volOptions.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): volOptions.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: *pvsource,
		},
	}, nil
}

// Delete a share from sfs
func (p *Provisioner) Delete(pv *v1.PersistentVolume) error {

	// init sfs client
	client, err := p.cloudconfig.SFSV2Client()
	if err != nil {
		return fmt.Errorf("failed to create SFS v2 client: %v", err)
	}

	// delete share
	err = DeleteShare(client, "")
	if err != nil {
		return fmt.Errorf("failed to delete share: %v", err)
	}

	return nil
}
