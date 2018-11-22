package azure

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/writeameer/aci/helpers"
)

// DeployArmTemplate Deploys and ARM template
func DeployArmTemplate(groupName string, location string, deploymentName string, template *map[string]interface{}, paramaters *map[string]interface{}) (deployment resources.DeploymentExtended, err error) {

	// Authenticate with Azure
	authorizer, sid := Auth()

	// Setup ARM Client
	armClient := resources.NewGroupsClient(sid)
	armClient.Authorizer = authorizer

	// Create ARM group
	params := resources.Group{
		Location: &location,
	}
	group, err := armClient.CreateOrUpdate(ctx, groupName, params)
	helpers.PrintError(err)
	log.Printf("%v arm group created", *group.Name)

	// Create deployment client
	dClient := resources.NewDeploymentsClient(sid)
	dClient.Authorizer = authorizer

	// Deploy ARM template deployment
	deploymentFuture, err := dClient.CreateOrUpdate(
		ctx,
		groupName,
		deploymentName,
		resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: paramaters,
				Mode:       resources.Incremental,
			},
		},
	)

	helpers.PrintError(err)

	// Wait for completion
	log.Printf("Wait for completion...")
	err = deploymentFuture.Future.WaitForCompletion(ctx, dClient.BaseClient.Client)
	if err != nil {
		return
	}
	log.Printf("Deployment completed...")
	return deploymentFuture.Result(dClient)
}
