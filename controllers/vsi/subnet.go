package vsi

import (
	"errors"
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"log"
)

// DeleteSubnet - DELETE
// /subnets/{id}
// Delete specified subnet
func (s *Service) deleteSubnet() (response *core.DetailedResponse, err error) {
	options := &vpcv1.DeleteSubnetOptions{}
	options.SetID(s.subnetID)
	response, err = s.vpcService.DeleteSubnet(options)
	return response, err
}

func (s *Service) getOrCreateSubnet(subnetName string) error {
	if s.vpcID == "" {
		return errors.New("error vpc is not initialized")
	}
	// List Subnets
	listSubnetsOptions := &vpcv1.ListSubnetsOptions{}
	subnets, _, err := s.vpcService.ListSubnets(listSubnetsOptions)

	if err != nil {
		log.Printf("Error from list call\n")
		return err
	}

	for _, subnet := range subnets.Subnets {
		if *subnet.Name == subnetName {
			s.subnetID = *subnet.ID
			return nil
		}
	}

	log.Printf("Subnet doesn't exist, creating one.")

	// Create a subnet
	vpcSubnet, _, err := s.createSubnet(s.vpcID, subnetName, "us-south-3", false)
	if err != nil {
		log.Printf("Error from Create Subnet call\n")
		return err
	}
	log.Printf(*vpcSubnet.ID)

	s.subnetID = *vpcSubnet.ID

	return nil
}

// CreateSubnet - POST
// /subnets
// Create a subnet
func (s *Service) createSubnet(vpcID, name, zone string, mock bool) (subnet *vpcv1.Subnet, response *core.DetailedResponse, err error) {
	options := &vpcv1.CreateSubnetOptions{}
	if mock {
		options.SetSubnetPrototype(&vpcv1.SubnetPrototype{
			Ipv4CIDRBlock: core.StringPtr("10.243.0.0/24"),
			Name:          &name,
			VPC: &vpcv1.VPCIdentity{
				ID: &vpcID,
			},
			Zone: &vpcv1.ZoneIdentity{
				Name: &zone,
			},
		})
	} else {
		options.SetSubnetPrototype(&vpcv1.SubnetPrototype{
			Name: &name,
			VPC: &vpcv1.VPCIdentity{
				ID: &vpcID,
			},
			Zone: &vpcv1.ZoneIdentity{
				Name: &zone,
			},
			TotalIpv4AddressCount: core.Int64Ptr(128),
		})
	}
	subnet, response, err = s.vpcService.CreateSubnet(options)
	return
}
