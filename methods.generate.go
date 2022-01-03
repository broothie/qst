// +build generate

package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var methods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

func main() {
	file, err := os.OpenFile("methods.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	tmpl, err := template.
		New("template").
		Funcs(map[string]interface{}{"lower": strings.ToLower, "title": strings.Title}).
		ParseFiles("methods.go.tmpl")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tmpl.ExecuteTemplate(file, "methods.go.tmpl", methods); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
