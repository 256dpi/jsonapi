package jsonapi

import (
	"net"
	"net/http"
	"time"
)

func withServer(cb func(base string, server *Server)) {
	serverConfig := ServerConfig{}

	srv := NewServer(serverConfig)

	lst, err := net.Listen("tcp", "0.0.0.0:1337")
	if err != nil {
		panic(err)
	}

	s := &http.Server{Handler: srv}
	go s.Serve(lst)

	cb("http://0.0.0.0:1337", srv)

	_ = s.Close()
	_ = lst.Close()

	time.Sleep(time.Millisecond)
}
