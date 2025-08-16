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
		smsByPhone      bool
	}{
		{
			name:            "Successful with phone number but without params",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "+989121234567",
				Payload:     "Hi",
			},
			wantErr: false,
		},
		{
			name:            "Successful with locale",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "+989121234567",
				Payload:     "Hi",
				Locale:      locale.FA,
			},
			wantErr: false,
		},
		{
			name:            "Successful with phone number and params",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09121234567",
				Payload:     "shs_template",
			},
			wantErr: false,
		},
		{
			name:            "Successful without params for passenger",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "1",
				Payload:     "Hi",
			},
			wantErr: false,
		},
		{
			name:            "Successful with good recipient type in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "Hi",
			},
			smsByPhone: true,
		},
		{
			name:            "Successful with driver recipient type in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "Hi",
			},
			smsByPhone: true,
		},
		{
			name:            "Successful with passenger recipient type in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "Hi",
			},
			smsByPhone: true,
		},
		{
			name:            "Successful with none recipient type in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "Hi",
			},
			smsByPhone: true,
		},
		{
			name:            "Successful with all recipient type in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "Hi",
			},
			smsByPhone: true,
		},
		{
			name:            "Successful with driver passenger recipient type in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "Hi",
			},
			smsByPhone: true,
		},
		{
			name:            "failed for invalid phone in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "1",
				Payload:     "Hi",
			},
			wantErr:    true,
			smsByPhone: true,
		},
		{
			name:            "failed without recipient type in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "Hi",
			},
			wantErr:    true,
			smsByPhone: true,
		},
		{
			name:            "failed with invalid recipient type in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "Hi",
			},
			wantErr:    true,
			smsByPhone: true,
		},
		{
			name:            "Successful with params for passenger",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "1",
				Payload:     "shs_template",
			},
			wantErr: false,
		},
		{
			name:            "Successful with params for passenger in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "shs_template",
			},
			wantErr:    false,
			smsByPhone: true,
		},
		{
			name:            "Successful without params for driver",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "1",
				Payload:     "Hi",
			},
			wantErr: false,
		},
		{
			name:            "Successful without params for driver in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "Hi",
			},
			wantErr:    false,
			smsByPhone: true,
		},
		{
			name:            "Successful with params for driver",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "1",
				Payload:     "shs_template",
			},
			wantErr: false,
		},
		{
			name:            "Successful with params for driver in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "09022123241",
				Payload:     "shs_template",
			},
			wantErr:    false,
			smsByPhone: true,
		},
		{
			name:            "Fail with invalid mobile number",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "12345678900",
				Payload:     "Hi",
			},
			wantErr: true,
		},
		{
			name:            "fail with Invalid mobile number in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "9151231232",
				Payload:     "Hi",
			},
			wantErr:    true,
			smsByPhone: true,
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
			name:            "Successful with valid locale",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "+989121234567",
				Payload:     "Hi",
				Locale:      locale.FA,
			},
		},
		{
			name:            "Fail with recipient id but without recipient type",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "1",
				Payload:     "Hi",
			},
			wantErr: true,
		},
		{
			name:            "Fail with invalid recipient type",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "1",
				Payload:     "Hi",
			},
			wantErr: true,
		},
		{
			name:            "Fail without recipient id",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				Payload: "Hi",
			},
			wantErr: true,
		},
		{
			name:            "Fail without recipient",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "",
				Payload:     "Hi",
			},
			wantErr: true,
		},
		{
			name:            "Fail without payload",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "+989121234567",
				Payload:     "",
			},
			wantErr: true,
		},
		{
			name:            "Fail without payload in sms by phone",
			regionWhiteList: []string{"arvan"},
			req: request.SMS{
				PhoneNumber: "+989121234567",
			},
			wantErr:    true,
			smsByPhone: true,
		},
		{
			name:            "Successful with phone number but without params in Iraq",
			regionWhiteList: []string{"arvan2"},
			req: request.SMS{
				PhoneNumber: "+9647890123456",
				Payload:     "World is yours",
			},
		},
		{
			name:            "Iranian in Iraq",
			regionWhiteList: []string{"arvan2"},
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
