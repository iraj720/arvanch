package request

import (
	"arvanch/i18n"
	"arvanch/pkg/locale"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type SMS struct {
	PhoneNumber string        `json:"phone_number"   validate:"required,phone_number,max=100"`
	Payload     string        `json:"payload"        validate:"required,payload,max=100"`
	Locale      locale.Locale `json:"locale"         validate:"omitempty,locale"`
}

func (r SMS) Validate(reqValidator *validator.Validate) error {
	if err := reqValidator.Struct(r); err != nil {
		return unwrapErrors(err)
	}

	if !i18n.MatchRegionRegexp([]string{"arvan", "turkey"}, r.PhoneNumber) {
		return fmt.Errorf("recipient format is not valid [recipient: %s]", r.PhoneNumber)
	}

	return nil
}
