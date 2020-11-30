package vsi

import (
	"errors"
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"log"
)

// CreateInstance POST
// /instances/{instance_id}
// Create an instance action
func (s *Service) createInstance(name, profileName, imageID, zoneName string) (instance *vpcv1.Instance, response *core.DetailedResponse, err error) {

	options := &vpcv1.CreateInstanceOptions{}
	options.SetInstancePrototype(&vpcv1.InstancePrototype{
		Name: &name,
		Image: &vpcv1.ImageIdentity{
			ID: &imageID,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &profileName,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &zoneName,
		},
		PrimaryNetworkInterface: &vpcv1.NetworkInterfacePrototype{
			Subnet: &vpcv1.SubnetIdentity{
				ID: &s.subnetID,
			},
		},
		Keys: []vpcv1.KeyIdentityIntf{
			&vpcv1.KeyIdentity{
				ID: &s.sshKeyID,
			},
		},
		VPC: &vpcv1.VPCIdentity{
			ID: &s.vpcID,
		},
	})
	instance, response, err = s.vpcService.CreateInstance(options)
	return
}

func (s *Service) getOrCreateInstance(instanceName, profileName, imageID, zoneName string) error {
	if s.subnetID == "" || s.sshKeyID == "" || s.vpcID == "" {
		return errors.New("initialization is not done")
	}
	// List Instances
	options := &vpcv1.ListInstancesOptions{}
	instances, _, err := s.vpcService.ListInstances(options)

	if err != nil {
		log.Printf("Error from list call\n")
		return err
	}

	for _, inst := range instances.Instances {
		if *inst.Name == instanceName {
			s.instanceID = *inst.ID
			s.status = *inst.Status
			s.rulesAdded = true
			return err
		}
	}

	log.Printf("Instance doesn't exist, creating one.")

	// Create a instance
	vsi, _, err := s.createInstance(instanceName, profileName, imageID, zoneName)
	if err != nil {
		log.Printf("Error from Create Instance call\n")
		return err
	}

	s.instanceID = *vsi.ID
	s.status = *vsi.Status

	return nil
}

// DeleteInstance DELETE
// /instances/{id}
// Delete specified instance
func (s *Service) deleteInstance() (response *core.DetailedResponse, err error) {
	options := &vpcv1.DeleteInstanceOptions{}
	options.SetID(s.instanceID)
	response, err = s.vpcService.DeleteInstance(options)
	return response, err
}
