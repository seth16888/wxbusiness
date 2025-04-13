package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidateErr struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value any    `json:"value"`
	Msg   string `json:"error"`
}

// string
func (v *ValidateErr) String() string {
	return fmt.Sprintf("{\"%s\": \"%v\",\"rule\": \"%s\"}", v.Field, v.Value, v.Tag)
}

type ValidateError struct {
	errs []ValidateErr
}

func (v *ValidateError) Error() string {
	errMsgs := make([]string, 0, len(v.errs))
	for _, err := range v.errs {
		errMsgs = append(errMsgs, err.String())
	}
	return strings.Join(errMsgs, ",")
}

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// Validate
func (v *Validator) Validate(data any) *ValidateError {
	rt := []ValidateErr{}

	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ValidateErr
			elem.Field = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Msg = err.Error()
			rt = append(rt, elem)
		}
	}

	if len(rt) == 0 {
		return nil
	}

	return &ValidateError{
		errs: rt,
	}
}

// ToString
func (v *Validator) ToString(errs []ValidateErr) string {
	errMsgs := make([]string, 0)
	for _, err := range errs {
		errMsgs = append(errMsgs, err.String())
	}
	return strings.Join(errMsgs, ",")
}
