package jsonapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/engine/standard"
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
	res, rec := constructResponseAndRecorder()

	WriteError(res, &Error{
		Status: http.StatusNotFound,
		Title:  "Resource Not Found",
		Detail: "The requested resource cannot be found",
	})

	result := rec.Result()

	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "404",
			"title": "Resource Not Found",
			"detail": "The requested resource cannot be found"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorEmpty(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	WriteError(res, &Error{})

	result := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorMissingStatus(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	WriteError(res, &Error{
		Title:  "Resource Not Found",
		Detail: "The requested resource cannot be found",
	})

	result := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "Resource Not Found",
			"detail": "The requested resource cannot be found"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorNonError(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	WriteError(res, errors.New("invalid"))

	result := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "Internal Server Error"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorNil(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	WriteError(res, nil)

	result := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "Internal Server Error"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListNone(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	WriteErrorList(res)

	result := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500",
			"title": "Internal Server Error"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListInvalid(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	WriteErrorList(res, &Error{})

	result := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "500"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListSame(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	WriteErrorList(res, &Error{
		Status: http.StatusMethodNotAllowed,
	}, &Error{
		Status: http.StatusMethodNotAllowed,
	})

	result := rec.Result()

	assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "405"
		}, {
			"status": "405"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListSettleOn400(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	WriteErrorList(res, &Error{
		Status: http.StatusUnauthorized,
	}, &Error{
		Status: http.StatusForbidden,
	})

	result := rec.Result()

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "401"
		}, {
			"status": "403"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListSettleOn500(t *testing.T) {
	res, rec := constructResponseAndRecorder()

	WriteErrorList(res, &Error{
		Status: http.StatusNotImplemented,
	}, &Error{
		Status: http.StatusBadGateway,
	})

	result := rec.Result()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
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

func TestWriteErrorHTTP(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteErrorHTTP(rec, &Error{
		Status: http.StatusNotFound,
		Title:  "Resource Not Found",
		Detail: "The requested resource cannot be found",
	})

	result := rec.Result()

	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "404",
			"title": "Resource Not Found",
			"detail": "The requested resource cannot be found"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListHTTP(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteErrorListHTTP(rec, &Error{
		Status: http.StatusMethodNotAllowed,
	}, &Error{
		Status: http.StatusMethodNotAllowed,
	})

	result := rec.Result()

	assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
	assert.Equal(t, MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "405"
		}, {
			"status": "405"
		}]
	}`, rec.Body.String())
}

func BenchmarkWriteError(b *testing.B) {
	err := &Error{
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		res := standard.NewResponse(httptest.NewRecorder(), nil)

		err := WriteError(res, err)
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
		res := standard.NewResponse(httptest.NewRecorder(), nil)

		err := WriteErrorList(res, err, err)
		if err != nil {
			panic(err)
		}
	}
}
