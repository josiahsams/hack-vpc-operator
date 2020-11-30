module github.ibm.com/josiah-sams/hack-vpc-operator

go 1.13

require (
	github.com/IBM-Cloud/bluemix-go v0.0.0-20201119073718-c3ed816a263b
	github.com/IBM/go-sdk-core/v4 v4.8.2
	github.com/IBM/vpc-go-sdk v0.3.1
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/openshift/api v3.9.0+incompatible
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.3
)
