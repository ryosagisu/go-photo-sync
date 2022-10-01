package google_photos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/photoslibrary/v1"
)

// GetService Retrieve a token, saves the token, then returns the generated client.
func GetService(credentialPath string) *photoslibrary.Service {
	credentialFile := fmt.Sprintf("%s/credential.json", credentialPath)
	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	oauthConfig, err := google.ConfigFromJSON(b, photoslibrary.PhotoslibraryScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := fmt.Sprintf("%s/token.json", credentialPath)
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(oauthConfig)
		saveToken(tokFile, tok)
	}

	client := oauthConfig.Client(context.Background(), tok)
	service, err := photoslibrary.New(client)
	if err != nil {
		log.Fatal(err)
	}

	return service
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the \n%v\n", authURL)
	fmt.Println("Type authorization code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()

	_ = json.NewEncoder(f).Encode(token)

	if err := f.Close(); err != nil {
		log.Fatalf("Unable to save oauth token: %v", err)
	}
}
