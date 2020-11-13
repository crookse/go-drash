package resources

import (
	"fmt"
	"../src/http"

	"github.com/valyala/fasthttp"
)

func UsersResource() *http.Resource {

	resource := new(http.Resource)

	resource.Uris = []string{"/"}

	resource.Methods = map[string]interface{}{
		"GET": func (ctx *fasthttp.RequestCtx) {
			fmt.Fprintf(ctx, "Hello World!")
		},
	}

	return resource;
}
