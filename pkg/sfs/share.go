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

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/sfs/v2/shares"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/kubernetes/pkg/controller/volume/persistentvolume"
)

// CreateShare in SFS
func CreateShare(client *golangsdk.ServiceClient, volOptions *controller.VolumeOptions) (*shares.Share, error) {
	// getStorageSizeInGiga
	/*size, err := getStorageSizeInGiga(volOptions.PVC)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve PVC storage size: %v", err)
	}*/

	// build share createOpts
	createOpts := shares.CreateOpts{}
	createOpts.ShareProto = "NFS"
	createOpts.Size = 10
	createOpts.Name = "pvc"
	createOpts.AvailabilityZone = ""
	createOpts.ShareNetworkID = ""
	createOpts.Metadata = map[string]string{
		persistentvolume.CloudVolumeCreatedForClaimNamespaceTag: volOptions.PVC.Namespace,
		persistentvolume.CloudVolumeCreatedForClaimNameTag:      volOptions.PVC.Name,
		persistentvolume.CloudVolumeCreatedForVolumeNameTag:     createOpts.Name,
	}

	// create share
	share, err := shares.Create(client, createOpts).Extract()
	if err != nil {
		return nil, fmt.Errorf("couldn't create share in SFS: %v", err)
	}
	return share, nil
}

// WaitForShareStatus wait for share desired status until timeout
func WaitForShareStatus(client *golangsdk.ServiceClient, shareID string, desiredStatus string, timeout int) error {
	return golangsdk.WaitFor(timeout, func() (bool, error) {
		share, err := GetShare(client, shareID)
		if err != nil {
			return false, err
		}
		return share.Status == desiredStatus, nil
	})
}

// GetShare in SFS
func GetShare(client *golangsdk.ServiceClient, shareID string) (*shares.Share, error) {
	return shares.Get(client, shareID).Extract()
}

// DeleteShare in SFS
func DeleteShare(client *golangsdk.ServiceClient, shareID string) error {
	result := shares.Delete(client, shareID)
	if result.Err != nil {
		return result.Err
	}
	return nil
}
