package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/fabiofalci/sconsify/ui"
	"github.com/howeyc/gopass"
)

func main() {
	username, pass := credentials()
	events := events.InitialiseEvents()

	go spotify.Initialise(username, pass, events)

	if <-events.Initialised {
		ui.Start(events)
	}
}

func credentials() (*string, *[]byte) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.Trim(username, " \n\r")
	fmt.Print("Password: ")
	pass := gopass.GetPasswd()
	return &username, &pass
}
