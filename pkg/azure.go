package azure

import "context"

type Client interface {
	AzureCreateNewUser(ctx context.Context, email, password, firstname, lastname string)
	AzureAddExtensionToUser(ctx context.Context, sessionkey string)
}
