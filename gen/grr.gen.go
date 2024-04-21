package gen

import (
	"fmt"
	"reflect"

	"github.com/whiskaway/grr/grr"
)

type FailedToExecuteTemplate struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &FailedToExecuteTemplate{}

func NewFailedToExecuteTemplate() grr.Error {
	return &FailedToExecuteTemplate{}
}

func (e *FailedToExecuteTemplate) Error() string {
	return fmt.Sprintf("something went wrong while generating")
}

func (e *FailedToExecuteTemplate) Unwrap() error {
	return e.err
}

func (e *FailedToExecuteTemplate) AsGrr(err grr.Error) (grr.Error, bool) {
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

func (e *FailedToExecuteTemplate) IsGrr(err grr.Error) bool {
	_, ok := e.AsGrr(err)
	return ok
}

func (e *FailedToExecuteTemplate) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *FailedToExecuteTemplate) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *FailedToExecuteTemplate) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *FailedToExecuteTemplate) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *FailedToExecuteTemplate) GetOp() string {
	return e.op
}

func (e *FailedToExecuteTemplate) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *FailedToExecuteTemplate) Trace() {
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
