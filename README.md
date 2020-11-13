# Go Drash

A REST microframework for Go.

## Quickstart

1. Create your resource.

```go
// File: /path/to/your/project/resources/home_resource.go

package resources

import (
	"fmt"
	"github.com/drashland/go-drash/http"
)

func HomeResource() *http.Resource {

	resource := new(http.Resource)

	resource.Uris = []string{"/"}

	resource.Methods = map[string]interface{}{
		"GET": get,
	}

	return resource;
}

// Registered, so this should output a response as expected
func get(request http.Request) {
	fmt.Fprintf(request.Ctx, "GET  World!")
}

// Not registered in Methods, so this should throw a 405
func post(request http.Request) {
	fmt.Fprintf(request.Ctx, "POST World!")
}
```

2. Create your app.

```go
// File: /path/to/your/project/app.go

package main

import (
	"flag"
	"fmt"

	"./resources"
	"github.com/drashland/go-drash/http"
)

var (
	addr = flag.String("addr", "localhost:1997", "TCP address to listen to")
)

func main() {
	flag.Parse()

	s := new(http.Server)

	s.AddResources(
		resources.HomeResource,
		resources.UsersResource,
	)

	fmt.Println("Server started at " + *addr)

	s.Run(*addr)
}
```

3. Run your app.

```shell
$ go get
$ go run app.go
```

4. Make a request.

```
$ curl localhost:1997
GET World!
```
