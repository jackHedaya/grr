package gen

import (
	"fmt"
	"github.com/jackHedaya/grr/grr"
	"strings"
)

// #############################################################################
// # ErrAlreadyDefined
// #############################################################################

type ErrAlreadyDefined struct {
	err     error
	op      string
	traits  map[grr.Trait]string
	errName string
}

var _ grr.Error = &ErrAlreadyDefined{}

func NewErrAlreadyDefined(errName string) grr.Error {
	return &ErrAlreadyDefined{
		errName: errName,
	}
}

func (e *ErrAlreadyDefined) Error() string {
	return fmt.Sprintf("error \"%s\" is already defined with the same arguments and message", e.errName)
}

func (e *ErrAlreadyDefined) Unwrap() error {
	return e.err
}

func (e *ErrAlreadyDefined) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrAlreadyDefined) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrAlreadyDefined) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrAlreadyDefined) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrAlreadyDefined) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrAlreadyDefined) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrAlreadyDefined) GetOp() string {
	return e.op
}

func (e *ErrAlreadyDefined) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrAlreadyDefined) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrConflict
// #############################################################################

type ErrConflict struct {
	err     error
	op      string
	traits  map[grr.Trait]string
	errName string
}

var _ grr.Error = &ErrConflict{}

func NewErrConflict(errName string) grr.Error {
	return &ErrConflict{
		errName: errName,
	}
}

func (e *ErrConflict) Error() string {
	return fmt.Sprintf("error \"%s\" is already defined with different arguments or message", e.errName)
}

func (e *ErrConflict) Unwrap() error {
	return e.err
}

func (e *ErrConflict) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrConflict) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrConflict) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrConflict) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrConflict) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrConflict) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrConflict) GetOp() string {
	return e.op
}

func (e *ErrConflict) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrConflict) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrFailedToExecuteTemplate
// #############################################################################

type ErrFailedToExecuteTemplate struct {
	err            error
	op             string
	traits         map[grr.Trait]string
	stringsBuilder strings.Builder
}

var _ grr.Error = &ErrFailedToExecuteTemplate{}

func NewErrFailedToExecuteTemplate(stringsBuilder strings.Builder) grr.Error {
	return &ErrFailedToExecuteTemplate{
		stringsBuilder: stringsBuilder,
	}
}

func (e *ErrFailedToExecuteTemplate) Error() string {
	return fmt.Sprintf("something went wrong while generating: %v", e.stringsBuilder)
}

func (e *ErrFailedToExecuteTemplate) Unwrap() error {
	return e.err
}

func (e *ErrFailedToExecuteTemplate) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToExecuteTemplate) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToExecuteTemplate) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToExecuteTemplate) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToExecuteTemplate) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToExecuteTemplate) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToExecuteTemplate) GetOp() string {
	return e.op
}

func (e *ErrFailedToExecuteTemplate) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToExecuteTemplate) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrFailedToFormat
// #############################################################################

type ErrFailedToFormat struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrFailedToFormat{}

func NewErrFailedToFormat() grr.Error {
	return &ErrFailedToFormat{}
}

func (e *ErrFailedToFormat) Error() string {
	return fmt.Sprintf("failed to format the generated code")
}

func (e *ErrFailedToFormat) Unwrap() error {
	return e.err
}

func (e *ErrFailedToFormat) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToFormat) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToFormat) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToFormat) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToFormat) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToFormat) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToFormat) GetOp() string {
	return e.op
}

func (e *ErrFailedToFormat) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToFormat) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrFailedToGenerateStruct
// #############################################################################

type ErrFailedToGenerateStruct struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrFailedToGenerateStruct{}

func NewErrFailedToGenerateStruct() grr.Error {
	return &ErrFailedToGenerateStruct{}
}

func (e *ErrFailedToGenerateStruct) Error() string {
	return fmt.Sprintf("failed to generate error struct")
}

func (e *ErrFailedToGenerateStruct) Unwrap() error {
	return e.err
}

func (e *ErrFailedToGenerateStruct) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToGenerateStruct) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToGenerateStruct) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToGenerateStruct) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToGenerateStruct) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToGenerateStruct) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToGenerateStruct) GetOp() string {
	return e.op
}

func (e *ErrFailedToGenerateStruct) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToGenerateStruct) Trace() {
	grr.Trace(e)
}

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
// # ErrFailedToLoad
// #############################################################################

type ErrFailedToLoad struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrFailedToLoad{}

func NewErrFailedToLoad() grr.Error {
	return &ErrFailedToLoad{}
}

func (e *ErrFailedToLoad) Error() string {
	return fmt.Sprintf("failed to load package")
}

func (e *ErrFailedToLoad) Unwrap() error {
	return e.err
}

func (e *ErrFailedToLoad) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToLoad) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToLoad) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToLoad) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToLoad) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToLoad) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToLoad) GetOp() string {
	return e.op
}

func (e *ErrFailedToLoad) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToLoad) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrFailedToLoadPackages
// #############################################################################

type ErrFailedToLoadPackages struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrFailedToLoadPackages{}

func NewErrFailedToLoadPackages() grr.Error {
	return &ErrFailedToLoadPackages{}
}

func (e *ErrFailedToLoadPackages) Error() string {
	return fmt.Sprintf("failed to load packages")
}

func (e *ErrFailedToLoadPackages) Unwrap() error {
	return e.err
}

func (e *ErrFailedToLoadPackages) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToLoadPackages) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToLoadPackages) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToLoadPackages) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToLoadPackages) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToLoadPackages) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToLoadPackages) GetOp() string {
	return e.op
}

func (e *ErrFailedToLoadPackages) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToLoadPackages) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrFailedToLoadPreviousErrors
// #############################################################################

type ErrFailedToLoadPreviousErrors struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrFailedToLoadPreviousErrors{}

func NewErrFailedToLoadPreviousErrors() grr.Error {
	return &ErrFailedToLoadPreviousErrors{}
}

func (e *ErrFailedToLoadPreviousErrors) Error() string {
	return fmt.Sprintf("failed to load previous errors")
}

func (e *ErrFailedToLoadPreviousErrors) Unwrap() error {
	return e.err
}

func (e *ErrFailedToLoadPreviousErrors) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToLoadPreviousErrors) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToLoadPreviousErrors) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToLoadPreviousErrors) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToLoadPreviousErrors) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToLoadPreviousErrors) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToLoadPreviousErrors) GetOp() string {
	return e.op
}

func (e *ErrFailedToLoadPreviousErrors) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToLoadPreviousErrors) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrFailedToWriteFile
// #############################################################################

type ErrFailedToWriteFile struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrFailedToWriteFile{}

func NewErrFailedToWriteFile() grr.Error {
	return &ErrFailedToWriteFile{}
}

func (e *ErrFailedToWriteFile) Error() string {
	return fmt.Sprintf("failed to write generated file")
}

func (e *ErrFailedToWriteFile) Unwrap() error {
	return e.err
}

func (e *ErrFailedToWriteFile) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrFailedToWriteFile) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrFailedToWriteFile) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrFailedToWriteFile) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrFailedToWriteFile) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrFailedToWriteFile) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrFailedToWriteFile) GetOp() string {
	return e.op
}

func (e *ErrFailedToWriteFile) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrFailedToWriteFile) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrGenerateErrorFile
// #############################################################################

type ErrGenerateErrorFile struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrGenerateErrorFile{}

func NewErrGenerateErrorFile() grr.Error {
	return &ErrGenerateErrorFile{}
}

func (e *ErrGenerateErrorFile) Error() string {
	return fmt.Sprintf("failed to generate error file")
}

func (e *ErrGenerateErrorFile) Unwrap() error {
	return e.err
}

func (e *ErrGenerateErrorFile) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrGenerateErrorFile) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrGenerateErrorFile) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrGenerateErrorFile) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrGenerateErrorFile) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrGenerateErrorFile) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrGenerateErrorFile) GetOp() string {
	return e.op
}

func (e *ErrGenerateErrorFile) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrGenerateErrorFile) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrNoErrorMessage
// #############################################################################

type ErrNoErrorMessage struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrNoErrorMessage{}

func NewErrNoErrorMessage() grr.Error {
	return &ErrNoErrorMessage{}
}

func (e *ErrNoErrorMessage) Error() string {
	return fmt.Sprintf("error message not found")
}

func (e *ErrNoErrorMessage) Unwrap() error {
	return e.err
}

func (e *ErrNoErrorMessage) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrNoErrorMessage) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrNoErrorMessage) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrNoErrorMessage) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrNoErrorMessage) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrNoErrorMessage) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrNoErrorMessage) GetOp() string {
	return e.op
}

func (e *ErrNoErrorMessage) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrNoErrorMessage) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrNoErrorName
// #############################################################################

type ErrNoErrorName struct {
	err    error
	op     string
	traits map[grr.Trait]string
}

var _ grr.Error = &ErrNoErrorName{}

func NewErrNoErrorName() grr.Error {
	return &ErrNoErrorName{}
}

func (e *ErrNoErrorName) Error() string {
	return fmt.Sprintf("error name not found in error message")
}

func (e *ErrNoErrorName) Unwrap() error {
	return e.err
}

func (e *ErrNoErrorName) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrNoErrorName) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrNoErrorName) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrNoErrorName) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrNoErrorName) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrNoErrorName) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrNoErrorName) GetOp() string {
	return e.op
}

func (e *ErrNoErrorName) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrNoErrorName) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrNoPackagesFound
// #############################################################################

type ErrNoPackagesFound struct {
	err            error
	op             string
	traits         map[grr.Trait]string
	stringsBuilder strings.Builder
}

var _ grr.Error = &ErrNoPackagesFound{}

func NewErrNoPackagesFound(stringsBuilder strings.Builder) grr.Error {
	return &ErrNoPackagesFound{
		stringsBuilder: stringsBuilder,
	}
}

func (e *ErrNoPackagesFound) Error() string {
	return fmt.Sprintf("no packages found in directory. string builder for testing: %v", e.stringsBuilder)
}

func (e *ErrNoPackagesFound) Unwrap() error {
	return e.err
}

func (e *ErrNoPackagesFound) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrNoPackagesFound) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrNoPackagesFound) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrNoPackagesFound) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrNoPackagesFound) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrNoPackagesFound) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrNoPackagesFound) GetOp() string {
	return e.op
}

func (e *ErrNoPackagesFound) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrNoPackagesFound) Trace() {
	grr.Trace(e)
}

// #############################################################################
// # ErrOneExpected
// #############################################################################

type ErrOneExpected struct {
	err    error
	op     string
	traits map[grr.Trait]string
	arg    int
}

var _ grr.Error = &ErrOneExpected{}

func NewErrOneExpected(arg int) grr.Error {
	return &ErrOneExpected{
		arg: arg,
	}
}

func (e *ErrOneExpected) Error() string {
	return fmt.Sprintf("expected one package, got %d", e.arg)
}

func (e *ErrOneExpected) Unwrap() error {
	return e.err
}

func (e *ErrOneExpected) UnwrapAll() error {
	return grr.UnwrapAll(e)
}

func (e *ErrOneExpected) AsGrr(err grr.Error) (grr.Error, bool) {
	return grr.AsGrr(e, err)
}

func (e *ErrOneExpected) AddTrait(trait grr.Trait, value string) grr.Error {
	e.traits[trait] = value
	return e
}

func (e *ErrOneExpected) GetTrait(key grr.Trait) (string, bool) {
	trait, ok := e.traits[key]
	return trait, ok
}

func (e *ErrOneExpected) GetTraits() map[grr.Trait]string {
	traits := make(map[grr.Trait]string)
	for k, v := range e.traits {
		traits[k] = v
	}
	return traits
}

func (e *ErrOneExpected) AddOp(op string) grr.Error {
	e.op = op
	return e
}

func (e *ErrOneExpected) GetOp() string {
	return e.op
}

func (e *ErrOneExpected) AddError(err error) grr.Error {
	e.err = err
	return e
}

func (e *ErrOneExpected) Trace() {
	grr.Trace(e)
}
