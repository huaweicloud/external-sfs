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
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// Provisioner struct, implements controller.Provisioner interface
type Provisioner struct {
	clientset clientset.Interface
}

// NewProvisioner creates a new instance of sfs provisioner
func NewProvisioner(c clientset.Interface) *Provisioner {
	return &Provisioner{
		clientset: c,
	}
}

// Provision a share in sfs
func (p *Provisioner) Provision(volOptions controller.VolumeOptions) (*v1.PersistentVolume, error) {
	return nil, nil
}

// Delete a share from sfs
func (p *Provisioner) Delete(pv *v1.PersistentVolume) error {
	return nil
}
