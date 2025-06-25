package main

import (
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const (
	fileName     = "methods.go"
	templateName = "methods.go.tmpl"
)

//go:embed methods.go.tmpl
var templateFile string

type Method string

func (m Method) String() string {
	return string(m)
}

func (m Method) Capitalized() string {
	lowerCase := strings.ToLower(m.String())

	return strings.ToUpper(lowerCase[0:1]) + lowerCase[1:]
}

func (m Method) HTTPMethodName() string {
	return fmt.Sprintf("http.Method%s", m.Capitalized())
}

func (m Method) ConstructorName() string {
	return fmt.Sprintf("New%s", m.Capitalized())
}

var methods = []Method{
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
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() (err error) {
	methodsGo, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := methodsGo.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	methodsGoTmpl, err := template.New(templateName).Parse(templateFile)
	if err != nil {
		return err
	}

	if err := methodsGoTmpl.ExecuteTemplate(methodsGo, templateName, methods); err != nil {
		return err
	}

	return nil
}
