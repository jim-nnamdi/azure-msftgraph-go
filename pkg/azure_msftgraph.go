package azure

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/jim-nnamdi/azure-msftgraph-go/pkg/web"
	a "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"go.uber.org/zap"
)

var _ Client = &azuremodel{}

var (
	ErrCouldNotGenerateToken = "could not generate token successfully"
	ErrGraphClient           = "could not initialize microsoft graph client successfully"
	ErrParsingJson           = "could not parse json successfully"
	ErrDeserializing         = "could not deserialize struct data"
)

type azuremodel struct {
	host          string
	logger        *zap.Logger
	clientID      string
	clientSecret  string
	loginClientID string
	graphUrl      string
	tenantID      string
	tenantUrl     string
	authorityUrl  string
	webClient     web.WebClient
}

func NewAzureModel(host string, clientID string, clientSecret string, loginClientID string, graphUrl string, tenantID string, authorityUrl string, tenantUrl string, logger *zap.Logger) *azuremodel {
	return &azuremodel{
		host:          host,
		logger:        logger,
		clientID:      clientID,
		clientSecret:  clientSecret,
		loginClientID: loginClientID,
		graphUrl:      graphUrl,
		tenantID:      tenantID,
		tenantUrl:     tenantUrl,
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
		log.Print(err.Error())
		model.logger.Debug(ErrCouldNotGenerateToken, zap.Any("credentials_secret", credentials_from_secret))
		return "", errors.New(ErrCouldNotGenerateToken)
	}
	app, err := confidential.New(model.clientID, credentials_from_secret, confidential.WithAuthority(model.authorityUrl))
	if err != nil {
		log.Print(err.Error())
		model.logger.Debug(ErrCouldNotGenerateToken, zap.Any("credentials_secret", credentials_from_secret))
		return "", errors.New(ErrCouldNotGenerateToken)
	}
	generate_token, err := app.AcquireTokenSilent(context.Background(), []string{model.graphUrl})
	if err != nil {
		generate_token, err := app.AcquireTokenByCredential(context.Background(), []string{model.graphUrl})
		if err != nil {
			log.Print(err.Error())
			model.logger.Debug(ErrCouldNotGenerateToken, zap.Any("acquire_token", credentials_from_secret))
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

func (model *azuremodel) AzureCreateNewUser(ctx context.Context, email, password, firstname, lastname string) (*AzUser, error) {
	requested_information := map[string]interface{}{
		"accountEnabled":   true,
		"displayName":      firstname + " " + lastname,
		"passwordPolicies": "DisablePasswordExpiration",
		"passwordProfile": map[string]interface{}{
			"password":                      password,
			"forceChangePasswordNextSignIn": false,
		},
		"preferredLanguage": "en-US",
		"identities": []map[string]interface{}{
			{
				"signInType":       "emailAddress",
				"issuer":           model.tenantID,
				"issuerAssignedId": email,
			},
		},
		"userType":     "guest",
		"creationType": "LocalAccount",
	}
	convert_map_to_byte_data, err := json.Marshal(requested_information)
	if err != nil {
		log.Print("failed to parse json")
		return nil, errors.New(ErrParsingJson)
	}
	convert_byte_data_to_string := string(convert_map_to_byte_data)
	log.Print(convert_byte_data_to_string)

	access_token, err := model.GetTokenUsingClientCredentials()
	res, err := model.webClient.ResponseUrlWithData(ctx, model.tenantUrl, convert_byte_data_to_string, access_token, http.MethodPost)
	if err != nil {
		log.Print(err.Error())
		return nil, errors.New(ErrParsingJson)
	}
	defer res.Body.Close()
	generate_data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Print(err.Error())
	}
	var user_data AzUser
	error := json.Unmarshal(generate_data, &user_data)
	if error != nil {
		log.Print("could not deserialize data into struct")
		return nil, errors.New(ErrDeserializing)
	}
	return &user_data, nil
}

func (model *azuremodel) AzureAddExtensionToUser(ctx context.Context, sessionkey string) {}
