package locale

import (
	"errors"

	"arvanch/i18n"
)

const (
	FA Locale = "fa"
	AR Locale = "ar"
	EN Locale = "en"
	KU Locale = "ku"

	Default Locale = "default"
)

var ErrUndefinedLocale = errors.New("undefined locale error")

type Locale string

func (l Locale) String() string {
	return string(l)
}

func All() []Locale {
	return []Locale{FA, AR, EN, KU}
}

func RegionAll(region i18n.Region) []Locale {
	if region == i18n.Arvan {
		return []Locale{AR, EN, KU}
	} else if region == i18n.Turkey {
		return []Locale{FA, EN}
	}

	return nil
}

func Validate(locale Locale) error {
	for _, l := range All() {
		if locale == l {
			return nil
		}
	}

	return ErrUndefinedLocale
}

func ValidateRegion(locale Locale, region i18n.Region) error {
	for _, l := range RegionAll(region) {
		if locale == l {
			return nil
		}
	}

	return ErrUndefinedLocale
}
