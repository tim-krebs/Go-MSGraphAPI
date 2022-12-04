package graphhelper

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type GraphHelper struct {
	deviceCodeCredential   *azidentity.DeviceCodeCredential
	userClient             *msgraphsdk.GraphServiceClient
	graphUserScopes        []string
	clientSecretCredential *azidentity.ClientSecretCredential
	appClient              *msgraphsdk.GraphServiceClient
}

func NewGraphHelper() *GraphHelper {
	g := &GraphHelper{}
	return g
}

func (g *GraphHelper) EnsureGraphForAppOnlyAuth() error {
	if g.clientSecretCredential == nil {
		clientId := os.Getenv("CLIENT_ID")
		tenantId := os.Getenv("TENANT_ID")
		clientSecret := os.Getenv("CLIENT_SECRET")
		credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
		if err != nil {
			return err
		}

		g.clientSecretCredential = credential
	}

	if g.appClient == nil {
		// Create an auth provider using the credential
		authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(g.clientSecretCredential, []string{
			"https://graph.microsoft.com/.default",
		})
		if err != nil {
			return err
		}

		// Create a request adapter using the auth provider
		adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
		if err != nil {
			return err
		}

		// Create a Graph client using request adapter
		client := msgraphsdk.NewGraphServiceClient(adapter)
		g.appClient = client
	}

	return nil
}

// ///////////////////////////////////////////////////////////////////////////////////
// List User in AAD
// ///////////////////////////////////////////////////////////////////////////////////
func (g *GraphHelper) GetUsers() (models.UserCollectionResponseable, error) {
	err := g.EnsureGraphForAppOnlyAuth()
	if err != nil {
		return nil, err
	}

	var topValue int32 = 25
	query := users.UsersRequestBuilderGetQueryParameters{
		// Only request specific properties
		Select: []string{"displayName", "id", "mail"},
		// Get at most 25 results
		Top: &topValue,
		// Sort by display name
		Orderby: []string{"displayName"},
	}

	return g.appClient.Users().
		Get(context.Background(),
			&users.UsersRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
}

// ///////////////////////////////////////////////////////////////////////////////////
// Create user in AAD
// ///////////////////////////////////////////////////////////////////////////////////
func (g *GraphHelper) CreateUsers() error {
	err := g.EnsureGraphForAppOnlyAuth()
	if err != nil {
		return err
	}

	requestBody := models.NewUser()
	accountEnabled := true
	requestBody.SetAccountEnabled(&accountEnabled)
	displayName := "Melissa Darrow"
	requestBody.SetDisplayName(&displayName)
	mailNickname := "MelissaD"
	requestBody.SetMailNickname(&mailNickname)
	userPrincipalName := "MelissaD@timkrebs9outlook.onmicrosoft.com"
	requestBody.SetUserPrincipalName(&userPrincipalName)
	passwordProfile := models.NewPasswordProfile()
	forceChangePasswordNextSignIn := true
	passwordProfile.SetForceChangePasswordNextSignIn(&forceChangePasswordNextSignIn)
	password := "xWwvJ]6NMw+bWH-d"
	passwordProfile.SetPassword(&password)
	requestBody.SetPasswordProfile(passwordProfile)

	result, err := g.appClient.Users().Post(context.Background(), requestBody, nil)

	fmt.Println(result)

	return err
}

// ///////////////////////////////////////////////////////////////////////////////////
// Delete user in AAD
// ///////////////////////////////////////////////////////////////////////////////////
func (g *GraphHelper) DeleteUser(username string) string {
	var succDelete string = "User deleted succesfully"
	users, err := g.GetUsers()
	if err != nil {
		log.Panicf("Error getting users: %v", err)
	}

	// Get user ID
	for _, user := range users.GetValue() {
		if *user.GetDisplayName() == username {
			g.appClient.UsersById(*user.GetId()).Delete(context.Background(), nil)
			return succDelete
		}
	}
	return "User deletion unsuccesful"
}

// ///////////////////////////////////////////////////////////////////////////////////
// Update user in AAD
// ///////////////////////////////////////////////////////////////////////////////////
func (g *GraphHelper) UpdateUser(username string) error {

	var id string

	err := g.EnsureGraphForAppOnlyAuth()
	if err != nil {
		return err
	}

	// Get User ID
	users, err := g.GetUsers()
	if err != nil {
		log.Panicf("Error getting users: %v", err)
	}
	for _, user := range users.GetValue() {
		if *user.GetDisplayName() == username {
			id = *user.GetId()
		}
	}

	// Create new Body
	requestBody := models.NewUser()
	businessPhones := []string{"+1 425 555 0109"}
	requestBody.SetBusinessPhones(businessPhones)
	officeLocation := "18/2111"
	requestBody.SetOfficeLocation(&officeLocation)

	result, err := g.appClient.UsersById(id).Patch(context.Background(), requestBody, nil)
	if result != nil {
		fmt.Println(result)
	}

	return err
}
