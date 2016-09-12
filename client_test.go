package jsonapi

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var lastRequest *http.Request
var nextResponse *http.Response

type testTransport struct{}

func (r *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	lastRequest = req
	return nextResponse, nil
}

func testClient() *Client {
	c := NewClient("")
	c.Client.Transport = &testTransport{}
	return c
}

func newTestResponse(status int, payload string) *http.Response {
	return &http.Response{
		StatusCode:    status,
		ContentLength: int64(len(payload)),
		Body:          ioutil.NopCloser(stringReader(payload)),
	}
}

func TestRequest(t *testing.T) {
	c := testClient()

	nextResponse = newTestResponse(http.StatusOK, `{
		"data": [
			{
				"type": "foo",
				"id": "1",
				"attributes": {
					"title": "bar"
				}
			}
		]
	}`)

	doc, err := c.Request(&Request{
		Intent:       ListResources,
		ResourceType: "foo",
	}, nil)
	assert.NoError(t, err)
	assert.Equal(t, &Document{
		Data: &HybridResource{
			Many: []*Resource{
				{
					Type: "foo",
					ID:   "1",
					Attributes: Map{
						"title": "bar",
					},
				},
			},
		},
	}, doc)

	assert.Equal(t, "GET", lastRequest.Method)
	assert.Equal(t, "/foo", lastRequest.URL.String())
	assert.Equal(t, "", readFullString(lastRequest.Body))
}

func TestRequestWithResource(t *testing.T) {
	c := testClient()

	nextResponse = newTestResponse(http.StatusOK, `{}`)

	doc, err := c.RequestWithResource(&Request{
		Intent:       CreateResource,
		ResourceType: "foo",
	}, &Resource{
		Type: "foo",
		ID:   "1",
		Attributes: Map{
			"title": "bar",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, &Document{}, doc)

	assert.Equal(t, "POST", lastRequest.Method)
	assert.Equal(t, "/foo", lastRequest.URL.String())
	assert.JSONEq(t, `{
		"data": {
			"type": "foo",
			"id": "1",
			"attributes": {
				"title": "bar"
			}
		}
	}`, readFullString(lastRequest.Body))
}

func TestRequestWithResources(t *testing.T) {
	c := testClient()

	nextResponse = newTestResponse(http.StatusOK, `{}`)

	doc, err := c.RequestWithResources(&Request{
		Intent:       AppendToRelationship,
		ResourceType: "foo",
		ResourceID:   "1",
		Relationship: "bar",
	}, []*Resource{
		{
			Type: "foo",
			ID:   "1",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, &Document{}, doc)

	assert.Equal(t, "POST", lastRequest.Method)
	assert.Equal(t, "/foo/1/relationships/bar", lastRequest.URL.String())
	assert.JSONEq(t, `{
		"data": [
			{
				"type": "foo",
				"id": "1"
			}
		]
	}`, readFullString(lastRequest.Body))
}

func TestClientRequestNoContent(t *testing.T) {
	c := testClient()

	nextResponse = newTestResponse(http.StatusNoContent, "")

	doc, err := c.Request(&Request{
		Intent:       ListResources,
		ResourceType: "foo",
	}, nil)
	assert.NoError(t, err)
	assert.Nil(t, doc)

	assert.Equal(t, "GET", lastRequest.Method)
	assert.Equal(t, "/foo", lastRequest.URL.String())
}

func TestClientRequestInvalidRequest(t *testing.T) {
	c := testClient()
	c.BaseURL = "%"

	doc, err := c.Request(&Request{}, nil)
	assert.Error(t, err)
	assert.Nil(t, doc)
}

func TestClientRequestUnderlyingError(t *testing.T) {
	c := testClient()
	c.Client.Transport = nil

	doc, err := c.Request(&Request{}, nil)
	assert.Error(t, err)
	assert.Nil(t, doc)
}
