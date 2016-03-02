package rpc

import (
	"net/rpc"
	"net"
	"net/http"
	"fmt"
	"github.com/fabiofalci/sconsify/sconsify"
)

type NoArgs struct {
}

type Server struct {
	events *sconsify.Events
}

func StartServer(ev *sconsify.Events) {
	server := new(Server)
	server.events = ev
	rpc.Register(server)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":45800")
	if err != nil {
		fmt.Printf("Cannot start the server: %v\n", err)
		return
	}
	go http.Serve(listener, nil)
}

func Client(command string) {
	var method string
	if command == "next" {
		method = "NextTrack"
	} else if command == "play_pause" {
		method = "PlayPause"
	} else if command == "replay" {
		method = "ReplayTrack"
	} else {
		fmt.Println("Unknown command")
		return
	}

	client, err := rpc.DialHTTP("tcp", ":45800")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	var reply string
	if err := client.Call("Server." + method, &NoArgs{}, &reply); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

func (t *Server) NextTrack(args *NoArgs, reply *string) error {
	t.events.NextPlay()
	return nil
}

func (t *Server) PlayPause(args *NoArgs, reply *string) error {
	t.events.Pause()
	return nil
}

func (t *Server) ReplayTrack(args *NoArgs, reply *string) error {
	t.events.Replay()
	return nil
}
