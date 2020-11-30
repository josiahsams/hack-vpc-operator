package vsi

import (
	"fmt"
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"log"
)

func (s *Service) bindFloatingIPToNetworkInterface() error {
	foptions := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{}
	foptions.SetInstanceID(s.instanceID)
	foptions.SetID(s.fipID)
	foptions.SetNetworkInterfaceID(s.networkID)
	fip, _, _ := s.vpcService.GetInstanceNetworkInterfaceFloatingIP(foptions)
	// if err != nil  {
	// 	return
	// }

	if fip == nil {
		options := &vpcv1.AddInstanceNetworkInterfaceFloatingIPOptions{}
		options.SetID(s.fipID)
		options.SetInstanceID(s.instanceID)
		options.SetNetworkInterfaceID(s.networkID)
		_, _, err := s.vpcService.AddInstanceNetworkInterfaceFloatingIP(options)

		if err != nil {
			return err
		}
	}
	return nil
}

// CreateFloatingIP - POST
// /floating_ips
// Reserve a floating IP
func (s *Service) createFloatingIP(zone, name string) (floatingIP *vpcv1.FloatingIP, response *core.DetailedResponse, err error) {
	options := &vpcv1.CreateFloatingIPOptions{}
	options.SetFloatingIPPrototype(&vpcv1.FloatingIPPrototype{
		Name: &name,
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
	})
	floatingIP, response, err = s.vpcService.CreateFloatingIP(options)
	return
}

func (s *Service) getOrCreateFloatingIP(zone, ipName string) error {
	// List Floating IP
	options := &vpcv1.ListFloatingIpsOptions{}
	fips, _, err := s.vpcService.ListFloatingIps(options)

	if err != nil {
		log.Printf("Error from list call\n")
		return err
	}

	for _, fip := range fips.FloatingIps {
		if *fip.Name == ipName {

			s.fipID = *fip.ID
			s.address = *fip.Address
			return nil
		}
	}

	log.Printf("Floating IP doesn't exist, creating one.")

	// Create a instance
	fip, _, err := s.createFloatingIP(zone, ipName)
	if err != nil {
		fmt.Printf("Error from Create Floating IP call\n")
		return err
	}

	s.fipID = *fip.ID
	s.address = *fip.Address

	return nil
}

// ReleaseFloatingIP - DELETE
// /floating_ips/{id}
// Release the specified floating IP
func (s *Service) releaseFloatingIP() (response *core.DetailedResponse, err error) {
	options := s.vpcService.NewDeleteFloatingIPOptions(s.fipID)
	response, err = s.vpcService.DeleteFloatingIP(options)
	return response, err
}
