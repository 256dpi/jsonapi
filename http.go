package jsonapi

import "net/http"

type requestBridge struct {
	request *http.Request
}

func (b *requestBridge) Method() string {
	return b.request.Method
}

func (b *requestBridge) Get(key string) string {
	return b.request.Header.Get(key)
}

func (b *requestBridge) Path() string {
	return b.request.URL.Path
}

func (b *requestBridge) QueryParams() map[string][]string {
	return b.request.URL.Query()
}

type responseWriterBridge struct {
	responseWriter http.ResponseWriter
}

func (b *responseWriterBridge) Set(key, value string) {
	b.responseWriter.Header().Set(key, value)
}

func (b *responseWriterBridge) WriteHeader(status int) {
	b.responseWriter.WriteHeader(status)
}

func (b *responseWriterBridge) Write(p []byte) (int, error) {
	return b.responseWriter.Write(p)
}

// BridgeRequest will return a bridge for the passed http request to be compatible
// with the Requester interface.
func BridgeRequest(r *http.Request) Requester {
	return &requestBridge{
		request: r,
	}
}

// BridgeResponseWriter will return a bride for the passed http response writer
// to be compatible with the Responder interface.
func BridgeResponseWriter(w http.ResponseWriter) Responder {
	return &responseWriterBridge{
		responseWriter: w,
	}
}
