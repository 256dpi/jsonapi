// This example uses the client to query the example API.
package main

import (
	"fmt"

	"github.com/256dpi/jsonapi"
)

var c = jsonapi.NewClient("http://0.0.0.0:4000/api")

type postModel struct {
	ID    string `json:"-"`
	Title string `json:"title"`
}

func main() {
	fmt.Println("==> Listing existing posts")
	posts := listPosts()
	fmt.Printf("%+v\n", posts)

	fmt.Println("==> Creating a new post")
	post := createPost("Hello world!")
	fmt.Printf("%+v\n", post)

	fmt.Println("==> Listing newly created posts")
	posts = listPosts()
	fmt.Printf("%+v\n", posts)

	fmt.Println("==> Updating created post")
	post.Title = "Amazing stuff!"
	post = updatePost(post)
	fmt.Printf("%+v\n", post)

	fmt.Println("==> Finding updated post")
	post = findPost(post.ID)
	fmt.Printf("%+v\n", post)

	fmt.Println("==> Deleting updated post")
	deletePost(post.ID)
	fmt.Println("ok")

	fmt.Println("==> Listing posts again")
	posts = listPosts()
	fmt.Printf("%+v\n", posts)
}

func listPosts() []postModel {
	doc, err := c.Request(&jsonapi.Request{
		Intent:       jsonapi.ListResources,
		ResourceType: "posts",
	}, nil)
	if err != nil {
		panic(err)
	}

	if doc == nil || doc.Data == nil {
		panic("missing resources")
	}

	posts := make([]postModel, len(doc.Data.Many))

	for i, resource := range doc.Data.Many {
		posts[i].ID = resource.ID
		resource.Attributes.Assign(&posts[i])
	}

	return posts
}

func createPost(title string) postModel {
	doc, err := c.RequestWithResource(&jsonapi.Request{
		Intent:       jsonapi.CreateResource,
		ResourceType: "posts",
	}, &jsonapi.Resource{
		Type: "posts",
		Attributes: jsonapi.Map{
			"title": title,
		},
	})
	if err != nil {
		panic(err)
	}

	if doc == nil || doc.Data == nil || doc.Data.One == nil {
		panic("missing resource")
	}

	return asPost(doc.Data.One)
}

func findPost(id string) postModel {
	doc, err := c.Request(&jsonapi.Request{
		Intent:       jsonapi.FindResource,
		ResourceType: "posts",
		ResourceID:   id,
	}, nil)
	if err != nil {
		panic(err)
	}

	if doc == nil || doc.Data == nil || doc.Data.One == nil {
		panic("missing resource")
	}

	return asPost(doc.Data.One)
}

func updatePost(post postModel) postModel {
	doc, err := c.RequestWithResource(&jsonapi.Request{
		Intent:       jsonapi.UpdateResource,
		ResourceType: "posts",
		ResourceID:   post.ID,
	}, &jsonapi.Resource{
		Type: "posts",
		Attributes: jsonapi.Map{
			"title": post.Title,
		},
	})
	if err != nil {
		panic(err)
	}

	if doc == nil || doc.Data == nil || doc.Data.One == nil {
		panic("missing resource")
	}

	return asPost(doc.Data.One)
}

func deletePost(id string) {
	_, err := c.Request(&jsonapi.Request{
		Intent:       jsonapi.DeleteResource,
		ResourceType: "posts",
		ResourceID:   id,
	}, nil)
	if err != nil {
		panic(err)
	}
}

func asPost(resource *jsonapi.Resource) postModel {
	var post postModel

	post.ID = resource.ID
	resource.Attributes.Assign(&post)

	return post
}
