package http

import (
)

type methods map[string]interface{}

type Resource struct {
	Uris []string
	response Response
	Methods methods
}
