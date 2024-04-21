package grr

import (
	"fmt"
	"reflect"
)

type Error interface {
	Error() string
	Unwrap() error
	AsGrr(err Error) (Error, bool)
	IsGrr(err Error) bool
	AddTrait(key Trait, value string) Error
	GetTrait(key Trait) (string, bool)
	AddOp(op string) Error
	GetOp() string
	AddError(err error) Error
	GetTraits() map[Trait]string
	Trace()
}

// grrError implements the Error interface
var _ Error = &grrError{}

type grrError struct {
	err    error
	msg    string
	op     string
	traits map[Trait]string
}

func Errorf(format string, args ...interface{}) Error {
	return &grrError{msg: fmt.Sprintf(format, args...), traits: make(map[Trait]string)}
}

func (e *grrError) Error() string {
	return e.msg
}

func (e *grrError) Unwrap() error {
	return e.err
}

func (e *grrError) UnwrapAll() error {
	last := e.err

	for {
		if last == nil {
			return nil
		}

		if _, ok := last.(Error); !ok {
			return last
		}

		last = last.(Error).Unwrap()
	}
}

func (e *grrError) AsGrr(err Error) (Error, bool) {
	E := reflect.TypeOf(err)

	var last error = e.err

	for {
		if last == nil {
			return nil, false
		}

		if _, ok := last.(Error); !ok {
			return nil, false
		}

		if reflect.TypeOf(last).ConvertibleTo(E) {
			return reflect.ValueOf(last).Convert(E).Interface().(Error), true
		}

		last = last.(Error).Unwrap()
	}
}

func (e *grrError) IsGrr(err Error) bool {
	_, ok := e.AsGrr(err)
	return ok
}

func (e *grrError) AddTrait(key Trait, value string) Error {
	e.traits[key] = value
	return e
}

func (e *grrError) GetTrait(key Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *grrError) AddOp(op string) Error {
	e.op = op
	return e
}

func (e *grrError) GetOp() string {
	return e.op
}

func (e *grrError) AddError(err error) Error {
	e.err = err
	return e
}

func (e *grrError) GetTraits() map[Trait]string {
	traits := make(map[Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *grrError) Trace() {
	// trace like so:
	// an error occured; op: SomeOp
	// |- the next level error
	// |- the next level error

	var err error = e
	var errs []error

	for {
		errs = append(errs, err)

		if casted, ok := err.(Error); !ok || casted.Unwrap() == nil {
			break
		}

		err = err.(Error).Unwrap()
	}

	for i := len(errs) - 1; i >= 0; i-- {
		if i != len(errs)-1 {
			fmt.Printf("|- ")
		}

		fmt.Print(errs[i].Error())

		if IsGrr(errs[i]) {
			op := errs[i].(Error).GetOp()

			if op != "" {
				fmt.Printf("; op: %s", op)
			}
		}

		fmt.Println()
	}
}

func Trace(err error) {
	if IsGrr(err) {
		err.(Error).Trace()
		return
	}
}

func IsGrr(err error) bool {
	_, ok := err.(Error)
	return ok
}