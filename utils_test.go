package jsonapi

import (
	"errors"
	"net"
	"net/http"
)

func withServer(cb func(client *Client, server *Server)) {
	server := NewServer(ServerConfig{})

	socket, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}

	go func() {
		err = http.Serve(socket, server)
		if !errors.Is(err, net.ErrClosed) {
			panic(err)
		}
	}()

	client := NewClient(ClientConfig{
		BaseURI: "http://" + socket.Addr().String(),
	})

	cb(client, server)

	err = socket.Close()
	if err != nil {
		panic(err)
	}
}
