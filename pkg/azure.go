package azure

import (
	"context"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type Client interface {
	GetTokenUsingClientCredentials() (string, error)
	InitializeClient() (*msgraphsdk.GraphServiceClient, error)
	AzureCreateNewUser(ctx context.Context, email, password, firstname, lastname string) (*AzUser, error)
	AzureAddExtensionToUser(ctx context.Context, sessionkey string)
}
