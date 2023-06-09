package errgen

import (
	"fmt"
	"text/template"
)

const errorsTmpl = `// Code generated by go tool eg; DO NOT EDIT.

package {{.PackageName}}

import (
	"errors"
	"fmt"
)

type Error struct {
	Status  int    {{jsonTag "-"}}
	Code    int    {{jsonTag "code"}}
	Message string {{jsonTag "message"}}
	Detail  string {{jsonTag "detail"}}
}

var (
{{- range $k, $v := .Errors}}
	{{$k}} = &Error{
		Code:    {{$v.Code}},
		Status:  {{$v.Status}},
		Message: "{{$v.Message}}",
	}
{{- end}}
)

func (err *Error) Error() string {
	return err.Message
}

// WithDetail Clone a new Error and set detail for it
func (err *Error) WithDetail(detail string) *Error {
	return &Error{
		Status:  err.Status,
		Code:    err.Code,
		Message: err.Message,
		Detail:  detail,
	}
}

// WithDetailf Clone a new Error and set detail for it
func (err *Error) WithDetailf(format string, detail ...any) *Error {
	return &Error{
		Status:  err.Status,
		Code:    err.Code,
		Message: err.Message,
		Detail:  fmt.Sprintf(format, detail...),
	}
}

func (err *Error) Is(e error) bool {
	var temp *Error
	if !errors.As(e, &temp) {
		return false
	}
	return err.Code == temp.Code
}


// Is a wrapper of built-in errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As a wrapper of built-in errors.As
func As(err error, target any) bool {
	return errors.As(err, target)
}

// New a wrapper of built-in errors.New
func New(text string) error {
	return errors.New(text)
}

`

var errorsTemplate = template.Must(template.New("errors").Funcs(template.FuncMap{
	"jsonTag": func(s string) string {
		return fmt.Sprintf("`json:\"%s\"`", s)
	},
}).Parse(errorsTmpl))
