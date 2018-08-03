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

package main

import (
	"flag"
	"os"

	"github.com/golang/glog"

	"github.com/kubernetes-incubator/external-storage/lib/controller"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/huaweicloud/external-sfs/pkg/config"
	"github.com/huaweicloud/external-sfs/pkg/sfs"
)

var (
	provisioner  = flag.String("provisioner", "external.k8s.io/sfs", "Name of the provisioner. The provisioner will only provision volumes for claims that request a StorageClass with a provisioner field set equal to this name.")
	master       = flag.String("master", "", "Master URL to build a client config from. Either this or kubeconfig needs to be set if the provisioner is being run out of cluster.")
	kubeconfig   = flag.String("kubeconfig", "", "Absolute path to the kubeconfig file. Either this or master needs to be set if the provisioner is being run out of cluster.")
	cloudconfig  = flag.String("cloudconfig", "/etc/origin/cloudprovider/openstack.conf", "Absolute path to the cloud config")
	sharetimeout = flag.Int("sharetimeout", 600, "Share operation timeout. Unit: second")
)

func main() {
	var restconfig *rest.Config
	var err error

	flag.Parse()
	flag.Set("logtostderr", "true")

	// get the KUBECONFIG from env if specified (useful for local/debug cluster)
	kubeconfigEnv := os.Getenv("KUBECONFIG")

	if kubeconfigEnv != "" {
		glog.Info("Found KUBECONFIG environment variable set, using that..")
		kubeconfig = &kubeconfigEnv
	}

	if *master != "" || *kubeconfig != "" {
		glog.Info("Either master or kubeconfig specified. building kube config from that..")
		restconfig, err = clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	} else {
		glog.Info("Building kube configs for running in cluster...")
		restconfig, err = rest.InClusterConfig()
	}
	if err != nil {
		glog.Fatalf("Failed to create restconfig: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(restconfig)
	if err != nil {
		glog.Fatalf("Failed to create client: %v", err)
	}

	cc, err := config.LoadConfig(*cloudconfig)
	if err != nil {
		glog.Fatalf("Failed to load cloud config: %v", err)
	}

	// The controller needs to know what the server version is because out-of-tree
	// provisioners aren't officially supported until 1.5
	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		glog.Fatalf("Error getting server version: %v", err)
	}

	glog.Infof("Get informations. server version: %s share time out: %d",
		serverVersion.GitVersion, *sharetimeout)

	provisionController := controller.NewProvisionController(
		clientset,
		*provisioner,
		sfs.NewProvisioner(clientset, cc, *sharetimeout),
		serverVersion.GitVersion,
	)

	provisionController.Run(wait.NeverStop)
}
