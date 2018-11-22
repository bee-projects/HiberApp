package azure

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-07-01/storage"
	"github.com/writeameer/aci/helpers"
)

// CreateStorageAccount Creates an Azure storage account
func CreateStorageAccount(resourceGroupName string, storageAccountName string) {
	// Authenticate with Azure
	authorizer, sid := Auth()

	// Setup ARM Client

	client := storage.NewAccountsClient(sid)
	client.Authorizer = authorizer

	location := "Australia East"
	sku := storage.Sku{
		Name: storage.StandardLRS,
		Tier: storage.Standard,
	}
	kind := storage.StorageV2
	accountCreateFuture, err := client.Create(ctx, resourceGroupName, storageAccountName, storage.AccountCreateParameters{
		Location: &location,
		Kind:     kind,
		Sku:      &sku,
	})

	helpers.PrintError(err)

	err = accountCreateFuture.WaitForCompletion(ctx, client.BaseClient.Client)
	helpers.PrintError(err)

	account, err := accountCreateFuture.Result(client)

	log.Println(account.Name)
}
