package clean

import (
	"fmt"

	"github.com/jackHedaya/grr/grr"
)

// #############################################################################
// # ErrFailedToClean
// #############################################################################

type ErrFailedToClean struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrFailedToClean{}

func NewErrFailedToClean() grr.Error {
	return &ErrFailedToClean{}
}

func (e *ErrFailedToClean) Error() string {
	return fmt.Sprintf("Failed to clean directory")
}

func (e *ErrFailedToClean) Unwrap() error {
	return e.err
}

func (e *ErrFailedToClean) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToClean) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToClean) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToClean) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToClean) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToClean) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToClean) GetOp() string {
	return e.op
}

func (e *ErrFailedToClean) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToClean) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrFailedToDelete
// #############################################################################

type ErrFailedToDelete struct {
	err    error
	op     string
	traits map[grr.Trait]string
	path   string
}

var _ grr.Error = &ErrFailedToDelete{}

func NewErrFailedToDelete(path string) grr.Error {
	return &ErrFailedToDelete{
		path: path,
	}
}

func (e *ErrFailedToDelete) Error() string {
	return fmt.Sprintf("Failed to delete %s", e.path)
}

func (e *ErrFailedToDelete) Unwrap() error {
	return e.err
}

func (e *ErrFailedToDelete) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToDelete) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToDelete) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToDelete) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToDelete) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToDelete) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToDelete) GetOp() string {
	return e.op
}

func (e *ErrFailedToDelete) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToDelete) Trace() {
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
	return fmt.Sprintf("Failed to resolve directory")
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

// #############################################################################
// # ErrFailedToWalk
// #############################################################################

type ErrFailedToWalk struct {
	err    error
	op     string
	traits map[grr.Trait]string
	path   string
}

var _ grr.Error = &ErrFailedToWalk{}

func NewErrFailedToWalk(path string) grr.Error {
	return &ErrFailedToWalk{
		path: path,
	}
}

func (e *ErrFailedToWalk) Error() string {
	return fmt.Sprintf("Failed to walk %s", e.path)
}

func (e *ErrFailedToWalk) Unwrap() error {
	return e.err
}

func (e *ErrFailedToWalk) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToWalk) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToWalk) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToWalk) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToWalk) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToWalk) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToWalk) GetOp() string {
	return e.op
}

func (e *ErrFailedToWalk) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToWalk) Trace() {
	grr.Trace(e)
}
