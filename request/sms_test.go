package request_test

import (
	"testing"

	"arvanch/pkg/locale"
	"arvanch/request"
)

// nolint:funlen
func TestSMS_Validate(t *testing.T) {
	reqValidator, err := request.NewValidator()
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name            string
		req             request.SMS
		regionWhiteList []string
		wantErr         bool
	}{
		{
			name:            "Successful with phone number",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09121234567",
				Locale:      locale.EN,
				Payload:     "Hi",
			},
		},
		{
			name:            "failed with invalid phone number",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022",
				Locale:      locale.EN,
				Payload:     "Hi",
			},
			wantErr: true,
		},
		{
			name:            "fail with Invalid mobile number 2",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "9151231232",
				Payload:     "Hi",
			},
			wantErr: true,
		},
		{
			name:            "Fail with invalid locale",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "+989121234567",
				Payload:     "Hi",
				Locale:      locale.Locale("I'm a locale"),
			},
			wantErr: true,
		},

		{
			name:            "Fail without recipient",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				Payload: "Hi",
				Locale:  locale.EN,
			},
			wantErr: true,
		},
		{
			name:            "Fail without payload",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "+989121234567",
				Locale:      locale.EN,
			},
			wantErr: true,
		},
		{
			name:            "Iranian in Iraq",
			regionWhiteList: []string{"turkey"},
			req: request.SMS{
				PhoneNumber: "+989121234567",
				Payload:     "World is yours",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.req.Validate(reqValidator); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
