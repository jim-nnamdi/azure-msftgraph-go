package azure

import (
	"context"
	"fmt"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

var _ Client = &azuremodel{}

type azuremodel struct {
	host          string
	clientID      string
	clientSecret  string
	loginClientID string
	graphUrl      string
}

func NewAzureModel(host string, clientID string, clientSecret string, loginClientID string, graphUrl string) *azuremodel {
	return &azuremodel{
		host:          host,
		clientID:      clientID,
		clientSecret:  clientSecret,
		loginClientID: loginClientID,
		graphUrl:      graphUrl,
	}
}

func (model *azuremodel) GetTokenUsingClientCredentials() string {
	creds, err := confidential.NewCredFromSecret(model.clientSecret)
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}
	app, err := confidential.New(model.clientID, creds, confidential.WithAuthority("https://login.microsoftonline.com/smscb2cgp.onmicrosoft.com"))
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}
	generate_token, err := app.AcquireTokenSilent(context.Background(), []string{model.graphUrl})
	if err != nil {
		generate_token, err := app.AcquireTokenByCredential(context.Background(), []string{model.graphUrl})
		if err != nil {
			fmt.Print(err.Error())
			return ""
		}
		return generate_token.AccessToken
	}
	return generate_token.AccessToken
}

func (model *azuremodel) AzureCreateNewUser(ctx context.Context, email, password, firstname, lastname string) {

}

func (model *azuremodel) AzureAddExtensionToUser(ctx context.Context, sessionkey string) {}
