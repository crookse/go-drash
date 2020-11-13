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
		"GET": GET,
	}

	return resource;
}

func GET(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello World!")
}
