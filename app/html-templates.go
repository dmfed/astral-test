package app

import (
	"astral/storage"
	_ "embed"
	"html/template"
)

//go:embed html/get.html
var stringGet string

//go:embed html/put.html
var stringPut string

//go:embed html/ok.html
var stringOK string

var htmlGet = template.Must(template.New("get endpoint").Parse(stringGet))
var htmlPut = template.Must(template.New("put endpoint").Parse(stringPut))
var htmlOK = template.Must(template.New("result page").Parse(stringOK))

type PageGet struct {
	Elements []storage.Element
}

type PageOK struct {
	Element storage.Element
}

type PagePut struct {
	Username string
}
