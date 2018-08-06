## Usage sfs-provisioner in openshift

1. Create a storage class named ```sfs-storage-class```.

```
oc create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/openshift/sc.yaml
```

2. Create a sfs pvc named ```sfs-pvc```.

```
oc create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/openshift/pvc.yaml
```

3. Create a nginx pod with sfs pvc.

```
oc create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/openshift/pod.yaml
```

If you want to create all of the above resources, you could run:

```
oc create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/openshift/
```
