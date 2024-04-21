package gen

import (
	_ "embed"
	"strings"

	"text/template"

	"github.com/whiskaway/grr/grr"
)

//go:embed errorFunc.tmpl
var tmpl string

var fns = template.FuncMap{
	"notlast": func(index int, len int) bool {
		return index+1 != len
	},
}

var tmplGet = template.Must(template.New("").Funcs(fns).Parse(tmpl))

type FnArg struct {
	Literal string
	Type    string
}

type TemplateData struct {
	ErrName string
	Vars    []FnArg
	Message string
}

var TrIsInternal = grr.NewTrait("IsInternal")

func GenerateErrorFunction(errMsg string, args ...FnArg) (string, error) {
	op := "Generate"

	if len(errMsg) == 0 || errMsg == "\"\"" {
		return "", grr.Errorf("NoErrorMessage: error message not found").
			AddOp(op)
	}

	if errMsg[0] == '"' {
		errMsg = errMsg[1:]
	}

	if errMsg[len(errMsg)-1] == '"' {
		errMsg = errMsg[:len(errMsg)-1]
	}

	// extract the new error name from message: "FileNotFound: a file with name %s was not found" => "FileNotFound"
	split := strings.Split(errMsg, ":")

	// extract the error name from the message
	errName := strings.TrimSpace(split[0])

	if errName == "" {
		return "", grr.Errorf("NoErrorName: error name not found in error message")
	}

	// extract the error message from the message
	errMsg = strings.TrimSpace(strings.Join(split[1:], ""))

	var buf strings.Builder

	err := tmplGet.Execute(&buf, TemplateData{
		ErrName: errName,
		Vars:    args,
		Message: errMsg,
	})

	if err != nil {
		return "", grr.Errorf("FailedToExecuteTemplate: something went wrong while generating").
			AddError(err).
			AddTrait(TrIsInternal, "true").
			AddOp(op)
	}

	return buf.String(), nil
}
