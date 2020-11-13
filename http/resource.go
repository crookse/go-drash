package http

type methods map[string]interface{}

type Resource struct {
	Uris []string
	Methods methods
	response Response
}
