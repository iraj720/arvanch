package i18n

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsMobileNumber(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "iran valid without country code",
			input:    "09193041055",
			expected: true,
		},
		{
			name:     "iran valid with country code",
			input:    "+989193041055",
			expected: true,
		},
		{
			name:     "iraq valid with country code",
			input:    "+9647773041055",
			expected: true,
		},
		{
			name:     "iran invalid without starting zero",
			input:    "9193041055",
			expected: false,
		},
		{
			name:     "invalid alphabetical",
			input:    "test",
			expected: false,
		},
		{
			name:     "invalid numerical",
			input:    "1234567890",
			expected: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expected, IsMobileNumber(tt.input))
		})
	}
}

func TestRegions(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		regions     Regions
		contains    []Region
		notContains []Region
	}{
		{
			name:        "both regions",
			regions:     Append(Arvan, Turkey),
			contains:    []Region{Arvan, Turkey},
			notContains: []Region{},
		},
		{
			name:        "arvan",
			regions:     Append(Arvan),
			contains:    []Region{Arvan},
			notContains: []Region{Turkey},
		},
		{
			name:        "turkey",
			regions:     Append(Turkey),
			contains:    []Region{Turkey},
			notContains: []Region{Arvan},
		},
		{
			name:        "none",
			regions:     Append(),
			contains:    []Region{},
			notContains: []Region{Arvan, Turkey},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for _, r := range tt.contains {
				require.True(t, tt.regions.Contains(r))
			}
			for _, r := range tt.notContains {
				require.False(t, tt.regions.Contains(r))
			}
		})
	}
}

// nolint:funlen
func TestMatchRegionRegexp(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		phoneNumber string
		expected    bool
		regions     []string
	}{
		{
			name:        "match - one region 1",
			phoneNumber: "+989124567891",
			regions:     []string{"arvan"},
			expected:    true,
		},
		{
			name:        "match - one region 2",
			phoneNumber: "+909124567891",
			regions:     []string{"turkey"},
			expected:    true,
		},
		{
			name:        "not match - one region",
			phoneNumber: "+9047852134560",
			regions:     []string{"arvan"},
			expected:    false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.expected, MatchRegionRegexp(tt.regions, tt.phoneNumber))
		})
	}
}

// nolint:funlen
func TestString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		expected string
		region   Region
	}{
		{
			name:     "arvan",
			region:   Arvan,
			expected: "arvan",
		},
		{
			name:     "turkey",
			region:   Turkey,
			expected: "turkey",
		},
		{
			name:     "invalid 1",
			region:   Invalid,
			expected: "",
		},
		{
			name:     "invalid 2",
			region:   Region(25),
			expected: "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.expected, tt.region.String())
		})
	}
}

// nolint:funlen
func TestRegionRegexp(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		expected *regexp.Regexp
		region   Region
	}{
		{
			name:     "arvan",
			region:   Arvan,
			expected: iranMobileRegexp,
		},
		{
			name:     "turkey",
			region:   Turkey,
			expected: turkeyMobileRegexp,
		},
		{
			name:     "invalid 1",
			region:   Invalid,
			expected: iranMobileRegexp,
		},
		{
			name:     "invalid 2",
			region:   Region(25),
			expected: iranMobileRegexp,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.expected, RegionRegexp(tt.region))
		})
	}
}

// nolint:funlen
func TestToRegion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		region        string
		expected      Region
		expectedError error
	}{
		{
			name:     "arvan",
			expected: Arvan,
			region:   "iran",
		},
		{
			name:          "invalid 1",
			expected:      Invalid,
			region:        "invalid",
			expectedError: ErrInvalidRegion,
		},
		{
			name:          "invalid 2",
			expected:      Invalid,
			region:        "ye esm e ajiib",
			expectedError: ErrInvalidRegion,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp, err := ToRegion(tt.region)

			require.Equal(t, tt.expectedError, err)
			require.Equal(t, tt.expected, resp)
		})
	}
}

func TestDetectLanguage(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		language string
	}{
		{
			name:     "persian 1",
			input:    "آ",
			language: PersianLanguage,
		},
		{
			name:     "persian 2",
			input:    "ی",
			language: PersianLanguage,
		},
		{
			name:     "arabic 1",
			input:    "؀",
			language: ArabicLanguage,
		},
		{
			name:     "arabic 2",
			input:    "ۿ",
			language: ArabicLanguage,
		},
		{
			name:     "english 1",
			input:    "a",
			language: EnglishLanguage,
		},
		{
			name:     "english 2",
			input:    "z",
			language: EnglishLanguage,
		},
		{
			name:     "english persian",
			input:    "سلام arvan",
			language: PersianLanguage,
		},
		{
			name:     "english arabic",
			input:    "؀ arvan",
			language: ArabicLanguage,
		},
		{
			name:     "arabic persian",
			input:    "سلام ؀",
			language: PersianLanguage,
		},
		{
			name:     "arabic persian english",
			input:    "سلام arvan ؀",
			language: PersianLanguage,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, DetectLanguage(tt.input) == tt.language)
		})
	}
}
