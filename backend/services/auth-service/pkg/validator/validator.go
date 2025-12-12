package validator

import (
    "fmt"
    "github.com/go-playground/validator/v10"
)

type Validator struct {
    validate *validator.Validate
}

func NewValidator() *Validator {
    return &Validator{
        validate: validator.New(),
    }
}

func (v *Validator) Validate(i interface{}) error {
    if err := v.validate.Struct(i); err != nil {
        return v.formatValidationError(err)
    }
    return nil
}

func (v *Validator) formatValidationError(err error) error {
    if validationErrs, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrs {
            switch e.Tag() {
            case "required":
                return fmt.Errorf("field '%s' is required", e.Field())
            case "email":
                return fmt.Errorf("field '%s' must be a valid email", e.Field())
            case "min":
                return fmt.Errorf("field '%s' must be at least %s characters", e.Field(), e.Param())
            case "oneof":
                return fmt.Errorf("field '%s' must be one of: %s", e.Field(), e.Param())
            default:
                return fmt.Errorf("field '%s' validation failed on '%s'", e.Field(), e.Tag())
            }
        }
    }
    return err
}
