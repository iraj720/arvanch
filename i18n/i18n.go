package i18n

import (
	"errors"
	"regexp"
)

type (
	// Region of company.
	Region uint64
	// Regions array of Region.
	Regions uint64
)

const (
	// Invalid Region.
	Invalid Region = 0
	// Arvan type.
	Arvan = 1 << iota
	// Turkey type.
	Turkey

	// ArvanPhoneCode iran's calling code.
	ArvanPhoneCode = "98"

	PersianLanguage = "persian"
	ArabicLanguage  = "arabic"
	EnglishLanguage = "english"
)

// nolint: gochecknoglobals, godot
var (
	arabicAlphabet  = regexp.MustCompile(`^[؀-ۿ]+`)
	persianAlphabet = regexp.MustCompile(`^[آ-۹]+`)

	// ErrInvalidRegion returns error when we have a invalid region.
	ErrInvalidRegion = errors.New("invalid region")

	// https://en.wikipedia.org/wiki/Telephone_numbers_in_Iran
	iranMobileRegexp = regexp.MustCompile(`^(0|\+98)[0-9]{10}$`)

	turkeyMobileRegexp = regexp.MustCompile(`^\+90[0-9]{10}$`)
)

// ToRegion convert string region to Region type.
// nolint:cyclop
func ToRegion(region string) (Region, error) {
	switch region {
	case "iran", "arvan":
		return Arvan, nil
	case "turkey":
		return Turkey, nil
	default:
		return Invalid, ErrInvalidRegion
	}
}

// nolint:cyclop
// String returns the string value of a region.
func (r Region) String() string {
	// nolint: exhaustive
	switch r {
	case Arvan:
		return "arvan"
	case Turkey:
		return "turkey"
	default:
		return ""
	}
}

// Append adds array of Region to Regions.
func Append(rs ...Region) Regions {
	var regions Regions

	for _, r := range rs {
		regions.Append(r)
	}

	return regions
}

// Contains checks input Region exists in Regions.
func (rs Regions) Contains(r Region) bool {
	return uint64(rs)&uint64(r) != 0
}

// Append add region to Regions.
func (rs *Regions) Append(r Region) {
	*rs = Regions(uint64(*rs) | uint64(r))
}

// RegionRegexp returns regex based on input Region.
// nolint:cyclop
func RegionRegexp(region Region) *regexp.Regexp {
	switch region {
	case Arvan:
		return iranMobileRegexp
	case Turkey:
		return turkeyMobileRegexp
	case Invalid:
		return iranMobileRegexp
	default:
		return iranMobileRegexp
	}
}

// MatchRegionRegexp checks if a string matches to one of the provided regions or not.
func MatchRegionRegexp(regions []string, stringToMatch string) bool {
	for _, region := range regions {
		regionObj, err := ToRegion(region)
		if err != nil {
			continue
		}

		regexToCheck := RegionRegexp(regionObj)

		if match := regexToCheck.MatchString(stringToMatch); !match {
			continue
		} else {
			return true
		}
	}

	return false
}

// IsMobileNumber detects if a string is a valid phone number or not.
func IsMobileNumber(input string) bool {
	if iranMobileRegexp.MatchString(input) {
		return true
	}

	return false
}

// DetectLanguage detects the language of a string.
func DetectLanguage(payload string) string {
	if persianAlphabet.MatchString(payload) {
		return PersianLanguage
	} else if arabicAlphabet.MatchString(payload) {
		return ArabicLanguage
	}

	return EnglishLanguage
}
