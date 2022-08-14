package azure

import (
	"context"
	"errors"
	"fmt"
	"log"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	a "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

var _ Client = &azuremodel{}

var (
	ErrCouldNotGenerateToken = "could not generate token successfully"
	ErrGraphClient           = "could not initialize microsoft graph client successfully"
)

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

func (model *azuremodel) GetTokenUsingClientCredentials() (string, error) {
	creds, err := confidential.NewCredFromSecret(model.clientSecret)
	if err != nil {
		fmt.Print(err.Error())
		return "", errors.New(ErrCouldNotGenerateToken)
	}
	app, err := confidential.New(model.clientID, creds, confidential.WithAuthority("https://login.microsoftonline.com/smscb2cgp.onmicrosoft.com"))
	if err != nil {
		fmt.Print(err.Error())
		return "", errors.New(ErrCouldNotGenerateToken)
	}
	generate_token, err := app.AcquireTokenSilent(context.Background(), []string{model.graphUrl})
	if err != nil {
		generate_token, err := app.AcquireTokenByCredential(context.Background(), []string{model.graphUrl})
		if err != nil {
			fmt.Print(err.Error())
			return "", errors.New(ErrCouldNotGenerateToken)
		}
		return generate_token.AccessToken, nil
	}
	return generate_token.AccessToken, nil
}

func (model *azuremodel) InitializeClient() (*msgraphsdk.GraphServiceClient, error) {
	cred, err := azidentity.NewClientSecretCredential(
		"smscb2cgp.onmicrosoft.com",
		model.clientID,
		model.clientSecret,
		nil,
	)
	if err != nil {
		log.Print("could not initialize azclient", err)
		return nil, errors.New(ErrGraphClient)
	}
	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		log.Print("authentication err", err)
		return nil, errors.New(ErrGraphClient)
	}
	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		log.Print("adapter error", err)
		return nil, errors.New(ErrGraphClient)
	}
	client := msgraphsdk.NewGraphServiceClient(adapter)
	return client, nil
}

func (model *azuremodel) AzureCreateNewUser(ctx context.Context, email, password, firstname, lastname string) {

}

func (model *azuremodel) AzureAddExtensionToUser(ctx context.Context, sessionkey string) {}
