package http

import (
)

type methods map[string]interface{}

type Resource struct {
	uris []string
	response Response
	Methods methods
}
