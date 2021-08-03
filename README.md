# algo-operator

# single group project - could be multi-group
```
operator-sdk create api --group cache --version v1alpha1 --kind Redisdb --resource --controller

./bin/controller-gen crd -w
./bin/controller-gen crd -h

operator-sdk create webhook --group batch --version v1 --kind CronJob --defaulting --programmatic-validation

```


# dbaas example:
```
operator-sdk init --domain redhat.com --repo github.com/RHEcosystemAppEng/dbaas-operator
operator-sdk create api --group dbaas --version v1 --kind DBaaSConnection --resource --controller
operator-sdk create api --group dbaas --version v1 --kind DBaaSInventory --resource --controller
make bundle
```