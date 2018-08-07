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
	"io/ioutil"
	"strings"

	"github.com/golang/glog"
	"github.com/huaweicloud/external-sfs/pkg/config"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/pagination"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
)

// InitVPC for share
func InitVPC(cc config.CloudCredentials) string {
	// define vpcid
	vpcid := ""

	// get instanceid
	glog.Info("get instance id...")
	instanceID := readInstanceID()
	if instanceID == "" {
		glog.Info("Failed to get instance id")
		return vpcid
	}

	// init compute client
	glog.Info("Init compute client...")
	computeClient, err := cc.ComputeV2Client()
	if err != nil {
		glog.Infof("Failed to create compute v2 client: %v", err)
		return vpcid
	}

	// init sfs client
	glog.Info("Init network client...")
	networkClient, err := cc.NetworkingV1Client()
	if err != nil {
		glog.Infof("Failed to create network v2 client: %v", err)
		return vpcid
	}

	// get interfaces
	interfaces, err := getAttachedInterfacesByID(computeClient, instanceID)
	if err != nil {
		glog.Infof("Failed to get interfaces: %v", err)
		return vpcid
	}

	// get subnet id
	subnetid := ""
	for _, intf := range interfaces {
		if intf.NetID != "" {
			subnetid = intf.NetID
			break
		}
	}
	if subnetid == "" {
		glog.Info("Failed to get subnet id")
		return vpcid
	}
	glog.Infof("Get subnet id: %s", subnetid)

	// get routers
	subnet, err := subnets.Get(networkClient, subnetid).Extract()
	if err != nil {
		glog.Infof("Failed to get subnet: %v", err)
		return vpcid
	}

	// get vpc id
	vpcid = subnet.VPC_ID
	glog.Infof("Get vpc id: %s", vpcid)

	return vpcid
}

// readInstanceID from local file
func readInstanceID() string {
	const instanceIDFile = "/var/lib/cloud/data/instance-id"
	idBytes, err := ioutil.ReadFile(instanceIDFile)
	if err != nil {
		glog.Infof("Failed to get instance id from file: %v", err)
		return ""
	} else {
		instanceID := string(idBytes)
		instanceID = strings.TrimSpace(instanceID)
		glog.Infof("Get instance id from file: %s", instanceID)
		return instanceID
	}
}

// getAttachedInterfacesByID returns the node interfaces of the specified instance.
func getAttachedInterfacesByID(client *gophercloud.ServiceClient, instanceID string) ([]attachinterfaces.Interface, error) {
	var interfaces []attachinterfaces.Interface

	pager := attachinterfaces.List(client, instanceID)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		s, err := attachinterfaces.ExtractInterfaces(page)
		if err != nil {
			return false, err
		}
		interfaces = append(interfaces, s...)
		return true, nil
	})
	if err != nil {
		return interfaces, err
	}

	return interfaces, nil
}
