package grr

import (
	"fmt"
	"reflect"
	"strings"
)

type Error interface {
	Error() string
	Unwrap() error
	UnwrapAll() error
	AsGrr(err Error) (Error, bool)
	AddTrait(key Trait, value any) Error
	GetTrait(key Trait) (any, bool)
	AddOp(op string) Error
	GetOp() string
	AddError(err error) Error
	GetTraits() map[Trait]any
	Trace()
	Strace() string
}

// grrError implements the Error interface
var _ Error = &grrError{}

type grrError struct {
	err    error
	msg    string
	op     string
	traits map[Trait]any
}

func Errorf(format string, args ...interface{}) Error {
	return &grrError{msg: fmt.Sprintf(format, args...), traits: map[Trait]any{}}
}

func (e *grrError) Error() string {
	return e.msg
}

func (e *grrError) Unwrap() error {
	return e.err
}

func (e *grrError) UnwrapAll() error {
	return UnwrapAll(e)
}

func (e *grrError) AsGrr(err Error) (Error, bool) {
	return AsGrr(e, err)
}

func (e *grrError) AddTrait(key Trait, value any) Error {
	e.traits[key] = value
	return e
}

func (e *grrError) GetTrait(key Trait) (any, bool) {
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

func (e *grrError) GetTraits() map[Trait]any {
	traits := make(map[Trait]any)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *grrError) Trace() {
	Trace(e)
}

func (e *grrError) Strace() string {
	return Strace(e)
}

func Trace(err error) {
	fmt.Println(Strace(err))
}

func Strace(err error) string {
	if !IsGrr(err) {
		return err.Error()
	}

	var errs []error

	for {
		errs = append(errs, err)

		if casted, ok := err.(Error); !ok || casted.Unwrap() == nil {
			break
		}

		err = err.(Error).Unwrap()
	}

	var trace strings.Builder

	for i := len(errs) - 1; i >= 0; i-- {
		if i != len(errs)-1 {
			trace.WriteString("|- ")
		}

		trace.WriteString(errs[i].Error())

		if IsGrr(errs[i]) {
			op := errs[i].(Error).GetOp()

			if op != "" {
				trace.WriteString(fmt.Sprintf("; op: %s", op))
			}
		}

		trace.WriteString("\n")
	}

	return trace.String()
}

func IsGrr(err error) bool {
	_, ok := err.(Error)
	return ok
}

func AsGrr(e Error, err error) (Error, bool) {
	E := reflect.TypeOf(err)

	var last error = e.UnwrapAll()

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

// Gets the trait value of the **innermost** grr.Error in the chain
// This let's you assign a trait to the root error and have it propogate down the stack
func GetTrait(err error, key Trait) (any, bool) {
	if !IsGrr(err) {
		return nil, false
	}

	bottomGrr := UnwrapAllGrr(err.(Error))

	return bottomGrr.GetTrait(key)
}

// Unwraps to the bottom-most grr.Error in the chain. This is the closest grr.Error to the root error
func UnwrapAllGrr(err Error) Error {
	for {
		if casted, ok := err.Unwrap().(Error); ok {
			err = casted
			continue
		}

		return err
	}
}

func UnwrapAll(e Error) error {
	last := e.Unwrap()

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
