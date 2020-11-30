package vsi

import (
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"log"
)

// CreateVPC - POST
// /vpcs
// Create a VPC
func (s *Service) createVPC(name, resourceGroup string) (vpc *vpcv1.VPC, response *core.DetailedResponse, err error) {
	options := &vpcv1.CreateVPCOptions{}

	options.SetResourceGroup(&vpcv1.ResourceGroupIdentity{
		ID: &resourceGroup,
	})
	options.SetName(name)
	vpc, response, err = s.vpcService.CreateVPC(options)
	return
}

// DeleteVPC - DELETE
// /vpcs/{id}
// Delete specified VPC
func (s *Service) deleteVPC() (response *core.DetailedResponse, err error) {
	deleteVpcOptions := &vpcv1.DeleteVPCOptions{}
	deleteVpcOptions.SetID(s.vpcID)
	response, err = s.vpcService.DeleteVPC(deleteVpcOptions)
	return response, err
}

// getOrCreateVPC ..
func (s *Service) getOrCreateVPC(vpcName, apiKey string) error {
	// List VPCs
	listVpcsOptions := &vpcv1.ListVpcsOptions{}
	vpcs, _, err := s.vpcService.ListVpcs(listVpcsOptions)
	if err != nil {
		log.Printf("Error from list call\n")
		return err
	}

	for _, vpc := range vpcs.Vpcs {
		if *vpc.Name == vpcName {
			s.vpcID = *vpc.ID
			return nil
		}
	}

	log.Printf("VPC doesn't exist, creating one.")
	resourceGroupID, err := getResourceGroupID(apiKey)
	if err != nil {
		log.Printf("Error from Create VPC call\n")
		return err
	}

	vpcNew, _, err := s.createVPC(vpcName, resourceGroupID)
	if err != nil {
		log.Printf("Error from Create VPC call\n")
		return err
	}

	s.vpcID = *vpcNew.ID

	return nil
}
