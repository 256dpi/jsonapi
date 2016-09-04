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
	assert.Equal(t, MediaType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "404",
			"title": "Resource Not Found",
			"detail": "The requested resource cannot be found"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorEmpty(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, &Error{})

	res := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, MediaType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500"
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
	assert.Equal(t, MediaType, res.Header.Get("Content-Type"))
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
	assert.Equal(t, MediaType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "Internal Server Error",
			"detail": "Unknown error encountered"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorNil(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, nil)

	res := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, MediaType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "Internal Server Error",
			"detail": "Unknown error encountered"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListNone(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteErrorList(rec)

	res := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, MediaType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "Internal Server Error",
			"detail": "Unknown error encountered"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListInvalid(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteErrorList(rec, &Error{})

	res := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, MediaType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListSame(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteErrorList(rec, &Error{
		Status: http.StatusMethodNotAllowed,
	}, &Error{
		Status: http.StatusMethodNotAllowed,
	})

	res := rec.Result()

	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	assert.Equal(t, MediaType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "405"
		}, {
			"status": "405"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListSettleOn400(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteErrorList(rec, &Error{
		Status: http.StatusUnauthorized,
	}, &Error{
		Status: http.StatusForbidden,
	})

	res := rec.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, MediaType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "401"
		}, {
			"status": "403"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListSettleOn500(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteErrorList(rec, &Error{
		Status: http.StatusNotImplemented,
	}, &Error{
		Status: http.StatusBadGateway,
	})

	res := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, MediaType, res.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "501"
		}, {
			"status": "502"
		}]
	}`, rec.Body.String())
}

func TestErrorGenerators(t *testing.T) {
	list := []error{
		ErrorFromStatus(0, "foo"),
		NotFound("foo"),
		BadRequest("foo"),
		BadRequestParam("foo", "bar"),
	}

	for _, err := range list {
		assert.Error(t, err)
	}
}

func BenchmarkWriteError(b *testing.B) {
	err := &Error{
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := WriteError(httptest.NewRecorder(), err)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkWriteErrorList(b *testing.B) {
	err := &Error{
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := WriteErrorList(httptest.NewRecorder(), err, err)
		if err != nil {
			panic(err)
		}
	}
}
