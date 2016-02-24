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
)

func Auth(spotifyClientId string, authRedirectUrl string) *spotify.Client {
	if spotifyClientId == "" {
		fmt.Print("Spotify Client ID not set")
		return nil
	}

	auth := spotify.NewAuthenticator(authRedirectUrl,
		spotify.ScopeUserLibraryRead,
		spotify.ScopeUserFollowRead,
		spotify.ScopePlaylistReadCollaborative,
		spotify.ScopePlaylistReadPrivate)

	auth.SetAuthInfo(spotifyClientId, "")

	token := loadToken()
	if token == nil || hasExpired(token.Expiry) {
		url := auth.AuthURL("")
		url = strings.Replace(url, "response_type=code", "response_type=token", -1)
		fmt.Printf("For web api authorization go to url:\n\n%v\n\n", url)

		fmt.Print("And paste the access token here: ")
		reader := bufio.NewReader(os.Stdin)
		accessToken, _ := reader.ReadString('\n')
		accessToken = strings.Trim(accessToken, " \n\r")

		result := strings.Split(accessToken, " ")

		seconds, err := strconv.ParseInt(strings.Split(result[2], ":")[1], 10, 64)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return nil
		}
		expiry := time.Now().Add(time.Duration(seconds) * time.Second)
		token = &oauth2.Token{
			AccessToken: strings.Split(result[0], ":")[1],
			TokenType:   strings.Split(result[1], ":")[1],
			Expiry:      expiry,
		}
		persistToken(token)
	}

	client := auth.NewClient(token)
	return &client
}

func hasExpired(expiry time.Time) bool {
	return expiry.Before(time.Now())
}

func loadToken() *oauth2.Token {
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
