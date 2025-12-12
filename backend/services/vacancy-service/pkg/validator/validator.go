package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	v *validator.Validate
}

func New() *Validator {
	return &Validator{
		v: validator.New(),
	}
}

func (val *Validator) ValidateStruct(s any) error {
	return val.v.Struct(s)
}
