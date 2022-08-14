package azure

import "context"

var _ Client = &azuremodel{}

type azuremodel struct {
	host          string
	clientID      string
	clientSecret  string
	loginClientID string
}

func NewAzureModel(host string, clientID string, clientSecret string, loginClientID string) *azuremodel {
	return &azuremodel{
		host:          host,
		clientID:      clientID,
		clientSecret:  clientSecret,
		loginClientID: loginClientID,
	}
}

func (model *azuremodel) AzureCreateNewUser(ctx context.Context, email, password, firstname, lastname string) {
}

func (model *azuremodel) AzureAddExtensionToUser(ctx context.Context, sessionkey string) {}
