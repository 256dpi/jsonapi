package adapter

import (
	"github.com/gonfire/jsonapi"
	"github.com/labstack/echo/engine"
)

type requestBridge struct {
	request engine.Request
}

func (b *requestBridge) Method() string {
	return b.request.Method()
}

func (b *requestBridge) Get(key string) string {
	return b.request.Header().Get(key)
}

func (b *requestBridge) Path() string {
	return b.request.URL().Path()
}

func (b *requestBridge) QueryParams() map[string][]string {
	return b.request.URL().QueryParams()
}

type responseBridge struct {
	response engine.Response
}

func (b *responseBridge) Set(key, value string) {
	b.response.Header().Set(key, value)
}

func (b *responseBridge) WriteHeader(status int) {
	b.response.WriteHeader(status)
}

func (b *responseBridge) Write(p []byte) (int, error) {
	return b.response.Write(p)
}

// BridgeRequest will return a bridge for the passed echo request to be
// compatible with the jsonapi.Requester interface.
func BridgeRequest(r engine.Request) jsonapi.Requester {
	return &requestBridge{
		request: r,
	}
}

// BridgeResponse will return a bridge for the passed echo response to be
// compatible with the jsonapi.Responder interface.
func BridgeResponse(w engine.Response) jsonapi.Responder {
	return &responseBridge{
		response: w,
	}
}
