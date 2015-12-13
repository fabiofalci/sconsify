package webapi

import (
	"fmt"
	"github.com/zmb3/spotify"
	"os"
	"bufio"
	"strings"
	"golang.org/x/oauth2"
	"strconv"
	"time"
)


func Auth() *spotify.Client {
	auth := spotify.NewAuthenticator("https://fabiofalci.github.io/sconsify/auth/",
		spotify.ScopeUserLibraryRead, spotify.ScopeUserFollowRead, spotify.ScopePlaylistReadCollaborative, spotify.ScopePlaylistReadPrivate)

	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	if clientId == "" {
		fmt.Print("Spotify Client ID not set")
		return nil
	}

	auth.SetAuthInfo(clientId, "")

	url := auth.AuthURL("")
	url = strings.Replace(url, "response_type=code", "response_type=token", -1)
	fmt.Printf("Go to url %v\n", url)

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
	token := &oauth2.Token{
		AccessToken:  strings.Split(result[0], ":")[1],
		TokenType:    strings.Split(result[1], ":")[1],
		Expiry:       expiry,
	}

	client := auth.NewClient(token)
	return &client
}

