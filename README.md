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
