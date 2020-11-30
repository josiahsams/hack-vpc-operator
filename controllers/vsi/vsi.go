package vsi

import (
	"fmt"
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

// Service ..
type Service struct {
	vpcService                                                                                           *vpcv1.VpcV1
	securityGroupID, vpcID, subnetID, gatewayID, instanceID, networkID, fipID, sshKeyID, address, status string
	rulesAdded                                                                                           bool
}

// getService ..
func getService(apiKey string) (*Service, error) {
	service, err := instantiateVPCService(apiKey)
	if err != nil {
		return nil, err
	}
	return &Service{vpcService: service}, nil
}

// instantiateVPCService - Instantiate VPC Gen2 service
func instantiateVPCService(apiKey string) (*vpcv1.VpcV1, error) {

	authenticator := &core.IamAuthenticator{
		ApiKey: apiKey,
		URL:    authenticatorURL,
	}

	options := &vpcv1.VpcV1Options{
		Authenticator: authenticator,
	}

	service, err := vpcv1.NewVpcV1UsingExternalConfig(options)
	if err != nil {
		fmt.Println("Gen2 Service creation failed.", err)
		return nil, err
	}

	service.SetServiceURL(iaasURL)
	// Check successful instantiation

	// return new vpc gen2 service
	return service, nil
}
