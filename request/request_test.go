package request

import "testing"

func TestMarshalRawRequest(t *testing.T) {
	tests := []struct {
		name       string
		expected   string
		rawRequest interface{}
	}{
		{
			name: "not nil request with json tag",
			rawRequest: struct {
				Name string `json:"name"`
				Type string `json:"type"`
			}{
				Name: "Sammy",
				Type: "Shark",
			},
			expected: `{"name":"Sammy","type":"Shark"}`,
		},
		{
			name: "not nil request without json tag",
			rawRequest: struct {
				Name string
				Type string
			}{
				Name: "Sammy",
				Type: "Shark",
			},
			expected: `{"Name":"Sammy","Type":"Shark"}`,
		},
		{
			name:       "nil request",
			rawRequest: nil,
			expected:   ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := MarshalRawRequest(tt.rawRequest)

			if out != tt.expected {
				t.Errorf("marshal raw request failed: expected %s: got: %s", tt.expected, out)
			}
		})
	}
}
