package vsi

import (
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"log"
)

// CreateSSHKey - POST
// /keys
// Create a key
func (s *Service) createSSHKey(name, publicKey string) (key *vpcv1.Key, response *core.DetailedResponse, err error) {
	options := &vpcv1.CreateKeyOptions{}
	options.SetName(name)
	options.SetPublicKey(publicKey)
	key, response, err = s.vpcService.CreateKey(options)
	return
}

func (s *Service) getOrCreateSSHKey(keyname, publicKey string) error {
	// List SSHKey
	listKeysOptions := &vpcv1.ListKeysOptions{}
	keys, _, err := s.vpcService.ListKeys(listKeysOptions)

	if err != nil {
		log.Printf("Error from list call\n")
		return err
	}

	for _, key := range keys.Keys {
		if *key.Name == keyname {
			s.sshKeyID = *key.ID
			return nil
		}
	}

	log.Printf("SSH Key doesn't exist, creating one.")

	// Create a subnet
	sshKey, _, err := s.createSSHKey(keyname, publicKey)
	if err != nil {
		log.Printf("Error from Create SSHKey call\n")
		return err
	}

	s.sshKeyID = *sshKey.ID

	return nil
}
