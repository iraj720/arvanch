package request

import (
	"github.com/go-playground/validator/v10"
)

type Charge struct {
	Amount int64 `json:"amount"        validate:"required"`
}

func (r Charge) Validate(reqValidator *validator.Validate) error {
	if err := reqValidator.Struct(r); err != nil {
		return unwrapErrors(err)
	}

	return nil
}

type Account struct {
	Name string `json:"name"        validate:"required,max=100"`
}

func (r Account) Validate(reqValidator *validator.Validate) error {
	if err := reqValidator.Struct(r); err != nil {
		return unwrapErrors(err)
	}

	return nil
}
