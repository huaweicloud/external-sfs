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
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/sfs/v2/shares"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/kubernetes/pkg/controller/volume/persistentvolume"
)

// CreateShare in SFS
func CreateShare(client *golangsdk.ServiceClient, volOptions *controller.VolumeOptions) (*shares.Share, error) {
	// build share createOpts
	createOpts := shares.CreateOpts{}
	// build name
	createOpts.Name = "pvc-" + string(volOptions.PVC.GetUID())
	// build share proto
	createOpts.ShareProto = volOptions.Parameters[SFSParametersProtocol]
	if createOpts.ShareProto == "" {
		createOpts.ShareProto = SFSParametersProtocolDefault
	}
	// build size
	size, err := getStorageSize(volOptions.PVC)
	if err != nil {
		return nil, fmt.Errorf("Couldn't retrieve PVC storage size: %v", err)
	}
	createOpts.Size = size
	// build availability
	az := volOptions.Parameters[SFSParametersAvailability]
	if az != "" {
		createOpts.AvailabilityZone = az
	}
	// build type
	tp := volOptions.Parameters[SFSParametersType]
	if tp != "" {
		createOpts.ShareType = tp
	}
	// build metadata
	createOpts.Metadata = map[string]string{
		persistentvolume.CloudVolumeCreatedForClaimNamespaceTag: volOptions.PVC.Namespace,
		persistentvolume.CloudVolumeCreatedForClaimNameTag:      volOptions.PVC.Name,
		persistentvolume.CloudVolumeCreatedForVolumeNameTag:     createOpts.Name,
	}

	// create share
	glog.Infof("Create share createOpts: %v", createOpts)
	share, err := shares.Create(client, createOpts).Extract()
	if err != nil {
		return nil, fmt.Errorf("Couldn't create share in SFS: %v", err)
	}
	return share, nil
}

// WaitForShareStatus wait for share desired status until timeout
func WaitForShareStatus(client *golangsdk.ServiceClient, shareID string, desiredStatus string, timeout int) error {
	return golangsdk.WaitFor(timeout, func() (bool, error) {
		// reduce the amount of API calls
		time.Sleep(2 * time.Second)
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

// GrantAccess in SFS
func GrantAccess(client *golangsdk.ServiceClient, volOptions *controller.VolumeOptions, shareID string, vpcid string) error {
	// build GrantAccessOpts
	grantAccessOpts := shares.GrantAccessOpts{}
	grantAccessOpts.AccessLevel = "rw"
	grantAccessOpts.AccessType = "cert"
	// build vpcid
	id := volOptions.Parameters[SFSParametersVPCID]
	if id != "" {
		grantAccessOpts.AccessTo = id
	} else {
		grantAccessOpts.AccessTo = vpcid
	}

	// grant access
	result := shares.GrantAccess(client, shareID, grantAccessOpts)
	if result.Err != nil {
		return result.Err
	}
	return nil
}

// DeleteShare in SFS
func DeleteShare(client *golangsdk.ServiceClient, shareID string) error {
	result := shares.Delete(client, shareID)
	if result.Err != nil {
		return result.Err
	}
	return nil
}

// getStorageSize from pvc
func getStorageSize(pvc *v1.PersistentVolumeClaim) (int, error) {
	errStorageSizeNotConfigured := fmt.Errorf("Requested storage capacity must be set")

	if pvc.Spec.Resources.Requests == nil {
		return 0, errStorageSizeNotConfigured
	}

	storageSize, ok := pvc.Spec.Resources.Requests[v1.ResourceStorage]
	if !ok {
		return 0, errStorageSizeNotConfigured
	}

	if storageSize.IsZero() {
		return 0, fmt.Errorf("Requested storage size must not have zero value")
	}

	if storageSize.Sign() == -1 {
		return 0, fmt.Errorf("Requested storage size must be greater than zero")
	}

	var buf []byte
	canonicalValue, _ := storageSize.AsScale(resource.Giga)
	storageSizeAsByteSlice, _ := canonicalValue.AsCanonicalBytes(buf)

	return strconv.Atoi(string(storageSizeAsByteSlice))
}
