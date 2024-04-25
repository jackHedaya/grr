package utils

import (
	"fmt"

	"github.com/jackHedaya/grr/grr"
)

// #############################################################################
// # ErrFailedToGetPackagePath
// #############################################################################

type ErrFailedToGetPackagePath struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrFailedToGetPackagePath{}

func NewErrFailedToGetPackagePath() grr.Error {
	return &ErrFailedToGetPackagePath{}
}

func (e *ErrFailedToGetPackagePath) Error() string {
	return fmt.Sprintf("failed to get package path")
}

func (e *ErrFailedToGetPackagePath) Unwrap() error {
	return e.err
}

func (e *ErrFailedToGetPackagePath) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToGetPackagePath) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToGetPackagePath) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToGetPackagePath) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToGetPackagePath) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToGetPackagePath) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToGetPackagePath) GetOp() string {
	return e.op
}

func (e *ErrFailedToGetPackagePath) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToGetPackagePath) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrFailedToResolveDir
// #############################################################################

type ErrFailedToResolveDir struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrFailedToResolveDir{}

func NewErrFailedToResolveDir() grr.Error {
	return &ErrFailedToResolveDir{}
}

func (e *ErrFailedToResolveDir) Error() string {
	return fmt.Sprintf("failed to resolve absolute path")
}

func (e *ErrFailedToResolveDir) Unwrap() error {
	return e.err
}

func (e *ErrFailedToResolveDir) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToResolveDir) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToResolveDir) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToResolveDir) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToResolveDir) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToResolveDir) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToResolveDir) GetOp() string {
	return e.op
}

func (e *ErrFailedToResolveDir) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToResolveDir) Trace() {
	grr.Trace(e)
}
