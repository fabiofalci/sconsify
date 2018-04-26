package webapi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fabiofalci/sconsify/infrastructure"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"os/exec"
)


var token *oauth2.Token
var auth spotify.Authenticator

func Auth(spotifyClientId string, authRedirectUrl string, cacheWebApiToken bool, openBrowserCommand string) (*spotify.Client, error) {
	if spotifyClientId == "" {
		fmt.Print("Spotify Client ID not set")
		return nil, nil
	}

	auth = spotify.NewAuthenticator(authRedirectUrl,
		spotify.ScopeUserLibraryRead,
		spotify.ScopeUserFollowRead,
		spotify.ScopePlaylistReadCollaborative,
		spotify.ScopePlaylistReadPrivate)

	auth.SetAuthInfo(spotifyClientId, "")

	LoadTokenFromFile()
	if token == nil || HasTokenExpired() {
		url := auth.AuthURL("")
		url = strings.Replace(url, "response_type=code", "response_type=token", -1)

		if openBrowserCommand != "" {
			fmt.Printf("\nOpen browser command provided: %v %v\n\n", openBrowserCommand, url)
			cmd := exec.Command(openBrowserCommand, url)
			if err := cmd.Run(); err != nil {
				return nil, err
			}
		} else {
			fmt.Printf("For web api authorization go to url:\n\n%v\n\n", url)
		}

		fmt.Print("And paste the access token here: ")
		reader := bufio.NewReader(os.Stdin)
		accessToken, _ := reader.ReadString('\n')
		accessToken = strings.Trim(accessToken, " \n\r")

		result := strings.Split(accessToken, " ")

		seconds, err := strconv.ParseInt(strings.Split(result[2], ":")[1], 10, 64)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return nil, err
		}
		expiry := time.Now().Add(time.Duration(seconds) * time.Second)
		token = &oauth2.Token{
			AccessToken: strings.Split(result[0], ":")[1],
			TokenType:   strings.Split(result[1], ":")[1],
			Expiry:      expiry,
		}
		if cacheWebApiToken {
			persistToken(token)
		}
	} else {
		fmt.Println("Token still valid")
	}

	client := NewClient()
	return &client, nil
}

func NewClient() spotify.Client {
	return auth.NewClient(token)
}

func HasTokenExpired() bool {
	return hasExpired(token.Expiry)
}

func hasExpired(expiry time.Time) bool {
	return expiry.Before(time.Now())
}

func LoadTokenFromFile() {
	token = loadTokenFromFile()
}

func loadTokenFromFile() *oauth2.Token {
	if fileLocation := infrastructure.GetWebApiTokenLocation(); fileLocation != "" {
		if b, err := ioutil.ReadFile(fileLocation); err == nil {
			var token *oauth2.Token
			if err := json.Unmarshal(b, &token); err == nil {
				return token
			}
		}
	}
	return nil
}

func persistToken(token *oauth2.Token) {
	if b, err := json.Marshal(token); err == nil {
		if fileLocation := infrastructure.GetWebApiTokenLocation(); fileLocation != "" {
			infrastructure.SaveFile(fileLocation, b)
		}
	}
}
