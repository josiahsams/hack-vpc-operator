# hack-vpc-operator

### Overview
This project is created as part of the Sangam project : https://github.ibm.com/isdl/sangam/issues/330

In this project, we will focus on making use of the VPC APIs and come up an Openshift Operator. 
All the virtual server instances (VSI) created within a namespace should be created under the same VPC. 
Proper gateways should be created so that the VSIs should be reachable within the Openshift cluster. 
The SSH Keys should be part of the cluster secrets and should be used to access the respective VSIs.

### Demo

[![Demo](others/Screenshot.jpg?raw=true "Click to watch this Demo ")](https://ibm.box.com/s/e6otsa86xqmpjurykawb6p3ov41kkbbv)

Link: https://ibm.box.com/s/e6otsa86xqmpjurykawb6p3ov41kkbbv

### Developer Guide.

1. Any update to the Permission should be followed by,
```
make manifests
```

2. To locally test a controller in a Openshift environment, login to oc and run the following commands
```
oc login --token=sha256~XYZ.... --server=https://9.x.x.x:6443
make install
make run
```

3. To test on Openshift cluster, 

a. create a controller image as follows,
```
export VERSION=0.0.4
docker build -t quay.io/josiahsams/vsi-operator:v${VERSION} .
docker push quay.io/josiahsams/vsi-operator:v${VERSION}
```

b. To create a operator bundle,
```
export VERSION=0.0.6
export UNAME="josiahsams"
export BUNDLE_IMG=quay.io/$UNAME/vsi-bundle:v$VERSION

make bundle-build BUNDLE_IMG=$BUNDLE_IMG
make manifests
operator-sdk generate kustomize manifests

kustomize build config/manifests | operator-sdk generate bundle --overwrite --version $VERSION

# Check config/manager/kustomization.yaml for manager image name. 
# Check bundle/manifests/joe.clusterserviceversion.yaml for mgr image name & cluster permission.

operator-sdk bundle validate ./bundle

docker build -f bundle.Dockerfile -t $BUNDLE_IMG .
make docker-push IMG=$BUNDLE_IMG

opm index add --bundles quay.io/josiahsams/vsi-bundle:v${VERSION} --tag quay.io/josiahsams/vsi-catalog:v${VERSION} -c docker
docker push quay.io/josiahsams/vsi-catalog:v${VERSION}
```

c. Install the catalog source to Openshift,
```
# delete existing source
oc delete catalogsource/vsi-catalog -n olm

# Check the yaml for the right version in the catalog-source.
oc apply -f catalog-source.yaml
```

d. uninstall existing operator from the Operatorhub and then cleanup from CRDs,
```
oc delete crds/vsis.cloud.ibm.com
```

e. Operator should be visible in the Operatorhub in the dashboard

f. Install the CRDs to the cluster
```
oc apply -f ./config/samples/cloud_v1_vsi.yaml
```

g. Check status,
```
oc get vsi.cloud.ibm.com/vsi-hack1
oc get vsi.cloud.ibm.com/vsi-hack1 -o yaml
oc get services
oc get endpoints
oc get routes
```

h. Cleanup
```
oc delete vsi.cloud.ibm.com/vsi-hack1
```
