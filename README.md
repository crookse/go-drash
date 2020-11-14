# Go Drash

A REST microframework for Go -- built on top of fasthttp.

## Quickstart

1. Create your resource.

```go
// File: /path/to/your/project/resources/home_resource.go

package resources

import (
  "github.com/drashland/go-drash/http"
)

func HomeResource() http.Resource {
	return http.Resource{

		Uris: []string{
			"/hello/:name",
		},

		GET: func (r http.Request) http.Response {
			r.Response.Body = "Hello World!"
			return r.Response
		},
	}
}

// This is registered, so it will output as expected
func get(r http.Request) http.Response {
  r.response.Body = "Hello World! Go + Drash is cool!"
  return r.response
}

// This is not registered, so it will throw a 405 error
func post(r http.Request) http.Response {
  r.response.Body = "test"
  return r.response
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

func main() {
	s := new(http.Server)

	s.AddResources(
		resources.HomeResource,
	)

	o := http.HttpOptions{
		Hostname: "localhost",
		Port: 1997,
	}

	fmt.Println("Server started at http://localhost:1997")

	s.Run(o)
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
Hello World! Go + Drash is cool!
```
