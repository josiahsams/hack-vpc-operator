package vsi

import (
	"errors"
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"log"
	"strings"
)

// CreateSubnetPublicGatewayBinding - PUT
// /subnets/{id}/public_gateway
// Attach a public gateway to a subnet
func (s *Service) bindPublicGatewayBindingToSubnet() error {
	subnetID := s.subnetID
	gatewayID := s.gatewayID
	if s.subnetID == "" {
		return errors.New("error subnet is not initialized")
	}
	goptions := &vpcv1.GetSubnetPublicGatewayOptions{}
	goptions.SetID(subnetID)
	pgw, _, err := s.vpcService.GetSubnetPublicGateway(goptions)
	if err != nil && !strings.Contains(err.Error(), "no public gateway associated") {
		return err
	}

	if pgw == nil {
		options := &vpcv1.SetSubnetPublicGatewayOptions{}
		options.SetID(subnetID)
		options.SetPublicGatewayIdentity(&vpcv1.PublicGatewayIdentity{ID: &gatewayID})

		if _, _, err = s.vpcService.SetSubnetPublicGateway(options); err != nil {
			return err
		}
	} else {
		if *pgw.ID != gatewayID {
			log.Printf("Gateway attached to a different Subnet")
		}
	}

	return nil
}

// CreatePublicGateway POST
// /public_gateways
// Create a public gateway
func (s *Service) createPublicGateway(name, zoneName string) (pgw *vpcv1.PublicGateway, response *core.DetailedResponse, err error) {
	if s.vpcID == "" {
		err = errors.New("error vpc is not initialized")
		return
	}
	options := &vpcv1.CreatePublicGatewayOptions{}
	options.SetVPC(&vpcv1.VPCIdentity{
		ID: &s.vpcID,
	})
	options.SetZone(&vpcv1.ZoneIdentity{
		Name: &zoneName,
	})
	options.SetName(name)
	pgw, response, err = s.vpcService.CreatePublicGateway(options)
	return
}

func (s *Service) getOrCreatePublicGateway(gatewayName string) error {
	if s.vpcID == "" {
		return errors.New("error vpc is not initialized")
	}

	// List Gateways
	listGatewayOptions := &vpcv1.ListPublicGatewaysOptions{}
	gateways, _, err := s.vpcService.ListPublicGateways(listGatewayOptions)

	if err != nil {
		log.Printf("Error from list call\n")
		return err
	}

	for _, gateway := range gateways.PublicGateways {
		if *gateway.Name == gatewayName {
			s.gatewayID = *gateway.ID
			return nil
		}
	}

	log.Printf("Gateway doesn't exist, creating one.")

	// Create a subnet
	vpcGateway, _, err := s.createPublicGateway(gatewayName, "us-south-3")
	if err != nil {
		log.Printf("Error from Create Gateway call\n")
		panic(err)
	}
	s.gatewayID = *vpcGateway.ID

	return nil
}

// DeleteSubnetPublicGatewayBinding - DELETE
// /subnets/{id}/public_gateway
// Detach a public gateway from a subnet
func (s *Service) deleteSubnetPublicGatewayBinding() (response *core.DetailedResponse, err error) {
	options := &vpcv1.UnsetSubnetPublicGatewayOptions{}
	options.SetID(s.subnetID)
	response, err = s.vpcService.UnsetSubnetPublicGateway(options)
	return response, err
}

// DeletePublicGateway DELETE
// /public_gateways/{id}
// Delete specified public gateway
func (s *Service) deletePublicGateway() (response *core.DetailedResponse, err error) {
	options := &vpcv1.DeletePublicGatewayOptions{}
	options.SetID(s.gatewayID)
	response, err = s.vpcService.DeletePublicGateway(options)
	return response, err
}
