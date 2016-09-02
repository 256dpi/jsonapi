package jsonapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

func Example() {
	http.HandleFunc("/api/posts/1", func(w http.ResponseWriter, r *http.Request) {
		req, err := ParseRequest(r, "/api/")
		if err != nil {
			WriteError(w, err)
			return
		}

		fmt.Println(req.Resource)
		fmt.Println(req.ResourceID)
	})

	go func() {
		http.ListenAndServe("0.0.0.0:4040", nil)
	}()

	time.Sleep(50 * time.Millisecond)

	_, str, err := gorequest.New().
		Get("http://0.0.0.0:4040/api/posts/1").
		Set("Accept", ContentType).
		End()
	if err != nil {
		panic(err[0])
	}

	fmt.Println(str)

	// Output:
	// posts
	// 1
}
