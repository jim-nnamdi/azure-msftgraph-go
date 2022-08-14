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
	tenantID      string
	authorityUrl  string
}

func NewAzureModel(host string, clientID string, clientSecret string, loginClientID string, graphUrl string, tenantID string, authorityUrl string) *azuremodel {
	return &azuremodel{
		host:          host,
		clientID:      clientID,
		clientSecret:  clientSecret,
		loginClientID: loginClientID,
		graphUrl:      graphUrl,
		tenantID:      tenantID,
		authorityUrl:  authorityUrl,
	}
}

// this token is generated using the credentials given to the client
// which includes the client secret, client id and authority url
// the authority url includes the tenant id. this token would be used as
// authorization bearer for querying the microsoft graph.

func (model *azuremodel) GetTokenUsingClientCredentials() (string, error) {
	credentials_from_secret, err := confidential.NewCredFromSecret(model.clientSecret)
	if err != nil {
		fmt.Print(err.Error())
		return "", errors.New(ErrCouldNotGenerateToken)
	}
	app, err := confidential.New(model.clientID, credentials_from_secret, confidential.WithAuthority(model.authorityUrl))
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

// when the client is initialized, it can be used as a base
// which contains different objects and functions to actually
// make direct calls to the microsoft SDK.

func (model *azuremodel) InitializeClient() (*msgraphsdk.GraphServiceClient, error) {
	credential_from_azidentity, err := azidentity.NewClientSecretCredential(
		model.tenantID,
		model.clientID,
		model.clientSecret,
		nil,
	)
	if err != nil {
		log.Print("could not initialize azclient", err)
		return nil, errors.New(ErrGraphClient)
	}
	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(credential_from_azidentity, []string{model.graphUrl})
	if err != nil {
		log.Print("authentication err", err)
		return nil, errors.New(ErrGraphClient)
	}
	graph_adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		log.Print("adapter error", err)
		return nil, errors.New(ErrGraphClient)
	}
	initialized_client_result := msgraphsdk.NewGraphServiceClient(graph_adapter)
	return initialized_client_result, nil
}

func (model *azuremodel) AzureCreateNewUser(ctx context.Context, email, password, firstname, lastname string) {

}

func (model *azuremodel) AzureAddExtensionToUser(ctx context.Context, sessionkey string) {}
