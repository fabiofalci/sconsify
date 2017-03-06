package rpc

import (
	"fmt"
	"github.com/fabiofalci/sconsify/sconsify"
	"net"
	"net/http"
	"net/rpc"
)

type NoArgs struct {
}

type Server struct {
	publisher *sconsify.Publisher
}

func StartServer(p *sconsify.Publisher) {
	server := new(Server)
	server.publisher = p
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
	} else if command == "pause" {
		method = "PauseTrack"
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
	if err := client.Call("Server."+method, &NoArgs{}, &reply); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

func (t *Server) NextTrack(args *NoArgs, reply *string) error {
	t.publisher.NextPlay()
	return nil
}

func (t *Server) PlayPause(args *NoArgs, reply *string) error {
	t.publisher.PlayPauseToggle()
	return nil
}

func (t *Server) PauseTrack(args *NoArgs, reply *string) error {
	t.publisher.Pause()
	return nil
}

func (t *Server) ReplayTrack(args *NoArgs, reply *string) error {
	t.publisher.Replay()
	return nil
}
