// This example implements a basic API using the standard HTTP package.
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gonfire/jsonapi"
	"github.com/gonfire/jsonapi/compat"
)

var counter = 1
var store = make(map[string]*postModel)

type postModel struct {
	ID    string `json:"-"`
	Title string `json:"title"`
}

func main() {
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func(star time.Time) {
			fmt.Printf("%6s  %-13s  %s\n", r.Method, r.URL.Path, time.Since(start).String())
		}(start)

		entryPoint(w, r)
	})

	http.ListenAndServe("0.0.0.0:4000", nil)
}

func entryPoint(w http.ResponseWriter, r *http.Request) {
	req, err := compat.ParseRequest(r, "/api/")
	if err != nil {
		compat.WriteError(w, err)
		return
	}

	if req.ResourceType != "posts" {
		compat.WriteError(w, jsonapi.NotFound("The requested resource is not available"))
		return
	}

	var doc *jsonapi.Document
	if req.Intent.DocumentExpected() {
		doc, err = jsonapi.ParseBody(r.Body)
		if err != nil {
			compat.WriteError(w, err)
			return
		}
	}

	if req.Intent == jsonapi.ListResources {
		err = listPosts(req, w)
	} else if req.Intent == jsonapi.FindResource {
		err = findPost(req, w)
	} else if req.Intent == jsonapi.CreateResource {
		err = createPost(req, doc, w)
	} else if req.Intent == jsonapi.UpdateResource {
		err = updatePost(req, doc, w)
	} else if req.Intent == jsonapi.DeleteResource {
		err = deletePost(req, w)
	} else {
		err = jsonapi.BadRequest("The requested method is not available")
	}

	if err != nil {
		compat.WriteError(w, err)
	}
}

func listPosts(_ *jsonapi.Request, w http.ResponseWriter) error {
	list := make([]*jsonapi.Resource, 0, len(store))
	for _, post := range store {
		list = append(list, &jsonapi.Resource{
			Type:       "posts",
			ID:         post.ID,
			Attributes: post,
		})
	}

	return compat.WriteResources(w, http.StatusOK, list, &jsonapi.DocumentLinks{
		Self: "/api/posts",
	})
}

func findPost(req *jsonapi.Request, w http.ResponseWriter) error {
	post, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	return writePost(w, http.StatusOK, post)
}

func createPost(_ *jsonapi.Request, doc *jsonapi.Document, w http.ResponseWriter) error {
	post := &postModel{
		ID: strconv.Itoa(counter),
	}

	err := jsonapi.MapToStruct(doc.Data.One.Attributes, post)
	if err != nil {
		return err
	}

	counter++
	store[post.ID] = post

	return writePost(w, http.StatusCreated, post)
}

func updatePost(req *jsonapi.Request, doc *jsonapi.Document, w http.ResponseWriter) error {
	post, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	err := jsonapi.MapToStruct(doc.Data.One.Attributes, post)
	if err != nil {
		return err
	}

	return writePost(w, http.StatusOK, post)
}

func deletePost(req *jsonapi.Request, w http.ResponseWriter) error {
	_, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	delete(store, req.ResourceID)

	w.WriteHeader(http.StatusOK)
	return nil
}

func writePost(w http.ResponseWriter, status int, post *postModel) error {
	return compat.WriteResource(w, status, &jsonapi.Resource{
		Type:       "posts",
		ID:         post.ID,
		Attributes: post,
	}, &jsonapi.DocumentLinks{
		Self: "/api/posts/" + post.ID,
	})
}
