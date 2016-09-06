package main

import (
	"net/http"
	"strconv"

	"github.com/gonfire/jsonapi"
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

	router.Any("/api/*all", entryPoint)
	router.Run(fasthttp.New("0.0.0.0:4000"))
}

func entryPoint(ctx echo.Context) error {
	req, err := jsonapi.ParseRequest(ctx.Request(), "/api/")
	if err != nil {
		return jsonapi.WriteError(ctx.Response(), err)
	}

	if req.ResourceType != "posts" {
		return jsonapi.WriteError(ctx.Response(), jsonapi.NotFound("The requested resource is not available"))
	}

	var doc *jsonapi.Document
	if req.Intent.DocumentExpected() {
		doc, err = jsonapi.ParseBody(ctx.Request().Body())
		if err != nil {
			return jsonapi.WriteError(ctx.Response(), err)
		}
	}

	if req.Intent == jsonapi.ListResources {
		err = listPosts(req, ctx)
	} else if req.Intent == jsonapi.FindResource {
		err = findPost(req, ctx)
	} else if req.Intent == jsonapi.CreateResource {
		err = createPost(req, ctx, doc)
	} else if req.Intent == jsonapi.UpdateResource {
		err = updatePost(req, ctx, doc)
	} else if req.Intent == jsonapi.DeleteResource {
		err = deletePost(req, ctx)
	} else {
		err = jsonapi.BadRequest("The requested method is not available")
	}

	if err != nil {
		return jsonapi.WriteError(ctx.Response(), err)
	}

	return nil
}

func listPosts(_ *jsonapi.Request, ctx echo.Context) error {
	list := make([]*jsonapi.Resource, 0, len(store))
	for _, post := range store {
		list = append(list, &jsonapi.Resource{
			Type:       "posts",
			ID:         post.ID,
			Attributes: post,
		})
	}

	return jsonapi.WriteResources(ctx.Response(), http.StatusOK, list, &jsonapi.DocumentLinks{
		Self: "/api/posts",
	})
}

func findPost(req *jsonapi.Request, ctx echo.Context) error {
	post, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	return writePost(ctx, http.StatusOK, post)
}

func createPost(_ *jsonapi.Request, ctx echo.Context, doc *jsonapi.Document) error {
	post := &postModel{
		ID: strconv.Itoa(counter),
	}

	err := jsonapi.MapToStruct(doc.Data.One.Attributes, post)
	if err != nil {
		return err
	}

	counter++
	store[post.ID] = post

	return writePost(ctx, http.StatusCreated, post)
}

func updatePost(req *jsonapi.Request, ctx echo.Context, doc *jsonapi.Document) error {
	post, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	err := jsonapi.MapToStruct(doc.Data.One.Attributes, post)
	if err != nil {
		return err
	}

	return writePost(ctx, http.StatusOK, post)
}

func deletePost(req *jsonapi.Request, ctx echo.Context) error {
	_, ok := store[req.ResourceID]
	if !ok {
		return jsonapi.NotFound("The requested resource does not exist")
	}

	delete(store, req.ResourceID)

	ctx.Response().WriteHeader(http.StatusOK)
	return nil
}

func writePost(ctx echo.Context, status int, post *postModel) error {
	return jsonapi.WriteResource(ctx.Response(), status, &jsonapi.Resource{
		Type:       "posts",
		ID:         post.ID,
		Attributes: post,
	}, &jsonapi.DocumentLinks{
		Self: "/api/posts/" + post.ID,
	})
}
