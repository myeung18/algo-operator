# algo-operator

# single group project - could be multi-group
```
operator-sdk create api --group cache --version v1alpha1 --kind Redisdb --resource --controller

./bin/controller-gen crd -w
./bin/controller-gen crd -h

operator-sdk create webhook --group batch --version v1 --kind CronJob --defaulting --programmatic-validation

```

```shell
# check role/sa/binding
oc get clusterrolebinding -o wide | grep algo
oc describe sa algo-operator-controller-manager

oc get clusterroles
oc get rolebindings -o wide | grep algo
oc get roles
oc describe clusterrolebinding algo-operator.v0.0.1-5f5ddccc66
oc get rolebinding.rbac
oc describe rolebing.rbac <...>
oc describe role algo-operator.v0.0.1-5f5ddccc66
oc describe clusterrole.rbac algo-operator.v0.0.1-5f5ddccc66

#add role to sa
oc create clusterrole route-create --verb=create --resource=route
oc policy add-role-to-user route-create -z algo-operator-controller-manager

```
```shell
# build with compatible operator-sdk and deploy with operator-sdk run
operator-sdk run bundle quay.io/myeung/algo-operator-bundle:v0.0.1

#apply CR, then expose srv
oc expose svc algocoding-sample-service-32000

```


# dbaas example:
```
operator-sdk init --domain redhat.com --repo github.com/RHEcosystemAppEng/dbaas-operator
operator-sdk create api --group dbaas --version v1 --kind DBaaSConnection --resource --controller
operator-sdk create api --group dbaas --version v1 --kind DBaaSInventory --resource --controller
make bundle
```