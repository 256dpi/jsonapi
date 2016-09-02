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

		defer func() {
			fmt.Printf("%s %s %s\n", pad(r.Method, 7), pad(r.URL.Path, 15), time.Since(start).String())
		}()

		req, err := jsonapi.ParseRequest(r, "/api/")
		if err != nil {
			jsonapi.WriteError(w, err)
			return
		}

		if req.ResourceType != "posts" {
			jsonapi.WriteErrorFromStatus(w, http.StatusNotFound)
			return
		}

		var doc *jsonapi.Document
		if req.DocumentExpected() {
			doc, err = jsonapi.ParseBody(r.Body)
			if err != nil {
				jsonapi.WriteError(w, err)
				return
			}
		}

		if req.Action == jsonapi.Fetch {
			if req.Target == jsonapi.ResourceCollection {
				fetchPostList(req, w)
				return
			} else if req.Target == jsonapi.SingleResource {
				fetchSinglePost(req, w)
				return
			}
		} else if req.Action == jsonapi.Create && req.Target == jsonapi.ResourceCollection {
			createPost(req, doc, w)
			return
		} else if req.Action == jsonapi.Update && req.Target == jsonapi.SingleResource {
			updatePost(req, doc, w)
			return
		} else if req.Action == jsonapi.Delete && req.Target == jsonapi.SingleResource {
			deletePost(req, w)
			return
		}

		jsonapi.WriteErrorFromStatus(w, http.StatusBadRequest)
	})

	http.ListenAndServe("0.0.0.0:4000", nil)
}

func fetchPostList(_ *jsonapi.Request, w http.ResponseWriter) {
	list := make([]*jsonapi.Resource, 0)
	for _, post := range store {
		list = append(list, &jsonapi.Resource{
			Type: "posts",
			ID:   post.ID,
			Attributes: jsonapi.Map{
				"title": post.Title,
			},
		})
	}

	jsonapi.WriteResponse(w, http.StatusOK, &jsonapi.Document{
		Data: &jsonapi.HybridResource{
			Many: list,
		},
	})
}

func fetchSinglePost(req *jsonapi.Request, w http.ResponseWriter) {
	post, ok := store[req.ResourceID]
	if !ok {
		jsonapi.WriteErrorFromStatus(w, http.StatusNotFound)
		return
	}

	writePost(w, post)
}

func createPost(_ *jsonapi.Request, doc *jsonapi.Document, w http.ResponseWriter) {
	post := &postModel{
		ID:    strconv.Itoa(counter),
		Title: doc.Data.One.Attributes["title"].(string),
	}

	counter++
	store[post.ID] = post

	writePost(w, post)
}

func updatePost(req *jsonapi.Request, doc *jsonapi.Document, w http.ResponseWriter) {
	post, ok := store[req.ResourceID]
	if !ok {
		jsonapi.WriteErrorFromStatus(w, http.StatusNotFound)
		return
	}

	post.Title = doc.Data.One.Attributes["title"].(string)

	writePost(w, post)
}

func deletePost(req *jsonapi.Request, w http.ResponseWriter) {
	_, ok := store[req.ResourceID]
	if !ok {
		jsonapi.WriteErrorFromStatus(w, http.StatusNotFound)
		return
	}

	delete(store, req.ResourceID)

	w.WriteHeader(http.StatusOK)
}

func writePost(w http.ResponseWriter, post *postModel) {
	jsonapi.WriteResponse(w, http.StatusOK, &jsonapi.Document{
		Data: &jsonapi.HybridResource{
			One: &jsonapi.Resource{
				Type: "posts",
				ID:   post.ID,
				Attributes: jsonapi.Map{
					"title": post.Title,
				},
			},
		},
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
