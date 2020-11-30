package vsi

import (
	"log"
)

// GetAddress ..
func (s *Service) GetAddress() string {
	return s.address
}

// GetStatus ..
func (s *Service) GetStatus() string {
	return s.status
}

// GetOrCreate ..
func GetOrCreate(apiKey, userKey string) (s *Service, err error) {

	s, err = getService(apiKey)
	if err != nil {
		log.Printf("Error from Get VPC service call\n")
		return
	}
	s.rulesAdded = false
	return s, GetOrCreateWithObj(s, apiKey, userKey)
}

// GetOrCreateWithObj ..
func GetOrCreateWithObj(s *Service, apiKey, userKey string) (err error) {

	err = s.getOrCreateVPC("hack-test-vpc", apiKey)
	if err != nil {
		log.Printf("Error from Create VPC call\n")
		return
	}

	log.Println("VPC ID: ", s.vpcID)

	err = s.getOrCreateSubnet("hack-test-subnet")
	if err != nil {
		log.Printf("Error from Create Subnet call\n")
		return
	}
	log.Println("Subnet ID: ", s.subnetID)

	// Create a public gateway
	err = s.getOrCreatePublicGateway("hack-test-pgw")
	if err != nil {
		log.Printf("Error from Create Public Gateway call\n")
		return
	}

	log.Println("Gateway ID: ", s.gatewayID)

	// Bind Subnet to the Public Gateway
	err = s.bindPublicGatewayBindingToSubnet()
	if err != nil {
		log.Printf("Error from Bind Subnet Public Gateway call\n")
		return
	}

	log.Println("Bind Public Gateway to Subnet: succeeded")

	// Set SSH key access
	err = s.getOrCreateSSHKey("mackey1", userKey)
	if err != nil {
		log.Printf("Error from Create Public Gateway call\n")
		return
	}

	log.Println("SSHKey  ID: ", s.sshKeyID)

	// VSI instance setup
	err = s.getOrCreateInstance("hack-test-vsi", plan,
		imageID, region)
	if err != nil {
		log.Printf("Error from Create Instance call\n")
		return
	}
	log.Println("Instance ID: ", s.instanceID)

	// Create a Network Interface
	err = s.listNetworkInterfaces()
	if err != nil {
		log.Printf("Error from Create Network Interface call\n")
		return
	}
	log.Println("Network interface: ", s.networkID)

	err = s.addSecurityGroupRules()
	if err != nil {
		log.Printf("Error from Create Security Group Rules call\n")
		return
	}
	log.Println("Add Security Group Rules succeeded")

	// Floating IP setup
	err = s.getOrCreateFloatingIP(region, "hack-test-fip")
	if err != nil {
		log.Printf("Error from Create Floating IP call\n")
		return
	}
	log.Println("FloatingIP: ", s.fipID)

	// Bind Floating IP to Network Interface
	err = s.bindFloatingIPToNetworkInterface()
	if err != nil {
		log.Printf("Error from Bind Floating IP call\n")
		return
	}

	return
}
