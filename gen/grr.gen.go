package gen

import (
	"fmt"
	"github.com/whiskaway/grr/grr"
	"reflect"
	"strings"
)

type FailedToFormatSource struct {
	err            error
	op             string
	traits         map[grr.Trait]string
	stringsBuilder strings.Builder
}

var _ grr.Error = &FailedToFormatSource{}

func NewFailedToFormatSource(stringsBuilder strings.Builder) grr.Error {
	return &FailedToFormatSource{
		stringsBuilder: stringsBuilder,
	}
}

func (e *FailedToFormatSource) Error() string {
	return fmt.Sprintf("something went wrong while formatting source %v", e.stringsBuilder)
}

func (e *FailedToFormatSource) Unwrap() error {
	return e.err
}

func (e *FailedToFormatSource) AsGrr(err grr.Error) (grr.Error, bool) {
	E := reflect.TypeOf(err)

	var last error = e.err

	for {
		if last == nil {
			return nil, false
		}

		if _, ok := last.(grr.Error); !ok {
			return nil, false
		}

		if reflect.TypeOf(last).ConvertibleTo(E) {
			return reflect.ValueOf(last).Convert(E).Interface().(grr.Error), true
		}

		last = last.(grr.Error).Unwrap()
	}
}

func (e *FailedToFormatSource) IsGrr(err grr.Error) bool {
	_, ok := e.AsGrr(err)
	return ok
}

func (e *FailedToFormatSource) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *FailedToFormatSource) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *FailedToFormatSource) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *FailedToFormatSource) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *FailedToFormatSource) GetOp() string {
	return e.op
}

func (e *FailedToFormatSource) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *FailedToFormatSource) Trace() {
	var err error = e
	var errs []error

	for {
		errs = append(errs, err)

		if casted, ok := err.(grr.Error); !ok || casted.Unwrap() == nil {
			break
		}

		err = err.(grr.Error).Unwrap()
	}

	for i := len(errs) - 1; i >= 0; i-- {
		if i != len(errs)-1 {
			fmt.Printf("|- ")
		}

		fmt.Print(errs[i].Error())

		if grr.IsGrr(errs[i]) {
			op := errs[i].(grr.Error).GetOp()

			if op != "" {
				fmt.Printf("; op: %s", op)
			}
		}

		fmt.Println()
	}
}
