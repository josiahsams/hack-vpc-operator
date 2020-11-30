package vsi

import (
	"fmt"
	bluemix "github.com/IBM-Cloud/bluemix-go"
	resMgr "github.com/IBM-Cloud/bluemix-go/api/resource/resourcev2/managementv2"
	"github.com/IBM-Cloud/bluemix-go/session"
)

func getResourceGroupID(apiKey string) (string, error) {
	fmt.Println("Get Token .. ")

	config := bluemix.Config{}
	config.BluemixAPIKey = apiKey
	tokenendpoint := tokenURL
	endpoint := resourceCtrlURL
	config.TokenProviderEndpoint = &tokenendpoint
	config.Endpoint = &endpoint
	s, err := session.New(&config)
	if err != nil {
		fmt.Println("1", err)
		return "", err
	}

	rMgr, err := resMgr.New(s)

	if err != nil {
		fmt.Println("2", err)
		return "", err
	}

	rGrp := rMgr.ResourceGroup()

	query := resMgr.ResourceGroupQuery{Default: true}
	rList, err := rGrp.List(&query)

	if err != nil {
		fmt.Println("3", err)
		return "", err
	}

	resourceGroupID := rList[0].ID
	fmt.Println("resourceGroupID: ", resourceGroupID)
	return resourceGroupID, err
}
