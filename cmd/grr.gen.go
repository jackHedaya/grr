package main

import (
	"fmt"
	"github.com/jackHedaya/grr/grr"
)

// #############################################################################
// # ErrPathNotFound
// #############################################################################

type ErrPathNotFound struct {
	err    error
	op     string
	traits map[grr.Trait]string
	path   string
}

var _ grr.Error = &ErrPathNotFound{}

func NewErrPathNotFound(path string) grr.Error {
	return &ErrPathNotFound{
		path: path,
	}
}

func (e *ErrPathNotFound) Error() string {
	return fmt.Sprintf("path %s not found", e.path)
}

func (e *ErrPathNotFound) Unwrap() error {
	return e.err
}

func (e *ErrPathNotFound) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrPathNotFound) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrPathNotFound) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrPathNotFound) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrPathNotFound) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrPathNotFound) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrPathNotFound) GetOp() string {
	return e.op
}

func (e *ErrPathNotFound) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrPathNotFound) Trace() {
	grr.Trace(e)
}
