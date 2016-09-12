package jsonapi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Client is a very basic client implementation that simplifies interacting with
// a JSON API compatible API.
type Client struct {
	base   string
	client *http.Client
}

// NewClient creates and returns a new client.
func NewClient(base string) *Client {
	return &Client{
		base:   base,
		client: &http.Client{},
	}
}

// Request will make a request using the specified request and an optional
// document.
func (c *Client) Request(req *Request, doc *Document) (*Document, error) {
	// prepare reader
	var buffer bytes.Buffer

	// make reader
	if doc != nil {
		json.NewEncoder(&buffer).Encode(doc)
	}

	// construct request
	r, err := http.NewRequest(req.Intent.RequestMethod(), c.base+req.Self(), &buffer)
	if err != nil {
		return nil, err
	}

	// set accept header
	r.Header.Set("Accept", MediaType)

	// set header if necessary
	if doc != nil {
		r.Header.Set("Content-Type", MediaType)
	}

	// do request
	res, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	// parse body if present
	if res.ContentLength > 0 {
		return ParseDocument(res.Body)
	}

	return nil, nil
}

// RequestWithResource will make a request using the specified request and the
// passed resource.
func (c *Client) RequestWithResource(req *Request, res *Resource) (*Document, error) {
	return c.Request(req, &Document{
		Data: &HybridResource{
			One: res,
		},
	})
}

// RequestWithResources will make a request using the specified request and the
// passed resources.
func (c *Client) RequestWithResources(req *Request, res []*Resource) (*Document, error) {
	return c.Request(req, &Document{
		Data: &HybridResource{
			Many: res,
		},
	})
}
