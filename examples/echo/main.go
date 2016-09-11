// This example implements a basic API using the echo framework.
package main

import (
	"net/http"
	"strconv"

	"github.com/gonfire/jsonapi"
	"github.com/gonfire/jsonapi/adapter"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
)

var counter = 1
var store = make(map[string]*postModel)

type postModel struct {
	ID    string `json:"-"`
	Title string `json:"title"`
}

func main() {
	router := echo.New()
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${latency_human}\n",
	}))

	router.Use(entryPoint)

	router.GET("/api/posts", listPosts)
	router.GET("/api/posts/:id", findPost)
	router.POST("/api/posts", createPost)
	router.PATCH("/api/posts/:id", updatePost)
	router.DELETE("/api/posts/:id", deletePost)

	router.Run(fasthttp.New("0.0.0.0:4000"))
}

func entryPoint(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		r := adapter.BridgeRequest(ctx.Request())
		w := adapter.BridgeResponse(ctx.Response())

		req, err := jsonapi.ParseRequest(r, "/api/")
		if err != nil {
			return jsonapi.WriteError(w, err)
		}

		ctx.Set("req", req)

		var doc *jsonapi.Document
		if req.Intent.DocumentExpected() {
			doc, err = jsonapi.ParseDocument(ctx.Request().Body())
			if err != nil {
				return jsonapi.WriteError(w, err)
			}

			ctx.Set("doc", doc)
		}

		return next(ctx)
	}
}

func listPosts(ctx echo.Context) error {
	list := make([]*jsonapi.Resource, 0, len(store))
	for _, post := range store {
		list = append(list, &jsonapi.Resource{
			Type:       "posts",
			ID:         post.ID,
			Attributes: post,
		})
	}

	w := adapter.BridgeResponse(ctx.Response())

	return jsonapi.WriteResources(w, http.StatusOK, list, &jsonapi.DocumentLinks{
		Self: "/api/posts",
	})
}

func findPost(ctx echo.Context) error {
	req := ctx.Get("req").(*jsonapi.Request)

	post, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	return writePost(ctx, http.StatusOK, post)
}

func createPost(ctx echo.Context) error {
	doc := ctx.Get("doc").(*jsonapi.Document)

	post := &postModel{
		ID: strconv.Itoa(counter),
	}

	err := doc.Data.One.AssignAttributes(post)
	if err != nil {
		return err
	}

	counter++
	store[post.ID] = post

	return writePost(ctx, http.StatusCreated, post)
}

func updatePost(ctx echo.Context) error {
	req := ctx.Get("req").(*jsonapi.Request)
	doc := ctx.Get("doc").(*jsonapi.Document)

	post, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	err := doc.Data.One.AssignAttributes(post)
	if err != nil {
		return err
	}

	return writePost(ctx, http.StatusOK, post)
}

func deletePost(ctx echo.Context) error {
	req := ctx.Get("req").(*jsonapi.Request)

	_, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	delete(store, req.ResourceID)

	ctx.Response().WriteHeader(http.StatusOK)
	return nil
}

func writePost(ctx echo.Context, status int, post *postModel) error {
	w := adapter.BridgeResponse(ctx.Response())

	return jsonapi.WriteResource(w, status, &jsonapi.Resource{
		Type:       "posts",
		ID:         post.ID,
		Attributes: post,
	}, &jsonapi.DocumentLinks{
		Self: "/api/posts/" + post.ID,
	})
}
