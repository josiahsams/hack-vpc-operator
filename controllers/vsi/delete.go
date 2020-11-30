package vsi

import (
	"log"
	"time"
)

// DeleteVSI ..
func DeleteVSI(s *Service) error {

	// unbind FloatingIP
	if _, err := s.deleteNetworkInterfaceFloatingIPBinding(); err != nil {
		return err
	}
	log.Println("DeleteNetworkInterfaceFloatingIPBinding Done ")
	// release floatingIp
	if _, err := s.releaseFloatingIP(); err != nil {
		return err
	}
	log.Println("ReleaseFloatingIP Done ")
	// delete instance
	if _, err := s.deleteInstance(); err != nil {
		return err
	}
	log.Println("DeleteInstance Done ")
	// unbind gateway to subnet
	if _, err := s.deleteSubnetPublicGatewayBinding(); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	log.Println("DeleteSubnetPublicGatewayBinding Done ")
	// delete gateway
	if _, err := s.deletePublicGateway(); err != nil {
		return err
	}
	log.Println("DeletePublicGateway Done ")
	// delete subnet
	if _, err := s.deleteSubnet(); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	log.Println("DeleteSubnet Done ")
	// delete VPC service
	if _, err := s.deleteVPC(); err != nil {
		return err
	}
	log.Println("DeleteVPC Done ")
	return nil
}
