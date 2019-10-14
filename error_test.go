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
		Title:  "resource not found",
		Detail: "the requested resource cannot be found",
	}
	assert.Equal(t, "resource not found: the requested resource cannot be found", err.Error())
}

func TestWriteError(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteError(res, &Error{
		Status: http.StatusNotFound,
		Title:  "resource not found",
		Detail: "the requested resource cannot be found",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	assert.Equal(t, MediaType, res.HeaderMap.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "404",
			"title": "resource not found",
			"detail": "the requested resource cannot be found"
		}]
	}`, res.Body.String())
}

func TestWriteErrorEmpty(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteError(res, &Error{})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Equal(t, MediaType, res.HeaderMap.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500"
		}]
	}`, res.Body.String())
}

func TestWriteErrorMissingStatus(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteError(res, &Error{
		Title:  "resource not found",
		Detail: "the requested resource cannot be found",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Equal(t, MediaType, res.HeaderMap.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "resource not found",
			"detail": "the requested resource cannot be found"
		}]
	}`, res.Body.String())
}

func TestWriteErrorNonError(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteError(res, errors.New("invalid"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Equal(t, MediaType, res.HeaderMap.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "internal server error"
		}]
	}`, res.Body.String())
}

func TestWriteErrorNil(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteError(res, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Equal(t, MediaType, res.HeaderMap.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "internal server error"
		}]
	}`, res.Body.String())
}

func TestWriteErrorListNone(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteErrorList(res)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Equal(t, MediaType, res.HeaderMap.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "internal server error"
		}]
	}`, res.Body.String())
}

func TestWriteErrorListInvalid(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteErrorList(res, &Error{})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Equal(t, MediaType, res.HeaderMap.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500"
		}]
	}`, res.Body.String())
}

func TestWriteErrorListSame(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteErrorList(res, &Error{
		Status: http.StatusMethodNotAllowed,
	}, &Error{
		Status: http.StatusMethodNotAllowed,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.Result().StatusCode)
	assert.Equal(t, MediaType, res.HeaderMap.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "405"
		}, {
			"status": "405"
		}]
	}`, res.Body.String())
}

func TestWriteErrorListSettleOn400(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteErrorList(res, &Error{
		Status: http.StatusUnauthorized,
	}, &Error{
		Status: http.StatusForbidden,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	assert.Equal(t, MediaType, res.HeaderMap.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "401"
		}, {
			"status": "403"
		}]
	}`, res.Body.String())
}

func TestWriteErrorListSettleOn500(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteErrorList(res, &Error{
		Status: http.StatusNotImplemented,
	}, &Error{
		Status: http.StatusBadGateway,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Equal(t, MediaType, res.HeaderMap.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "501"
		}, {
			"status": "502"
		}]
	}`, res.Body.String())
}

func TestErrorGenerators(t *testing.T) {
	list := []error{
		ErrorFromStatus(0, "foo"),
		NotFound("foo"),
		BadRequest("foo"),
		BadRequestParam("foo", "bar"),
		BadRequestPointer("foo", "bar"),
	}

	for _, err := range list {
		assert.Error(t, err)
	}
}

func BenchmarkWriteError(b *testing.B) {
	err := &Error{
		Title:  "internal server error",
		Status: http.StatusInternalServerError,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()

		err := WriteError(res, err)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkWriteErrorList(b *testing.B) {
	err := &Error{
		Title:  "internal server error",
		Status: http.StatusInternalServerError,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()

		err := WriteErrorList(res, err, err)
		if err != nil {
			panic(err)
		}
	}
}
