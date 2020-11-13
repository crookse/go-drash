# Go Drash

A REST microframework for Go.

## Quickstart

1. Create your resource.

```go
// File: /path/to/your/project/home_resource.go

package resources

import (
	"fmt"
	"github.com/drashland/godrash/http"
)

func HomeResource() *http.Resource {

	resource := new(http.Resource)

	resource.Uris = []string{"/"}

	resource.Methods = map[string]interface{}{
		"GET": get,
	}

	return resource;
}

func get(request http.Request) {
	fmt.Fprintf(request.Ctx, "GET  World!")
}

func post(request http.Request) {
	fmt.Fprintf(request.Ctx, "POST World!")
}
```

2. Create your app.

```go
package main

import (
	"flag"
	"fmt"

	"./resources"
	"github.com/drashland/godrash/http"
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
