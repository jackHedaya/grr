package gen

import (
	"bytes"
	_ "embed"
	"go/format"
	"os"
	"regexp"
	"slices"
	"strings"

	"text/template"

	"github.com/jackHedaya/grr/grr"
	"github.com/jackHedaya/grr/utils"
)

//go:embed errorStruct.tmpl
var errorStructTemplateStr string

//go:embed errorFile.tmpl
var errorFileTemplateStr string

var templateFuncs = template.FuncMap{
	"notlast": func(index int, len int) bool {
		return index+1 != len
	},
}

var errorStructTemplate = template.Must(template.New("").Funcs(templateFuncs).Parse(errorStructTemplateStr))
var errorFileTemplate = template.Must(template.New("").Funcs(templateFuncs).Parse(errorFileTemplateStr))

type GrrGenErrorField struct {
	Expr string
	Name string
	Type string
}

type StructTemplateData struct {
	ErrName string
	Vars    []GrrGenErrorField
	Message string
}

type HeaderTemplateData struct {
	PkgName         string
	Imports         []string
	GeneratedErrors []GeneratedError
}

var TrIsInternal = grr.NewTrait("IsInternal")
var TrIsNonFatal = grr.NewTrait("IsNonFatal")

type GenerateFileArgs struct {
	ErrMsg string
	Args   []GrrGenErrorField
}

func (f *grrWalker) GenerateErrorStruct(params GenerateFileArgs) (*GeneratedError, error) {
	op := "GenerateStruct"

	errMsg := params.ErrMsg
	args := params.Args

	if len(errMsg) == 0 || errMsg == "\"\"" {
		return nil, grr.Errorf("NoErrorMessage: error message not found").
			AddOp(op)
	}

	if errMsg[0] == '"' {
		errMsg = errMsg[1:]
	}

	if errMsg[len(errMsg)-1] == '"' {
		errMsg = errMsg[:len(errMsg)-1]
	}

	// extract the new error name from message
	// For example, "FileNotFound: a file with name %s was not found" =>
	// Error Name: FileNotFound
	// Error Message: a file with name %s was not found
	re := regexp.MustCompile(`^([A-Z][a-zA-Z]+):\s(.*)`)
	matches := re.FindStringSubmatch(errMsg)

	if len(matches) < 3 {
		return nil, grr.Errorf("NoErrorName: error name not found in error message")
	}

	errName := "Err" + matches[1]

	if errName == "" {
		return nil, grr.Errorf("NoErrorName: error name not found in error message")
	}

	errMsg = strings.TrimSpace(matches[2])

	// we need to check if the error name, args, and message are already defined in the package
	isDefined, isConflict := isAlreadyDefined(f, errName, args, errMsg)

	if isConflict {
		return nil, grr.Errorf("Conflict: error \"%s\" is already defined with different arguments or message", errName).
			AddTrait(TrIsInternal, "false").
			AddOp(op)
	}

	if isDefined {
		return nil, grr.Errorf("AlreadyDefined: error \"%s\" is already defined with the same arguments and message", errName).
			AddTrait(TrIsInternal, "false").
			AddTrait(TrIsNonFatal, "true").
			AddOp(op)
	}

	var buf bytes.Buffer

	err := errorStructTemplate.Execute(&buf, StructTemplateData{
		ErrName: errName,
		Vars:    args,
		Message: errMsg,
	})

	if err != nil {
		return nil, grr.Errorf("FailedToExecuteTemplate: something went wrong while generating: %v", strings.Builder{}).
			AddError(err).
			AddTrait(TrIsInternal, "true").
			AddOp(op)
	}

	return &GeneratedError{
		Name:          errName,
		Args:          args,
		Msg:           errMsg,
		GeneratedCode: buf.String(),
	}, nil
}

func GenerateErrorFile(pkgName string, imports []string, errors map[string]GeneratedError) ([]byte, error) {
	op := "GenerateErrorFile"

	pairs := utils.MapToPairs(errors)

	slices.SortStableFunc(pairs, func(i, j utils.Pair[string, GeneratedError]) int {
		return strings.Compare(i.Key, j.Key)
	})

	var headerBuff bytes.Buffer

	err := errorFileTemplate.Execute(&headerBuff, HeaderTemplateData{
		PkgName:         pkgName,
		Imports:         imports,
		GeneratedErrors: utils.PairValues(pairs),
	})

	if err != nil {
		return nil, grr.Errorf("FailedToExecuteTemplate: something went wrong while generating: %v", strings.Builder{}).
			AddError(err).
			AddTrait(TrIsInternal, "true").
			AddOp(op)
	}

	fmted, err := format.Source(headerBuff.Bytes())

	if err != nil {
		// write for debugging
		os.WriteFile("grr.failed.gen.go", headerBuff.Bytes(), 0644)

		return nil, grr.Errorf("FailedToFormat: failed to format the generated code").
			AddError(err).
			AddOp(op)
	}

	return fmted, nil
}

func GenDefaultImports() []string {
	return []string{"fmt"}
}

func isAlreadyDefined(f *grrWalker, errName string, args []GrrGenErrorField, errMsg string) (bool, bool) {
	errs := utils.Merge(f.prevErrors, f.generatedErrors)

	// check if the error name is already defined
	if _, ok := errs[errName]; !ok {
		return false, false
	}

	prevErr := errs[errName]

	// if it is defined, check if the arguments and message are the same
	if len(prevErr.Args) != len(args) {
		return true, true
	}

	if prevErr.Msg != errMsg {
		return true, true
	}

	for i, arg := range args {
		if prevErr.Args[i].Type != arg.Type {
			return true, true
		}

		if prevErr.Args[i].Name != arg.Name {
			return true, true
		}
	}

	return false, false
}
