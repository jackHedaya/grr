package gen

import (
	"bytes"
	_ "embed"
	"go/format"
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
	Expr string
	Name string
	Type string
}

type TemplateData struct {
	PkgName string
	ErrName string
	Vars    []FnArg
	Message string
	Imports []string
}

var TrIsInternal = grr.NewTrait("IsInternal")

type GenerateFileArgs struct {
	Imports []string
	PkgName string
	ErrMsg  string
	Args    []FnArg
}

func GenerateErrorFile(params GenerateFileArgs) (string, error) {
	op := "Generate"

	errMsg := params.ErrMsg
	pkgName := params.PkgName
	args := params.Args
	imports := params.Imports

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
	errName := "Err" + strings.TrimSpace(split[0])

	if errName == "" {
		return "", grr.Errorf("NoErrorName: error name not found in error message")
	}

	// extract the error message from the message
	errMsg = strings.TrimSpace(strings.Join(split[1:], ""))

	var buf bytes.Buffer

	err := tmplGet.Execute(&buf, TemplateData{
		PkgName: pkgName,
		ErrName: errName,
		Vars:    args,
		Message: errMsg,
		Imports: imports,
	})

	if err != nil {
		return "", grr.Errorf("FailedToExecuteTemplate: something went wrong while generating: %v", strings.Builder{}).
			AddError(err).
			AddTrait(TrIsInternal, "true").
			AddOp(op)
	}

	fmted, err := format.Source(buf.Bytes())

	if err != nil {
		return "", grr.Errorf("FailedToFormatSource: something went wrong while formatting source: %v", strings.Builder{}).
			AddError(err).
			AddOp(op)
	}

	return string(fmted), nil
}

func GenDefaultImports() []string {
	return []string{"fmt", "reflect"}
}
