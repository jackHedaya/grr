package main


import (
  "fmt"
  "reflect"
  "grr/grr"
)

type MyError struct {
  err error
  op string
  traits map[grr.Trait]string
  num int
}

var _ grr.Error = &MyError{}

func NewMyError(num int) grr.Error {
  return &MyError{
    num: num,
  }
}

func (e *MyError) Error() string {
  return fmt.Sprintf("this is your number %v", e.num)
}

func (e *MyError) Unwrap() error {
  return e.err
}

func (e *MyError) AsGrr(err grr.Error) (grr.Error, bool) {
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


func (e *MyError) IsGrr(err grr.Error) bool {
  _, ok := e.AsGrr(err)
  return ok
}

func (e *MyError) AddTrait(trait grr.Trait, value string) grr.Error {
  e.traits[trait] = value
  return e
}

func (e *MyError) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *MyError) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}


func (e *MyError) AddOp(op string) grr.Error {
  e.op = op
  return e
}

func (e *MyError) GetOp() string {
  return e.op
}

func (e *MyError) AddError(err error) grr.Error {
  e.err = err
  return e
}

func (e *MyError) Trace() {
	// trace like so:
	// an error occured; op: SomeOp
	// |- the next level error
	// |- the next level error

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
