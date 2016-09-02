package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gonfire/jsonapi"
)

type postModel struct {
	ID    string
	Title string
}

func main() {
	router := gin.Default()

	counter := 1
	store := make(map[string]*postModel)

	router.GET("/api/posts", func(ctx *gin.Context) {
		_, err := jsonapi.ParseRequest(ctx.Request, "/api/")
		if err != nil {
			jsonapi.WriteError(ctx.Writer, err)
			return
		}

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

		jsonapi.WriteResponse(ctx.Writer, http.StatusOK, &jsonapi.Document{
			Data: &jsonapi.HybridResource{
				Many: list,
			},
		})
	})

	router.POST("/api/posts", func(ctx *gin.Context) {
		_, err := jsonapi.ParseRequest(ctx.Request, "/api/")
		if err != nil {
			jsonapi.WriteError(ctx.Writer, err)
			return
		}

		doc, err := jsonapi.ParseBody(ctx.Request.Body)
		if err != nil {
			jsonapi.WriteError(ctx.Writer, err)
			return
		}

		post := &postModel{
			ID:    strconv.Itoa(counter),
			Title: doc.Data.One.Attributes["title"].(string),
		}

		counter++
		store[post.ID] = post

		writePost(ctx, post)
	})

	router.GET("/api/posts/:id", func(ctx *gin.Context) {
		req, err := jsonapi.ParseRequest(ctx.Request, "/api/")
		if err != nil {
			jsonapi.WriteError(ctx.Writer, err)
			return
		}

		post, ok := store[req.ResourceID]
		if !ok {
			jsonapi.WriteErrorFromStatus(ctx.Writer, http.StatusNotFound)
			return
		}

		writePost(ctx, post)
	})

	router.PATCH("/api/posts/:id", func(ctx *gin.Context) {
		req, err := jsonapi.ParseRequest(ctx.Request, "/api/")
		if err != nil {
			jsonapi.WriteError(ctx.Writer, err)
			return
		}

		doc, err := jsonapi.ParseBody(ctx.Request.Body)
		if err != nil {
			jsonapi.WriteError(ctx.Writer, err)
			return
		}

		post, ok := store[req.ResourceID]
		if !ok {
			jsonapi.WriteErrorFromStatus(ctx.Writer, http.StatusNotFound)
			return
		}

		post.Title = doc.Data.One.Attributes["title"].(string)

		writePost(ctx, post)
	})

	router.DELETE("/api/posts/:id", func(ctx *gin.Context) {
		req, err := jsonapi.ParseRequest(ctx.Request, "/api/")
		if err != nil {
			jsonapi.WriteError(ctx.Writer, err)
			return
		}

		_, ok := store[req.ResourceID]
		if !ok {
			jsonapi.WriteErrorFromStatus(ctx.Writer, http.StatusNotFound)
			return
		}

		delete(store, req.ResourceID)

		ctx.Status(http.StatusOK)
	})

	router.Run("0.0.0.0:4000")
}

func writePost(ctx *gin.Context, post *postModel) {
	jsonapi.WriteResponse(ctx.Writer, http.StatusOK, &jsonapi.Document{
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
