# external-sfs
[![Go Report Card](https://goreportcard.com/badge/github.com/huaweicloud/external-sfs)](https://goreportcard.com/badge/github.com/huaweicloud/external-sfs)
[![Build Status](https://travis-ci.org/huaweicloud/external-sfs.svg?branch=master)](https://travis-ci.org/huaweicloud/external-sfs)
[![LICENSE](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://github.com/huaweicloud/external-sfs/blob/master/LICENSE)

Scalable File Service (SFS) provides completely hosted sharable file storage for Elastic Cloud Servers (ECSs)
on huawei clouds.
Compatible with the Network File System protocol, SFS is expandable to petabytes, features high performance,
and seamlessly handles data-intensive and bandwidth-intensive applications.

This repository houses external sfs provisioner for OpenShift and Kubernetes.

## Getting Started on OpenShift

### Deploy

external-sfs should be deployed in the OpenShift Master after OpenShift is deployed successfully.
In default, the Cloud Tenant informations are stored in the file ```/etc/origin/cloudprovider/openstack.conf``` of OpenShift Master. If your OpenShift Master contains the file ```/etc/origin/cloudprovider/openstack.conf```, please directly run the following command in your OpenShift Master.

```
oc adm policy add-scc-to-user privileged system:serviceaccount:default:sfs-provisioner
oc create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/deploy/sfs-provisioner/openshift/statefulset.yaml
```

If not, please firstly run the following command to download this repository,
```
git clone https://github.com/huaweicloud/external-sfs
```
and modify the statefulset.yaml,
```
vi external-sfs/deploy/sfs-provisioner/openshift/statefulset.yaml
```
and replace ```/etc/origin/cloudprovider/openstack.conf``` with your Cloud Config file in the line 73 of statefulset.yaml and replace the path ```/etc/origin``` with your Cloud Config directory in the line 82 of statefulset.yaml,

if you want to increase the log level, please add the following two lines after the line 73 of statefulset.yaml.

```
            - name: OS_DEBUG
              value: true
```

finally you can run the following command.
```
oc adm policy add-scc-to-user privileged system:serviceaccount:default:sfs-provisioner
oc create -f external-sfs/deploy/sfs-provisioner/openshift/statefulset.yaml
```

### Usage

```
oc create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/openshift/example.yaml
```

## Getting Started on Kubernetes

### Deploy

```
kubectl create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/deploy/sfs-provisioner/kubernetes/statefulset.yaml
```

### Usage

```
kubectl create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/kubernetes/example.yaml
```

## License

See the [LICENSE](LICENSE) file for details.
