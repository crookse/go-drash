# Go Drash

A REST microframework for Go built on top of [fasthttp](https://github.com/valyala/fasthttp).

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
			r.Response.Body = "Hello World! Go + Drash is cool!"
			return r.Response
		},
	}
}
```

2. Create your app.

```go
// File: /path/to/your/project/app.go

package main

import (
	"fmt"

	"./resources"
	"github.com/drashland/go-drash/http"
)

func main() {
	server := http.Server{
		Resources: []func() http.Resource{
			resources.HomeResource,
		},
	}

	options := http.HttpOptions{
		Hostname: "localhost",
		Port: 1997,
	}

	fmt.Println(fmt.Sprintf(
		"Server started at http://%s:%d", options.Hostname, options.Port,
	))

	server.Run(options)
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
