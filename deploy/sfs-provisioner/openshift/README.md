## Deploy sfs-provisioner in openshift

```
oc adm policy add-scc-to-user privileged system:serviceaccount:default:sfs-provisioner
oc create -f https://raw.githubusercontent.com/huaweicloud/external-sfs/master/deploy/sfs-provisioner/openshift/statefulset.yaml
```
