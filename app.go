package main

import (
	"flag"
	"fmt"
	"log"

	"./resources"
	"./src/http"

	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", ":1997", "TCP address to listen to")
)

func main() {
	flag.Parse()

	s := new(http.Server)

	s.AddResources(
		resources.HomeResource())

	fmt.Println("Server started at " + *addr)

	if err := fasthttp.ListenAndServe(*addr, s.HandleRequest); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
