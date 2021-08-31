package jsonapi

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
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

func unescape(str string) string {
	str, err := url.QueryUnescape(str)
	if err != nil {
		panic(err)
	}
	return str
}

func escape(str string) string {
	str = strings.ReplaceAll(str, "[", "%5B")
	return strings.ReplaceAll(str, "]", "%5D")
}
