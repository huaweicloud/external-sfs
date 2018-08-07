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

// Provisioner implements controller.Provisioner interface
type Provisioner struct {
	clientset    clientset.Interface
	cloudconfig  config.CloudCredentials
	sharetimeout int
	vpcid        string
}

// NewProvisioner creates a new instance of sfs provisioner
func NewProvisioner(c clientset.Interface, cc config.CloudCredentials, timeout int, vpcid string) *Provisioner {

	// init backends for provisioner
	InitBackends()

	// init vpc for provisioner
	if vpcid == "" {
		vpcid = InitVPC(cc)
	}

	// return provisioner instance
	return &Provisioner{
		clientset:    c,
		cloudconfig:  cc,
		sharetimeout: timeout,
		vpcid:        vpcid,
	}
}

// Provision a share in sfs
func (p *Provisioner) Provision(volOptions controller.VolumeOptions) (*v1.PersistentVolume, error) {

	// selector check
	glog.Infof("Provision volOptions: %v", volOptions)
	if volOptions.PVC.Spec.Selector != nil {
		return nil, fmt.Errorf("Claim Selector is not supported")
	}

	// init sfs client
	glog.Info("Init sfs client...")
	client, err := p.cloudconfig.SFSV2Client()
	if err != nil {
		return nil, fmt.Errorf("Failed to create SFS v2 client: %v", err)
	}

	// create share
	glog.Info("Create share begin...")
	share, err := CreateShare(client, &volOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to create share: %v", err)
	}

	// wait fo share available
	glog.Infof("Wait fo share available: %s", share.ID)
	err = WaitForShareStatus(client, share.ID, SFSStatusAvailable, p.sharetimeout)
	if err != nil {
		return nil, fmt.Errorf("Waiting for share %s to become created failed: %v", share.ID, err)
	}

	// get new share
	glog.Infof("Get share: %s", share.ID)
	share, err = GetShare(client, share.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get share: %v", err)
	}

	// grant access
	glog.Infof("Grant access: %s", share.ID)
	err = GrantAccess(client, &volOptions, share.ID, p.vpcid)
	if err != nil {
		return nil, fmt.Errorf("Failed to grant access: %v", err)
	}

	// get location
	location := share.ExportLocation
	if (len(location) == 0) && (len(share.ExportLocations) > 0) {
		location = share.ExportLocations[0]
	}
	if len(location) == 0 {
		return nil, fmt.Errorf("Failed to get share %s location", share.ID)
	}
	glog.Infof("Get share: %s location: %s", share.ID, location)

	// get backend
	glog.Infof("Get backend: %s share protocal: %s", share.ID, share.ShareProto)
	b, err := GetBackend(share.ShareProto)
	if err != nil {
		return nil, fmt.Errorf("failed to get backend: %v", err)
	}

	// get persistent volume source
	glog.Infof("Build source from share: %v", share)
	pvsource, err := b.BuildSource(&backends.BuildSourceArgs{Location: location})
	if err != nil {
		return nil, fmt.Errorf("Failed to build source from backend: %v", err)
	}

	return &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: volOptions.PVName,
			Annotations: map[string]string{
				SFSAnnotationID: share.ID,
			},
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
	glog.Info("Init sfs client...")
	client, err := p.cloudconfig.SFSV2Client()
	if err != nil {
		return fmt.Errorf("Failed to create SFS v2 client: %v", err)
	}

	// get share id
	var shareid string
	shareid, ok := pv.ObjectMeta.Annotations[SFSAnnotationID]
	if (!ok) || (shareid == "") {
		return fmt.Errorf("Failed to get share id: %v", pv)
	}

	// delete share
	glog.Infof("Delete share: %s", shareid)
	err = DeleteShare(client, shareid)
	if err != nil {
		return fmt.Errorf("failed to delete share: %v", err)
	}

	return nil
}
