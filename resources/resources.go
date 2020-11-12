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
		"GET": func(ctx *fasthttp.RequestCtx) {
			fmt.Println(string(ctx.Method()))
			fmt.Fprintf(ctx, "Hi, Sara!")
		},
	}

	return resource;
}
