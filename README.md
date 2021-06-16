# algo-operator

# single group project - could be multi-group
operator-sdk create api --group cache --version v1alpha1 --kind Redisdb --resource --controller


# dbaas example:
operator-sdk init --domain redhat.com --repo github.com/RHEcosystemAppEng/dbaas-operator
operator-sdk create api --group dbaas --version v1 --kind DBaaSConnection --resource --controller
operator-sdk create api --group dbaas --version v1 --kind DBaaSInventory --resource --controller
make bundle