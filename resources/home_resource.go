package resources

import (
	"fmt"
	"../src/http"

	"github.com/valyala/fasthttp"
)

func HomeResource() *http.Resource {

	resource := new(http.Resource)

	resource.Uris = []string{"/"}

	resource.Methods = map[string]interface{}{
		"GET": get,
	}

	return resource;
}

func get(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "GET  World!")
}

func post(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "POST World!")
}
