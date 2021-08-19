package jsonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TODO: Provide helpers to navigate paginated responses? Add an iterator for
//  navigating the resources.

// ClientConfig is used to configure a client.
type ClientConfig struct {
	BaseURI       string
	Authorizer    func(*http.Request)
	ResponseLimit int64
}

// Client is a low-level jsonapi client.
type Client struct {
	config ClientConfig
	client *http.Client
}

// NewClient will create and return a new client.
func NewClient(config ClientConfig) *Client {
	return NewClientWithClient(config, new(http.Client))
}

// NewClientWithClient will create and return a new client with the specified
// client.
func NewClientWithClient(config ClientConfig, client *http.Client) *Client {
	// cleanup config
	config.BaseURI = strings.TrimSuffix(config.BaseURI, "/")

	// set default response limit
	if config.ResponseLimit == 0 {
		config.ResponseLimit = 8192
	}

	return &Client{
		config: config,
		client: client,
	}
}

// List will list the specified resources. The additional requests are merged
// with the base request.
func (c *Client) List(typ string, reqs ...Request) (*Document, error) {
	return c.Do(Request{
		Intent:       ListResources,
		ResourceType: typ,
	}.Merge(reqs...), nil)
}

// Find will find the specified resource. The additional requests are merged
// // with the base request.
func (c *Client) Find(typ, id string, reqs ...Request) (*Document, error) {
	return c.Do(Request{
		Intent:       FindResource,
		ResourceType: typ,
		ResourceID:   id,
	}.Merge(reqs...), nil)
}

// Create will create the specified resource.
func (c *Client) Create(res *Resource) (*Document, error) {
	return c.Do(Request{
		Intent:       CreateResource,
		ResourceType: res.Type,
	}, &Document{
		Data: &HybridResource{
			One: res,
		},
	})
}

// Update will update the specified resource.
func (c *Client) Update(res *Resource) (*Document, error) {
	return c.Do(Request{
		Intent:       UpdateResource,
		ResourceType: res.Type,
		ResourceID:   res.ID,
	}, &Document{
		Data: &HybridResource{
			One: res,
		},
	})
}

// Delete will delete the specified resource.
func (c *Client) Delete(typ, id string) error {
	_, err := c.Do(Request{
		Intent:       DeleteResource,
		ResourceType: typ,
		ResourceID:   id,
	}, nil)
	return err
}

// Do will perform the specified request and return the result.
func (c *Client) Do(req Request, doc *Document) (*Document, error) {
	// check doc
	if req.Intent.DocumentExpected() && doc == nil {
		return nil, fmt.Errorf("missing document")
	}

	// prepare url
	url := c.config.BaseURI + req.Self()

	// prepare body
	var body io.Reader
	if doc != nil {
		data, err := json.Marshal(doc)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}

	// create request
	r, err := http.NewRequest(req.Intent.RequestMethod(), url, body)
	if err != nil {
		return nil, err
	}

	// set content type if body is set
	if body != nil {
		r.Header.Set("Content-Type", MediaType)
	}

	// authorize request if available
	if c.config.Authorizer != nil {
		c.config.Authorizer(r)
	}

	// perform request
	res, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	// ensure body is closed
	defer func() {
		_ = res.Body.Close()
	}()

	// allow other status codes for some requests
	switch req.Intent {
	case CreateResource, UpdateResource, DeleteResource:
		switch res.StatusCode {
		case http.StatusAccepted, http.StatusNoContent:
			return nil, nil
		}
	}

	// prepare decoder
	dec := json.NewDecoder(io.LimitReader(res.Body, c.config.ResponseLimit))
	dec.UseNumber()

	// decode response
	var response Document
	err = dec.Decode(&response)
	if err != nil {
		return nil, err
	}

	// check errors
	if len(response.Errors) > 0 {
		return &response, response.Errors[0]
	}

	// handle errors
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("missing error")
	}

	return &response, nil
}
