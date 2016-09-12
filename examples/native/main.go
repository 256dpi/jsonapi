// This example implements a basic API using the standard HTTP package.
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gonfire/jsonapi"
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

func entryPoint(writer http.ResponseWriter, r *http.Request) {
	w := jsonapi.BridgeResponseWriter(writer)

	req, err := jsonapi.ParseRequest(jsonapi.BridgeRequest(r), "/api/")
	if err != nil {
		jsonapi.WriteError(w, err)
		return
	}

	if req.ResourceType != "posts" {
		jsonapi.WriteError(w, jsonapi.NotFound("The requested resource is not available"))
		return
	}

	var doc *jsonapi.Document
	if req.Intent.DocumentExpected() {
		doc, err = jsonapi.ParseDocument(r.Body)
		if err != nil {
			jsonapi.WriteError(w, err)
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
		jsonapi.WriteError(w, err)
	}
}

func listPosts(req *jsonapi.Request, w jsonapi.Responder) error {
	list := make([]*jsonapi.Resource, 0, len(store))
	for _, post := range store {
		list = append(list, &jsonapi.Resource{
			Type:       "posts",
			ID:         post.ID,
			Attributes: jsonapi.StructToMap(post, req.Fields["posts"]),
		})
	}

	return jsonapi.WriteResources(w, http.StatusOK, list, &jsonapi.DocumentLinks{
		Self: "/api/posts",
	})
}

func findPost(req *jsonapi.Request, w jsonapi.Responder) error {
	post, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	return writePost(req, w, http.StatusOK, post)
}

func createPost(req *jsonapi.Request, doc *jsonapi.Document, w jsonapi.Responder) error {
	post := &postModel{
		ID: strconv.Itoa(counter),
	}

	err := doc.Data.One.AssignAttributes(post)
	if err != nil {
		return err
	}

	counter++
	store[post.ID] = post

	return writePost(req, w, http.StatusCreated, post)
}

func updatePost(req *jsonapi.Request, doc *jsonapi.Document, w jsonapi.Responder) error {
	post, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	err := doc.Data.One.AssignAttributes(post)
	if err != nil {
		return err
	}

	return writePost(req, w, http.StatusOK, post)
}

func deletePost(req *jsonapi.Request, w jsonapi.Responder) error {
	_, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	delete(store, req.ResourceID)

	w.WriteHeader(http.StatusOK)
	return nil
}

func writePost(req *jsonapi.Request, w jsonapi.Responder, status int, post *postModel) error {
	return jsonapi.WriteResource(w, status, &jsonapi.Resource{
		Type:       "posts",
		ID:         post.ID,
		Attributes: jsonapi.StructToMap(post, req.Fields["posts"]),
	}, &jsonapi.DocumentLinks{
		Self: "/api/posts/" + post.ID,
	})
}
