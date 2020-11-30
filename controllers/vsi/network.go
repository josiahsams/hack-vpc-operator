package vsi

import (
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"log"
)

// DeleteNetworkInterfaceFloatingIPBinding DELETE
// /instances/{instance_id}/network_interfaces/{network_interface_id}/floating_ips/{id}
// Disassociate specified floating IP
func (s *Service) deleteNetworkInterfaceFloatingIPBinding() (response *core.DetailedResponse, err error) {
	options := &vpcv1.RemoveInstanceNetworkInterfaceFloatingIPOptions{}
	options.SetID(s.fipID)
	options.SetInstanceID(s.instanceID)
	options.SetNetworkInterfaceID(s.networkID)
	response, err = s.vpcService.RemoveInstanceNetworkInterfaceFloatingIP(options)
	return response, err
}

func (s *Service) addSecurityGroupRules() error {
	if s.rulesAdded {
		return nil
	}

	soptions := &vpcv1.GetSecurityGroupOptions{}
	soptions.SetID(s.securityGroupID)
	sg, _, err := s.vpcService.GetSecurityGroup(soptions)
	if err != nil {
		log.Printf("Error from GetSecurityGroup\n")
		return err
	}

	roptions1 := &vpcv1.CreateSecurityGroupRuleOptions{}
	roptions1.SetSecurityGroupID(*sg.ID)
	roptions1.SetSecurityGroupRulePrototype(&vpcv1.SecurityGroupRulePrototype{
		Direction: core.StringPtr("inbound"),
		Protocol:  core.StringPtr("tcp"),
		IPVersion: core.StringPtr("ipv4"),
		PortMin:   core.Int64Ptr(1),
		PortMax:   core.Int64Ptr(65535),
	})
	_, _, err = s.vpcService.CreateSecurityGroupRule(roptions1)
	if err != nil {
		log.Printf("Error from CreateSecurityGroupRule 1\n")
		return err
	}

	roptions2 := &vpcv1.CreateSecurityGroupRuleOptions{}
	roptions2.SetSecurityGroupID(*sg.ID)
	roptions2.SetSecurityGroupRulePrototype(&vpcv1.SecurityGroupRulePrototype{
		Direction: core.StringPtr("inbound"),
		Protocol:  core.StringPtr("icmp"),
		IPVersion: core.StringPtr("ipv4"),
	})
	_, _, err = s.vpcService.CreateSecurityGroupRule(roptions2)
	if err != nil {
		log.Printf("Error from CreateSecurityGroupRule 2\n")
		return err

	}

	s.rulesAdded = true

	return nil
}

// ListNetworkInterfaces GET
// /instances/{instance_id}/network_interfaces
// List all network interfaces on an instance
func (s *Service) listNetworkInterfaces() error {
	options := &vpcv1.ListInstanceNetworkInterfacesOptions{}
	options.SetInstanceID(s.instanceID)
	networkInterfaces, _, err := s.vpcService.ListInstanceNetworkInterfaces(options)
	if err != nil {
		return err
	}
	s.networkID = *networkInterfaces.NetworkInterfaces[0].ID
	s.securityGroupID = *networkInterfaces.NetworkInterfaces[0].SecurityGroups[0].ID
	return err
}
