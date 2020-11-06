package jsonapi

import (
	"net"
	"net/http"
	"time"
)

func withServer(cb func(client *Client, server *Server)) {
	server := NewServer(ServerConfig{})

	lst, err := net.Listen("tcp", "0.0.0.0:1337")
	if err != nil {
		panic(err)
	}

	s := &http.Server{Handler: server}
	go s.Serve(lst)

	client := NewClient(ClientConfig{
		BaseURI: "http://0.0.0.0:1337",
	})

	cb(client, server)

	_ = s.Close()
	_ = lst.Close()

	time.Sleep(time.Millisecond)
}
