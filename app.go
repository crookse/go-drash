package main

import (
	"flag"
	"fmt"
	"log"

	"./src/http"

	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", ":1997", "TCP address to listen to")
)

// Create the server.
func createServer() *http.Server {
	resource := createResources()
	
	s := new(http.Server)
	s.AddResource(resource)

	return s
}

// Create the resources for the server.
func createResources() *http.Resource {

	homeResource := new(http.Resource)

	homeResource.Methods = map[string]interface{}{
		"GET": func(ctx *fasthttp.RequestCtx) {
			fmt.Fprintf(ctx, "Hi, Sara!")
		},
	}

	return homeResource
}

func main() {
	flag.Parse()

	s := createServer();

	fmt.Println("Server started at " + *addr)

	if err := fasthttp.ListenAndServe(*addr, s.HandleRequest); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
