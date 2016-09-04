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
	ID    string
	Title string
}

func main() {
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func(star time.Time) {
			fmt.Printf("%s %s %s\n", pad(r.Method, 7), pad(r.URL.Path, 15), time.Since(start).String())
		}(start)

		req, err := jsonapi.ParseRequest(r, "/api/")
		if err != nil {
			jsonapi.WriteError(w, err)
			return
		}

		if req.ResourceType != "posts" {
			jsonapi.WriteErrorNotFound(w, "The requested resource is not available")
			return
		}

		var doc *jsonapi.Document
		if req.Intent.DocumentExpected() {
			doc, err = jsonapi.ParseBody(r.Body)
			if err != nil {
				jsonapi.WriteError(w, err)
				return
			}
		}

		if req.Intent == jsonapi.ListResources {
			listPosts(req, w)
			return
		} else if req.Intent == jsonapi.FindResource {
			findPost(req, w)
			return
		} else if req.Intent == jsonapi.CreateResource {
			createPost(req, doc, w)
			return
		} else if req.Intent == jsonapi.UpdateResource {
			updatePost(req, doc, w)
			return
		} else if req.Intent == jsonapi.DeleteResource {
			deletePost(req, w)
			return
		}

		jsonapi.WriteErrorFromStatus(w, http.StatusBadRequest, "The requested method is not available")
	})

	http.ListenAndServe("0.0.0.0:4000", nil)
}

func listPosts(_ *jsonapi.Request, w http.ResponseWriter) {
	list := make([]*jsonapi.Resource, 0, len(store))
	for _, post := range store {
		list = append(list, &jsonapi.Resource{
			Type: "posts",
			ID:   post.ID,
			Attributes: jsonapi.Map{
				"title": post.Title,
			},
		})
	}

	jsonapi.WriteResources(w, http.StatusOK, list, &jsonapi.DocumentLinks{
		Self: "/api/posts",
	})
}

func findPost(req *jsonapi.Request, w http.ResponseWriter) {
	post, ok := store[req.ResourceID]
	if !ok {
		jsonapi.WriteErrorNotFound(w, "The requested resource does not exist")
		return
	}

	writePost(w, http.StatusOK, post)
}

func createPost(_ *jsonapi.Request, doc *jsonapi.Document, w http.ResponseWriter) {
	post := &postModel{
		ID:    strconv.Itoa(counter),
		Title: doc.Data.One.Attributes["title"].(string), // FIXME
	}

	counter++
	store[post.ID] = post

	writePost(w, http.StatusCreated, post)
}

func updatePost(req *jsonapi.Request, doc *jsonapi.Document, w http.ResponseWriter) {
	post, ok := store[req.ResourceID]
	if !ok {
		jsonapi.WriteErrorNotFound(w, "The requested resource does not exist")
		return
	}

	post.Title = doc.Data.One.Attributes["title"].(string) // FIXME

	writePost(w, http.StatusOK, post)
}

func deletePost(req *jsonapi.Request, w http.ResponseWriter) {
	_, ok := store[req.ResourceID]
	if !ok {
		jsonapi.WriteErrorNotFound(w, "The requested resource does not exist")
		return
	}

	delete(store, req.ResourceID)

	w.WriteHeader(http.StatusOK)
}

func writePost(w http.ResponseWriter, status int, post *postModel) {
	jsonapi.WriteResource(w, status, &jsonapi.Resource{
		Type: "posts",
		ID:   post.ID,
		Attributes: jsonapi.Map{
			"title": post.Title,
		},
	}, &jsonapi.DocumentLinks{
		Self: "/api/posts/" + post.ID,
	})
}

func pad(str string, n int) string {
	for {
		if len(str) < n {
			str += " "
		}

		if len(str) >= n {
			return str
		}
	}
}
