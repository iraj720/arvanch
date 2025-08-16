package log

import "testing"

func TestApplyFilter(t *testing.T) {
	t.Parallel()

	t.Run("replace no word", func(t *testing.T) {
		filterSingleWord := []string{}
		textRequest := "this is a normal request"
		textResponse := "this is a normal response"
		expectedRequest := "this is a normal request"
		expectedResponse := "this is a normal response"

		p := NewProviderLogger(nil, filterSingleWord...)

		actualRequest, actualResponse := p.filterRequestResponse(textRequest, textResponse)

		if actualRequest != expectedRequest || actualResponse != expectedResponse {
			t.Errorf("wrong expected request and/or response")
		}
	})

	t.Run("replace single word", func(t *testing.T) {
		filterSingleWord := []string{"sensitive"}
		textRequest := "this is a sensitive request"
		textResponse := "this is a sensitive response"
		expectedRequest := "this is a " + maskedFilter + " request"
		expectedResponse := "this is a " + maskedFilter + " response"

		p := NewProviderLogger(nil, filterSingleWord...)

		actualRequest, actualResponse := p.filterRequestResponse(textRequest, textResponse)

		if actualRequest != expectedRequest || actualResponse != expectedResponse {
			t.Errorf("wrong expected request and/or response")
		}
	})

	t.Run("replace multi word", func(t *testing.T) {
		filterSingleWord := []string{"sensitive", ":D", ":O"}
		textRequest := "this is a sensitive request :D"
		textResponse := "this :O is a sensitive response"
		expectedRequest := "this is a " + maskedFilter + " request " + maskedFilter
		expectedResponse := "this " + maskedFilter + " is a " + maskedFilter + " response"

		p := NewProviderLogger(nil, filterSingleWord...)

		actualRequest, actualResponse := p.filterRequestResponse(textRequest, textResponse)

		if actualRequest != expectedRequest || actualResponse != expectedResponse {
			t.Errorf("wrong expected request and/or response")
		}
	})
}
