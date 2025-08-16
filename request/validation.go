package request

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"

	"arvanch/pkg/locale"

	"github.com/go-playground/validator/v10"
)

var (
	accountRegex   = regexp.MustCompile("^[a-zA-Z_.]+$")
	recipientRegex = regexp.MustCompile("^[0-9]+$")
)

const (
	minPayloadCharacterLen = 1
	maxPayloadCharacterLen = 100

	minPayloadByteLen = 1
	maxPayloadByteLen = 100
)

func NewValidator() (*validator.Validate, error) {
	reqValidator := validator.New()

	// and represents `account` validator.
	validations := map[string]func(fl validator.FieldLevel) bool{
		"phone_number": phoneNumberValidation,
		"payload":      payloadValidation,
		"account":      accountValidation,
		"locale":       localeValid,
	}

	for name, validationFunc := range validations {
		if err := reqValidator.RegisterValidation(name, validationFunc); err != nil {
			return nil, err
		}
	}

	return reqValidator, nil
}

// localeValid checks the validity of the locale and represents `locale` validator.
func localeValid(fl validator.FieldLevel) bool {
	return locale.Validate(locale.Locale(fl.Field().String())) == nil
}

// accountValidation checks the validity of the account and represents `account` validator.
func accountValidation(fl validator.FieldLevel) bool {
	return accountRegex.MatchString(fl.Field().String())
}

// accountValidation checks the validity of the account and represents `account` validator.
func phoneNumberValidation(fl validator.FieldLevel) bool {
	return recipientRegex.MatchString(fl.Field().String())
}

// payloadValidation checks the validity of the payload and represents `payload` validator.
func payloadValidation(fl validator.FieldLevel) bool {
	payload := fl.Field().String()

	return utf8.RuneCountInString(payload) >= minPayloadCharacterLen &&
		utf8.RuneCountInString(payload) <= maxPayloadCharacterLen &&
		len(payload) >= minPayloadByteLen && len(payload) <= maxPayloadByteLen
}

// nolint:err113
func unwrapErrors(err error) error {
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) || len(validationErrors) == 0 {
		return errors.New("unknown error")
	}

	errorsStr := make([]string, len(validationErrors))
	for i := range validationErrors {
		errorsStr[i] = validationErrors[i].Field()
	}

	return fmt.Errorf("invalid parameters: %v", errorsStr)
}
