/**
 * (C) Copyright IBM Corp. 2020.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cloudhandler

import (
	"fmt"
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

// InstantiateVPCService - Instantiate VPC Gen2 service
func InstantiateVPCService() *vpcv1.VpcV1 {

	authenticator := &core.IamAuthenticator{
		ApiKey: "H20bJtI7rcyG9VgbhE3ohom6KbmYHZ8TT-dUu3M0oZrP",
		URL:    "https://iam.test.cloud.ibm.com/identity/token",
	}

	options := &vpcv1.VpcV1Options{
		Authenticator: authenticator,
	}

	service, serviceErr := vpcv1.NewVpcV1UsingExternalConfig(options)

	service.SetServiceURL("https://us-south-stage01.iaasdev.cloud.ibm.com/v1")
	// Check successful instantiation
	if serviceErr != nil {
		fmt.Println("Gen2 Service creation failed.", serviceErr)
		return nil
	}
	// return new vpc gen2 service
	return service
}

// ListRegions - List all regions
// GET
// /regions
func ListRegions(gen2 *vpcv1.VpcV1) (regions *vpcv1.RegionCollection, response *core.DetailedResponse, err error) {
	listRegionsOptions := &vpcv1.ListRegionsOptions{}
	regions, response, err = gen2.ListRegions(listRegionsOptions)
	return
}

// GetFloatingIP - GET
// /floating_ips/{id}
// Retrieve the specified floating IP
func GetFloatingIP(vpcService *vpcv1.VpcV1, id string) (floatingIP *vpcv1.FloatingIP, response *core.DetailedResponse, err error) {
	options := vpcService.NewGetFloatingIPOptions(id)
	floatingIP, response, err = vpcService.GetFloatingIP(options)
	return
}

// ReleaseFloatingIP - DELETE
// /floating_ips/{id}
// Release the specified floating IP
func ReleaseFloatingIP(vpcService *vpcv1.VpcV1, id string) (response *core.DetailedResponse, err error) {
	options := vpcService.NewDeleteFloatingIPOptions(id)
	response, err = vpcService.DeleteFloatingIP(options)
	return response, err
}

// UpdateFloatingIP - PATCH
// /floating_ips/{id}
// Update the specified floating IP
func UpdateFloatingIP(vpcService *vpcv1.VpcV1, id, name string) (floatingIP *vpcv1.FloatingIP, response *core.DetailedResponse, err error) {
	body := &vpcv1.FloatingIPPatch{
		Name: &name,
	}
	patchBody, _ := body.AsPatch()
	options := &vpcv1.UpdateFloatingIPOptions{
		ID:              &id,
		FloatingIPPatch: patchBody,
	}

	floatingIP, response, err = vpcService.UpdateFloatingIP(options)
	return
}

// CreateFloatingIP - POST
// /floating_ips
// Reserve a floating IP
func CreateFloatingIP(vpcService *vpcv1.VpcV1, zone, name string) (floatingIP *vpcv1.FloatingIP, response *core.DetailedResponse, err error) {
	options := &vpcv1.CreateFloatingIPOptions{}
	options.SetFloatingIPPrototype(&vpcv1.FloatingIPPrototype{
		Name: &name,
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
	})
	floatingIP, response, err = vpcService.CreateFloatingIP(options)
	return
}

/**
 * SSH Keys
 *
 */

// ListKeys - GET
// /keys
// List all keys
func ListKeys(vpcService *vpcv1.VpcV1) (keys *vpcv1.KeyCollection, response *core.DetailedResponse, err error) {
	listKeysOptions := &vpcv1.ListKeysOptions{}
	keys, response, err = vpcService.ListKeys(listKeysOptions)
	return
}

// GetSSHKey - GET
// /keys/{id}
// Retrieve specified key
func GetSSHKey(vpcService *vpcv1.VpcV1, id string) (key *vpcv1.Key, response *core.DetailedResponse, err error) {
	getKeyOptions := &vpcv1.GetKeyOptions{}
	getKeyOptions.SetID(id)
	key, response, err = vpcService.GetKey(getKeyOptions)
	return
}

// UpdateSSHKey - PATCH
// /keys/{id}
// Update specified key
func UpdateSSHKey(vpcService *vpcv1.VpcV1, id, name string) (key *vpcv1.Key, response *core.DetailedResponse, err error) {
	body := &vpcv1.KeyPatch{
		Name: &name,
	}
	patchBody, _ := body.AsPatch()
	updateKeyOptions := vpcService.NewUpdateKeyOptions(id, patchBody)
	key, response, err = vpcService.UpdateKey(updateKeyOptions)
	return
}

// DeleteSSHKey - DELETE
// /keys/{id}
// Delete specified key
func DeleteSSHKey(vpcService *vpcv1.VpcV1, id string) (response *core.DetailedResponse, err error) {
	deleteKeyOptions := &vpcv1.DeleteKeyOptions{}
	deleteKeyOptions.SetID(id)
	response, err = vpcService.DeleteKey(deleteKeyOptions)
	return response, err
}

// CreateSSHKey - POST
// /keys
// Create a key
func CreateSSHKey(vpcService *vpcv1.VpcV1, name, publicKey string) (key *vpcv1.Key, response *core.DetailedResponse, err error) {
	options := &vpcv1.CreateKeyOptions{}
	options.SetName(name)
	options.SetPublicKey(publicKey)
	key, response, err = vpcService.CreateKey(options)
	return
}

// ListVpcs - GET
// /vpcs
// List all VPCs
func ListVpcs(vpcService *vpcv1.VpcV1) (vpcs *vpcv1.VPCCollection, response *core.DetailedResponse, err error) {
	listVpcsOptions := &vpcv1.ListVpcsOptions{}
	vpcs, response, err = vpcService.ListVpcs(listVpcsOptions)
	return
}

// GetVPC - GET
// /vpcs/{id}
// Retrieve specified VPC
func GetVPC(vpcService *vpcv1.VpcV1, id string) (vpc *vpcv1.VPC, response *core.DetailedResponse, err error) {
	getVpcOptions := &vpcv1.GetVPCOptions{}
	getVpcOptions.SetID(id)
	vpc, response, err = vpcService.GetVPC(getVpcOptions)
	return
}

// DeleteVPC - DELETE
// /vpcs/{id}
// Delete specified VPC
func DeleteVPC(vpcService *vpcv1.VpcV1, id string) (response *core.DetailedResponse, err error) {
	deleteVpcOptions := &vpcv1.DeleteVPCOptions{}
	deleteVpcOptions.SetID(id)
	response, err = vpcService.DeleteVPC(deleteVpcOptions)
	return response, err
}

// UpdateVPC - PATCH
// /vpcs/{id}
// Update specified VPC
func UpdateVPC(vpcService *vpcv1.VpcV1, id, name string) (vpc *vpcv1.VPC, response *core.DetailedResponse, err error) {
	body := &vpcv1.VPCPatch{
		Name: &name,
	}
	patchBody, _ := body.AsPatch()
	options := vpcService.NewUpdateVPCOptions(id, patchBody)
	vpc, response, err = vpcService.UpdateVPC(options)
	return
}

// CreateVPC - POST
// /vpcs
// Create a VPC
func CreateVPC(vpcService *vpcv1.VpcV1, name, resourceGroup string) (vpc *vpcv1.VPC, response *core.DetailedResponse, err error) {
	options := &vpcv1.CreateVPCOptions{}

	options.SetResourceGroup(&vpcv1.ResourceGroupIdentity{
		ID: &resourceGroup,
	})
	options.SetName(name)
	vpc, response, err = vpcService.CreateVPC(options)
	return
}

/**
 * Public Gateway
 *
 */

// ListPublicGateways GET
// /public_gateways
// List all public gateways
func ListPublicGateways(vpcService *vpcv1.VpcV1) (pgws *vpcv1.PublicGatewayCollection, response *core.DetailedResponse, err error) {
	options := &vpcv1.ListPublicGatewaysOptions{}
	pgws, response, err = vpcService.ListPublicGateways(options)
	return
}

// CreatePublicGateway POST
// /public_gateways
// Create a public gateway
func CreatePublicGateway(vpcService *vpcv1.VpcV1, name, vpcID, zoneName string) (pgw *vpcv1.PublicGateway, response *core.DetailedResponse, err error) {
	options := &vpcv1.CreatePublicGatewayOptions{}
	options.SetVPC(&vpcv1.VPCIdentity{
		ID: &vpcID,
	})
	options.SetZone(&vpcv1.ZoneIdentity{
		Name: &zoneName,
	})
	pgw, response, err = vpcService.CreatePublicGateway(options)
	return
}

// DeletePublicGateway DELETE
// /public_gateways/{id}
// Delete specified public gateway
func DeletePublicGateway(vpcService *vpcv1.VpcV1, id string) (response *core.DetailedResponse, err error) {
	options := &vpcv1.DeletePublicGatewayOptions{}
	options.SetID(id)
	response, err = vpcService.DeletePublicGateway(options)
	return response, err
}

// GetPublicGateway GET
// /public_gateways/{id}
// Retrieve specified public gateway
func GetPublicGateway(vpcService *vpcv1.VpcV1, id string) (pgw *vpcv1.PublicGateway, response *core.DetailedResponse, err error) {
	options := &vpcv1.GetPublicGatewayOptions{}
	options.SetID(id)
	pgw, response, err = vpcService.GetPublicGateway(options)
	return
}

// UpdatePublicGateway PATCH
// /public_gateways/{id}
// Update a public gateway's name
func UpdatePublicGateway(vpcService *vpcv1.VpcV1, id, name string) (pgw *vpcv1.PublicGateway, response *core.DetailedResponse, err error) {
	body := &vpcv1.PublicGatewayPatch{
		Name: &name,
	}
	patchBody, _ := body.AsPatch()
	options := &vpcv1.UpdatePublicGatewayOptions{
		PublicGatewayPatch: patchBody,
		ID:                 &id,
	}
	pgw, response, err = vpcService.UpdatePublicGateway(options)
	return
}

/**
 * Instances
 *
 */

// ListInstanceProfiles - GET
// /instance/profiles
// List all instance profiles
func ListInstanceProfiles(vpcService *vpcv1.VpcV1) (profiles *vpcv1.InstanceProfileCollection, response *core.DetailedResponse, err error) {
	options := &vpcv1.ListInstanceProfilesOptions{}
	profiles, response, err = vpcService.ListInstanceProfiles(options)
	return
}

// GetInstanceProfile - GET
// /instance/profiles/{name}
// Retrieve specified instance profile
func GetInstanceProfile(vpcService *vpcv1.VpcV1, profileName string) (profile *vpcv1.InstanceProfile, response *core.DetailedResponse, err error) {
	options := &vpcv1.GetInstanceProfileOptions{}
	options.SetName(profileName)
	profile, response, err = vpcService.GetInstanceProfile(options)
	return
}

// ListInstances GET
// /instances
// List all instances
func ListInstances(vpcService *vpcv1.VpcV1) (instances *vpcv1.InstanceCollection, response *core.DetailedResponse, err error) {
	options := &vpcv1.ListInstancesOptions{}
	instances, response, err = vpcService.ListInstances(options)
	return
}

// GetInstance GET
// instances/{id}
// Retrieve an instance
func GetInstance(vpcService *vpcv1.VpcV1, instanceID string) (instance *vpcv1.Instance, response *core.DetailedResponse, err error) {
	options := &vpcv1.GetInstanceOptions{}
	options.SetID(instanceID)
	instance, response, err = vpcService.GetInstance(options)
	return
}

// DeleteInstance DELETE
// /instances/{id}
// Delete specified instance
func DeleteInstance(vpcService *vpcv1.VpcV1, id string) (response *core.DetailedResponse, err error) {
	options := &vpcv1.DeleteInstanceOptions{}
	options.SetID(id)
	response, err = vpcService.DeleteInstance(options)
	return response, err
}

// UpdateInstance PATCH
// /instances/{id}
// Update specified instance
func UpdateInstance(vpcService *vpcv1.VpcV1, id, name string) (instance *vpcv1.Instance, response *core.DetailedResponse, err error) {
	body := &vpcv1.InstancePatch{
		Name: &name,
	}
	patchBody, _ := body.AsPatch()
	options := &vpcv1.UpdateInstanceOptions{
		InstancePatch: patchBody,
		ID:            &id,
	}
	instance, response, err = vpcService.UpdateInstance(options)
	return
}

// CreateInstance POST
// /instances/{instance_id}
// Create an instance action
func CreateInstance(vpcService *vpcv1.VpcV1, name, profileName, imageID, zoneName, subnetID, sshkeyID, vpcID string) (instance *vpcv1.Instance, response *core.DetailedResponse, err error) {
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
				ID: &subnetID,
			},
		},
		Keys: []vpcv1.KeyIdentityIntf{
			&vpcv1.KeyIdentity{
				ID: &sshkeyID,
			},
		},
		VPC: &vpcv1.VPCIdentity{
			ID: &vpcID,
		},
	})
	instance, response, err = vpcService.CreateInstance(options)
	return
}

// CreateInstanceAction PATCH
// /instances/{instance_id}/actions
// Update specified instance
func CreateInstanceAction(vpcService *vpcv1.VpcV1, instanceID, typeOfAction string) (action *vpcv1.InstanceAction, response *core.DetailedResponse, err error) {
	options := &vpcv1.CreateInstanceActionOptions{}
	options.SetInstanceID(instanceID)
	options.SetType(typeOfAction)
	action, response, err = vpcService.CreateInstanceAction(options)
	return
}

// GetInstanceInitialization GET
// /instances/{id}/initialization
// Retrieve configuration used to initialize the instance.
func GetInstanceInitialization(vpcService *vpcv1.VpcV1, instanceID string) (initData *vpcv1.InstanceInitialization, response *core.DetailedResponse, err error) {
	options := &vpcv1.GetInstanceInitializationOptions{}
	options.SetID(instanceID)
	initData, response, err = vpcService.GetInstanceInitialization(options)
	return
}

// ListNetworkInterfaces GET
// /instances/{instance_id}/network_interfaces
// List all network interfaces on an instance
func ListNetworkInterfaces(vpcService *vpcv1.VpcV1, id string) (networkInterfaces *vpcv1.NetworkInterfaceUnpaginatedCollection, response *core.DetailedResponse, err error) {
	options := &vpcv1.ListInstanceNetworkInterfacesOptions{}
	options.SetInstanceID(id)
	networkInterfaces, response, err = vpcService.ListInstanceNetworkInterfaces(options)
	return
}

// CreateNetworkInterface POST
// /instances/{instance_id}/network_interfaces
// List all network interfaces on an instance
func CreateNetworkInterface(vpcService *vpcv1.VpcV1, id, subnetID string) (networkInterface *vpcv1.NetworkInterface, response *core.DetailedResponse, err error) {
	options := &vpcv1.CreateInstanceNetworkInterfaceOptions{}
	options.SetInstanceID(id)
	options.SetName("eth1")
	options.SetSubnet(&vpcv1.SubnetIdentityByID{
		ID: &subnetID,
	})
	networkInterface, response, err = vpcService.CreateInstanceNetworkInterface(options)
	return
}

// DeleteNetworkInterface Delete
// /instances/{instance_id}/network_interfaces/{id}
// Retrieve specified network interface
func DeleteNetworkInterface(vpcService *vpcv1.VpcV1, instanceID, vnicID string) (response *core.DetailedResponse, err error) {
	options := &vpcv1.DeleteInstanceNetworkInterfaceOptions{}
	options.SetID(vnicID)
	options.SetInstanceID(instanceID)
	response, err = vpcService.DeleteInstanceNetworkInterface(options)
	return response, err
}

// GetNetworkInterface GET
// /instances/{instance_id}/network_interfaces/{id}
// Retrieve specified network interface
func GetNetworkInterface(vpcService *vpcv1.VpcV1, instanceID, networkID string) (networkInterface *vpcv1.NetworkInterface, response *core.DetailedResponse, err error) {
	options := &vpcv1.GetInstanceNetworkInterfaceOptions{}
	options.SetID(networkID)
	options.SetInstanceID(instanceID)
	networkInterface, response, err = vpcService.GetInstanceNetworkInterface(options)
	return
}

// UpdateNetworkInterface PATCH
// /instances/{instance_id}/network_interfaces/{id}
// Update a network interface
func UpdateNetworkInterface(vpcService *vpcv1.VpcV1, instanceID, networkID, name string) (networkInterface *vpcv1.NetworkInterface, response *core.DetailedResponse, err error) {
	body := &vpcv1.NetworkInterfacePatch{
		Name: &name,
	}
	patchBody, _ := body.AsPatch()
	options := &vpcv1.UpdateInstanceNetworkInterfaceOptions{
		NetworkInterfacePatch: patchBody,
		ID:                    &networkID,
		InstanceID:            &instanceID,
	}
	networkInterface, response, err = vpcService.UpdateInstanceNetworkInterface(options)
	return
}

// ListNetworkInterfaceFloatingIPs GET
// /instances/{instance_id}/network_interfaces
// List all network interfaces on an instance
func ListNetworkInterfaceFloatingIPs(vpcService *vpcv1.VpcV1, instanceID, networkID string) (fips *vpcv1.FloatingIPUnpaginatedCollection, response *core.DetailedResponse, err error) {
	options := &vpcv1.ListInstanceNetworkInterfaceFloatingIpsOptions{}
	options.SetInstanceID(instanceID)
	options.SetNetworkInterfaceID(networkID)
	fips, response, err = vpcService.ListInstanceNetworkInterfaceFloatingIps(options)
	return
}

// GetNetworkInterfaceFloatingIP GET
// /instances/{instance_id}/network_interfaces/{network_interface_id}/floating_ips
// List all floating IPs associated with a network interface
func GetNetworkInterfaceFloatingIP(vpcService *vpcv1.VpcV1, instanceID, networkID, fipID string) (fip *vpcv1.FloatingIP, response *core.DetailedResponse, err error) {
	options := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{}
	options.SetID(fipID)
	options.SetInstanceID(instanceID)
	options.SetNetworkInterfaceID(networkID)
	fip, response, err = vpcService.GetInstanceNetworkInterfaceFloatingIP(options)
	return
}

// DeleteNetworkInterfaceFloatingIPBinding DELETE
// /instances/{instance_id}/network_interfaces/{network_interface_id}/floating_ips/{id}
// Disassociate specified floating IP
func DeleteNetworkInterfaceFloatingIPBinding(vpcService *vpcv1.VpcV1, instanceID, networkID, fipID string) (response *core.DetailedResponse, err error) {
	options := &vpcv1.RemoveInstanceNetworkInterfaceFloatingIPOptions{}
	options.SetID(fipID)
	options.SetInstanceID(instanceID)
	options.SetNetworkInterfaceID(networkID)
	response, err = vpcService.RemoveInstanceNetworkInterfaceFloatingIP(options)
	return response, err
}

// CreateNetworkInterfaceFloatingIPBinding PUT
// /instances/{instance_id}/network_interfaces/{network_interface_id}/floating_ips/{id}
// Associate a floating IP with a network interface
func CreateNetworkInterfaceFloatingIPBinding(vpcService *vpcv1.VpcV1, instanceID, networkID, fipID string) (fip *vpcv1.FloatingIP, response *core.DetailedResponse, err error) {
	options := &vpcv1.AddInstanceNetworkInterfaceFloatingIPOptions{}
	options.SetID(fipID)
	options.SetInstanceID(instanceID)
	options.SetNetworkInterfaceID(networkID)
	fip, response, err = vpcService.AddInstanceNetworkInterfaceFloatingIP(options)
	return
}

/**
 * Subnets
 *
 */

// ListSubnets - GET
// /subnets
// List all subnets
func ListSubnets(vpcService *vpcv1.VpcV1) (subnets *vpcv1.SubnetCollection, response *core.DetailedResponse, err error) {
	options := &vpcv1.ListSubnetsOptions{}
	subnets, response, err = vpcService.ListSubnets(options)
	return
}

// GetSubnet - GET
// /subnets/{id}
// Retrieve specified subnet
func GetSubnet(vpcService *vpcv1.VpcV1, subnetID string) (subnet *vpcv1.Subnet, response *core.DetailedResponse, err error) {
	options := &vpcv1.GetSubnetOptions{}
	options.SetID(subnetID)
	subnet, response, err = vpcService.GetSubnet(options)
	return
}

// DeleteSubnet - DELETE
// /subnets/{id}
// Delete specified subnet
func DeleteSubnet(vpcService *vpcv1.VpcV1, id string) (response *core.DetailedResponse, err error) {
	options := &vpcv1.DeleteSubnetOptions{}
	options.SetID(id)
	response, err = vpcService.DeleteSubnet(options)
	return response, err
}

// UpdateSubnet - PATCH
// /subnets/{id}
// Update specified subnet
func UpdateSubnet(vpcService *vpcv1.VpcV1, id, name string) (subnet *vpcv1.Subnet, response *core.DetailedResponse, err error) {
	body := &vpcv1.SubnetPatch{
		Name: &name,
	}
	patchBody, _ := body.AsPatch()
	options := &vpcv1.UpdateSubnetOptions{
		SubnetPatch: patchBody,
	}
	options.SetID(id)
	subnet, response, err = vpcService.UpdateSubnet(options)
	return
}

// CreateSubnet - POST
// /subnets
// Create a subnet
func CreateSubnet(vpcService *vpcv1.VpcV1, vpcID, name, zone string, mock bool) (subnet *vpcv1.Subnet, response *core.DetailedResponse, err error) {
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
	subnet, response, err = vpcService.CreateSubnet(options)
	return
}

// GetSubnetNetworkACL -GET
// /subnets/{id}/network_acl
// Retrieve a subnet's attached network ACL
func GetSubnetNetworkACL(vpcService *vpcv1.VpcV1, subnetID string) (subnetACL *vpcv1.NetworkACL, response *core.DetailedResponse, err error) {
	options := &vpcv1.GetSubnetNetworkACLOptions{}
	options.SetID(subnetID)
	subnetACL, response, err = vpcService.GetSubnetNetworkACL(options)
	return
}

// SetSubnetNetworkACLBinding - PUT
// /subnets/{id}/network_acl
// Attach a network ACL to a subnet
func SetSubnetNetworkACLBinding(vpcService *vpcv1.VpcV1, subnetID, id string) (nacl *vpcv1.NetworkACL, response *core.DetailedResponse, err error) {
	options := &vpcv1.ReplaceSubnetNetworkACLOptions{}
	options.SetID(subnetID)
	options.SetNetworkACLIdentity(&vpcv1.NetworkACLIdentity{ID: &id})
	nacl, response, err = vpcService.ReplaceSubnetNetworkACL(options)
	return
}

// DeleteSubnetPublicGatewayBinding - DELETE
// /subnets/{id}/public_gateway
// Detach a public gateway from a subnet
func DeleteSubnetPublicGatewayBinding(vpcService *vpcv1.VpcV1, id string) (response *core.DetailedResponse, err error) {
	options := &vpcv1.UnsetSubnetPublicGatewayOptions{}
	options.SetID(id)
	response, err = vpcService.UnsetSubnetPublicGateway(options)
	return response, err
}

// GetSubnetPublicGateway - GET
// /subnets/{id}/public_gateway
// Retrieve a subnet's attached public gateway
func GetSubnetPublicGateway(vpcService *vpcv1.VpcV1, subnetID string) (pgw *vpcv1.PublicGateway, response *core.DetailedResponse, err error) {
	options := &vpcv1.GetSubnetPublicGatewayOptions{}
	options.SetID(subnetID)
	pgw, response, err = vpcService.GetSubnetPublicGateway(options)
	return
}

// CreateSubnetPublicGatewayBinding - PUT
// /subnets/{id}/public_gateway
// Attach a public gateway to a subnet
func CreateSubnetPublicGatewayBinding(vpcService *vpcv1.VpcV1, subnetID, id string) (pgw *vpcv1.PublicGateway, response *core.DetailedResponse, err error) {
	options := &vpcv1.SetSubnetPublicGatewayOptions{}
	options.SetID(subnetID)
	options.SetPublicGatewayIdentity(&vpcv1.PublicGatewayIdentity{ID: &id})
	pgw, response, err = vpcService.SetSubnetPublicGateway(options)
	return
}

func main() {

	// Create a VPC Service
	vpcService := InstantiateVPCService()

	// List VPCs
	listVpcsOptions := &vpcv1.ListVpcsOptions{}
	vpcs, _, err := vpcService.ListVpcs(listVpcsOptions)
	if err != nil {
		log.Printf("Error from list call\n")
		panic(err)
	}
	//fmt.Printf("%T\n", vpcs)
	log.Printf("Num VPCs: %#v\n", *vpcs.TotalCount)

	// Create a new VPC
	vpcNew, _, err := GetVPC(vpcService, "hack-test-vpc")
	if err != nil {
		log.Printf("VPC doesn't exist, creating one.")
		vpcNew, _, err = CreateVPC(vpcService, "hack-test-vpc", "fc456662944c44f58610fac9a65c1fd5")
		if err != nil {
			log.Printf("Error from Create VPC call\n")
			panic(err)
		}
	} else {
		log.Printf("VPC already existed, skipping all steps")
		log.Printf(*vpcNew.ID)
	}

	// Create a subnet
	// subnetCollection,  := ListSubnets(vpcService)
	vpcSubnet, _, err := CreateSubnet(vpcService, *vpcNew.ID, "hack-test-subnet", "us-south-3", false)
	if err != nil {
		log.Printf("Error from Create Subnet call\n")
		panic(err)
	}
	log.Printf(*vpcSubnet.ID)

	// Create a public gateway
	pgw, _, err := CreatePublicGateway(vpcService, "hack-test-pgw", *vpcNew.ID, "us-south-3")
	if err != nil {
		log.Printf("Error from Create Public Gateway call\n")
		panic(err)
	}

	// Bind Subnet to the Public Gateway
	_, _, err = CreateSubnetPublicGatewayBinding(vpcService, *vpcSubnet.ID, *pgw.ID)
	if err != nil {
		fmt.Printf("Error from Bind Subnet Public Gateway call\n")
		panic(err)
	}

	log.Printf(*pgw.ID)

	// Set SSH key access
	sshKey, _, err := CreateSSHKey(vpcService, "hacktest", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDJV+b6WKJ3jbM9YeoM6uynlOf2NweFlNeoW/2kltUnwcSFldpELd2XDZ7Gg4SK1l4tHIVKaY1EYvgFA8CWGHwMUOLm8eHvhHK2aLOV8wZ++VCKKv2A15gFPl/yPYAYeSybE53zXWl7PTllE3P7H7/x5zKhrXYyQwVWU61wBLI1EL7cwbrYAK9M0f7cnb0fWsa7Qstxv0O5HrYFNB1f9MtXhJ/YZh5HfHBMzlKqvaal0fztVbSkP6a1ueqsyyT4kvEbWw/KZBQMAob2p2EZ6nSMJ+MxYnDEHC1wZGSaej57KB2smJ1aOnPIIJFYMxBb3UCEHHaXkzWnopFAPx2scfNQ3DmnY3DGVlKgb3oaySb/f3fhD4WdCTpd3JliKKEShQwxF6ncJd9lS+wb3O4zP7OTEjltCf0qPp4KQ8Hnx8d4feCRrWlQInlEi9LtJIBWjg7Unb0nnquxngjPENbeUZL2pQuPvDcVxXLA+9TcjvV/E06xrj8hqZprP4bxyGSZXb0= thousandsunny@tp450-7245q-ibm-com")
	if err != nil {
		fmt.Printf("Error from Create SSH Key call\n")
		panic(err)
	}
	log.Printf(*sshKey.ID)

	// instances, _, err := ListInstances(vpcService)
	// log.Printf(*instances.)
	// VSI instance setup
	vsi, _, err := CreateInstance(vpcService, "hack-test-vsi", "bx2-8x32", "r134-be4e8fa4-5004-47a2-b47d-397566d08e14", "us-south-3", *vpcSubnet.ID, *sshKey.ID, *vpcNew.ID)
	if err != nil {
		fmt.Printf("Error from Create Instance call\n")
		panic(err)
	}
	log.Printf(*vsi.ID)

	// Create a Network Interface
	netInterface, _, err := ListNetworkInterfaces(vpcService, *vsi.ID)
	// netInterface, _, err := CreateNetworkInterface(vpcService, *vsi.ID, *vpcSubnet.ID)
	if err != nil {
		fmt.Printf("Error from Create Network Interface call\n")
		panic(err)
	}
	log.Printf(*netInterface.NetworkInterfaces[0].Name)

	// Floating IP setup
	fip, _, err := CreateFloatingIP(vpcService, "us-south-3", "hack-test-fip")
	if err != nil {
		fmt.Printf("Error from Create Floating IP call\n")
		panic(err)
	}
	log.Printf(*fip.ID)

	// Bind Floating IP to Network Interface
	_, _, err = CreateNetworkInterfaceFloatingIPBinding(vpcService, *vsi.ID, *netInterface.NetworkInterfaces[0].ID, *fip.ID)
	if err != nil {
		fmt.Printf("Error from Create Floating IP call\n")
		panic(err)
	}
	log.Printf(*fip.Address)
	log.Printf(*vsi.Status)

	// Send SSH key access and list

	// time.Sleep(120 * time.Second)

	// // Delete the ssh key
	// _, err = DeleteSSHKey(vpcService, "office-tp")

	// // Delete the Subnet Gateway binding
	// _, err = DeleteSubnetPublicGatewayBinding(vpcService, *pgw.ID)
	// // Delete the Public Gateway
	// _, err = DeletePublicGateway(vpcService, *pgw.ID)
	// // Delete a Subnet
	// _, err = DeleteSubnet(vpcService, *vpcSubnet.ID)

	// // Delete VPC
	// _, err = DeleteVPC(vpcService, *vpcNew.ID)

	// // Check for errors
	// if err != nil {
	// 	fmt.Printf("Error from Delete calls\n")
	// 	panic(err)
	// }
}
