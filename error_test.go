package jsonapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	var err error = &Error{
		Status: http.StatusNotFound,
		Title:  "Resource Not Found",
		Detail: "The requested resource cannot be found",
	}
	assert.Equal(t, "Resource Not Found: The requested resource cannot be found", err.Error())
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, &Error{
		Status: http.StatusNotFound,
		Title:  "Resource Not Found",
		Detail: "The requested resource cannot be found",
	})

	res := rec.Result()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, ContentType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "404",
			"title": "Resource Not Found",
			"detail": "The requested resource cannot be found"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorMissingStatus(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, &Error{
		Title:  "Resource Not Found",
		Detail: "The requested resource cannot be found",
	})

	res := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, ContentType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "Resource Not Found",
			"detail": "The requested resource cannot be found"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorNonError(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, errors.New("invalid"))

	res := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, ContentType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "Internal Server Error"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorFromStatus(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteErrorFromStatus(rec, http.StatusBadRequest)

	res := rec.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, ContentType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "400",
			"title": "Bad Request"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorFromStatusInvalidStatus(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteErrorFromStatus(rec, 0)

	res := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, ContentType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "Internal Server Error"
		}]
	}`, rec.Body.String())
}
