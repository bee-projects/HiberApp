package azure

import (
	"context"
	"log"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/writeameer/aci/helpers"
)

var (
	ctx = context.Background()
)

// Auth Checks creds are provided in the ENV and returns an Azure token and Subscription ID
func Auth() (authorizer autorest.Authorizer, sid string) {
	// Check env for creds and read env
	helpers.FatalError(helpers.CheckEnv())
	sid = os.Getenv("AZURE_SUBSCRIPTION_ID")

	// Authenticate with Azure
	log.Println("Starting azure auth...")
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	helpers.FatalError(err)
	log.Println("After azure auth...")

	return
}
