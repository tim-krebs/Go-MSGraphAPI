package main

import (
	"fmt"
	"log"

	"GraphAPI/graphhelper"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/joho/godotenv"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type GraphHelper struct {
	deviceCodeCredential   *azidentity.DeviceCodeCredential
	userClient             *msgraphsdk.GraphServiceClient
	graphUserScopes        []string
	clientSecretCredential *azidentity.ClientSecretCredential
	appClient              *msgraphsdk.GraphServiceClient
}

func main() {
	// Load .env files
	// .env.local takes precedence (if present)
	godotenv.Load(".env.local")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	graphHelper := graphhelper.NewGraphHelper()

	var choice int64 = -1
	choice = 1

	for {
		fmt.Println()
		fmt.Println("Please choose one of the following options:")
		fmt.Println("0. Exit")
		fmt.Println("1. List all user")
		fmt.Println("2. Create a user")
		fmt.Println("3. Delete a specific user")
		fmt.Println("4. Update user credentials")
		fmt.Println()
		//_, err = fmt.Scanf("%d", &choice)
		//if err != nil {
		//	choice = -1
		//}

		switch choice {
		case 0:
			// Exit the program
			fmt.Println("Goodbye...")
		case 1:
			// Display access token
			listUsers(graphHelper)
		case 2:
			// Display access token
			createUsers(graphHelper)
		case 3:
			// Display access token
			deleteUser(graphHelper)
		case 4:
			// Display access token
			updateUser(graphHelper)

		default:
			fmt.Println("Invalid choice! Please try again.")
		}

		if choice == 0 {
			break
		}
	}
}

// ///////////////////////////////////////////////////////////////////////////////////
// List user in AAD
// ///////////////////////////////////////////////////////////////////////////////////
func listUsers(graphHelper *graphhelper.GraphHelper) {
	users, err := graphHelper.GetUsers()
	if err != nil {
		log.Panicf("Error getting users: %v", err)
	}

	// Output each user's details
	for _, user := range users.GetValue() {
		fmt.Printf("User: %s\n", *user.GetDisplayName())
		fmt.Printf("  ID: %s\n", *user.GetId())

		noEmail := "NO EMAIL"
		email := user.GetMail()
		if email == nil {
			email = &noEmail
		}
		fmt.Printf("  Email: %s\n", *email)
	}

	// If GetOdataNextLink does not return nil,
	// there are more users available on the server
	nextLink := users.GetOdataNextLink()

	fmt.Println()
	fmt.Printf("More users available? %t\n", nextLink != nil)
	fmt.Println()
}

// ///////////////////////////////////////////////////////////////////////////////////
// Create user in AAD
// ///////////////////////////////////////////////////////////////////////////////////
func createUsers(graphHelper *graphhelper.GraphHelper) {
	err := graphHelper.CreateUsers()
	if err != nil {
		fmt.Println(err)
	}
}

// ///////////////////////////////////////////////////////////////////////////////////
// Delete user in AAD
// ///////////////////////////////////////////////////////////////////////////////////
func deleteUser(graphHelper *graphhelper.GraphHelper) {
	var username string

	//fmt.Scanf("%s", &username)

	str := graphHelper.DeleteUser(username)
	fmt.Println(str)
}

// ///////////////////////////////////////////////////////////////////////////////////
// Update user in AAD
// ///////////////////////////////////////////////////////////////////////////////////
func updateUser(graphHelper *graphhelper.GraphHelper) {
	var username string

	//fmt.Scanf("%s", &username)

	err := graphHelper.UpdateUser(username)
	if err != nil {
		fmt.Println(err)
	}
}
