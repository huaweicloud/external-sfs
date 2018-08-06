## Usage sfs-provisioner in kubernetes

1. Create a storage class named ```sfs-storage-class```.

```
kubectl create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/kubernetes/sc.yaml
```

2. Create a sfs pvc named ```sfs-pvc```.

```
kubectl create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/kubernetes/pvc.yaml
```

3. Create a nginx pod with sfs pvc.

```
kubectl create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/kubernetes/pod.yaml
```

If you want to create all of the above resources, you could run:

```
kubectl create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/examples/sfs-provisioner/kubernetes/example.yaml
```
