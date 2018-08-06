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

```
oc adm policy add-scc-to-user privileged system:serviceaccount:default:sfs-provisioner
oc create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/deploy/sfs-provisioner/openshift/statefulset.yaml
```

### Usage

```
oc create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/openshift/
```

## Getting Started on Kubernetes

### Deploy

```
kubectl create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/deploy/sfs-provisioner/kubernetes/statefulset.yaml
```

### Usage

```
kubectl create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/kubernetes/
```

## License

See the [LICENSE](LICENSE) file for details.
