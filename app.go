package main

import (
	"flag"
	"fmt"
	"log"

	"./resources"

	godrash "github.com/drashland/go-drash"
	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", ":1997", "TCP address to listen to")
)

func main() {
	flag.Parse()

	s := new(godrash.http.Server)

	s.AddResources(
		resources.HomeResource,
		resources.UsersResource,
	)

	fmt.Println("Server started at " + *addr)

	err := fasthttp.ListenAndServe(*addr, s.HandleRequest)

	if err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
