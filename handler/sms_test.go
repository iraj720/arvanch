package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"arvanch/config"
	"arvanch/db"
	"arvanch/i18n"
	"arvanch/repository"
	"arvanch/request"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

var regionWhiteList = []string{"arvan"}

func TestToOTP(t *testing.T) {
	tcs := []struct {
		providers []string
		expected  []string
		name      string
	}{
		{
			providers: []string{"fake"},
			expected:  []string{"fake"},
			name:      "fake",
		},
		{
			providers: []string{"magfa", "rahyab"},
			expected:  []string{"rahyab_otp", "magfa", "rahyab"},
			name:      "magfa-rahyab",
		},
		{
			providers: []string{"rahyab", "magfa"},
			expected:  []string{"rahyab_otp", "rahyab", "magfa"},
			name:      "rahyab-magfa",
		},
		{
			providers: []string{"rahyab"},
			expected:  []string{"rahyab_otp", "rahyab"},
			name:      "rahyab",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// assert.Equal(t, tc.expected, ToOTP(tc.providers, i18n.ara))
		})
	}
}

const (
	helloSMSTemplate = "hello"
	defaultCellphone = "09193041055"
)

type FakeProducer struct {
	rabbitFail bool
}

func (f *FakeProducer) Send(_ context.Context) error {
	if f.rabbitFail {
		return errors.New("some error")
	}

	return nil
}

type FakeSMSPriorityCache struct {
	priority []string
	cacheErr bool
}

type fakeUserAuthClient struct {
	shouldFail bool
}

func (f fakeUserAuthClient) PhoneNumber(_ context.Context, recipient, _ string) (string, error) {
	if f.shouldFail {
		return "", errors.New("fake user auth client failed")
	}

	return recipient, nil
}

type SMSTestSuite struct {
	suite.Suite
	Producer      *FakeProducer
	PriorityCache *FakeSMSPriorityCache
	engine        *echo.Echo
	reqValidator  *validator.Validate
}

func (suite *SMSTestSuite) SetupSuite() {
	suite.engine = echo.New()
	suite.Producer = &FakeProducer{}

	g := suite.engine.Group("api")

	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// c.Set("client", suite.Client)
			return next(c)
		}
	})

	var err error
	suite.reqValidator, err = request.NewValidator()
	suite.NoError(err)

	repo := repository.NewMessageRepo(db.WithRetry(db.Create, config.Init().Postgres))

	g.POST("/sms/phone", NewSMSHandler(repo, i18n.Turkey, nil, suite.reqValidator).Sms)
}

// nolint:funlen,gocognit
func (suite *SMSTestSuite) TestSMSHandler() {
	tcs := []struct {
		name            string
		request         request.SMS
		providers       []string
		providerHeader  string
		status          int
		expectedPayload string
		isOTP           bool
		errExpected     bool
	}{
		{
			name: "successful send",
			request: request.SMS{
				PhoneNumber: "09375080734",
				Payload:     "Hello World",
			},
			providers: []string{"fake"},
			status:    http.StatusOK,
		},
		{
			name: "successful send with new line",
			request: request.SMS{
				PhoneNumber: "09375080734",
				Payload:     " \n Hello World \n\r ",
			},
			expectedPayload: "Hello World",
			providers:       []string{"fake"},
			status:          http.StatusOK,
		},
		{
			name: "successful send with forced provider",
			request: request.SMS{
				PhoneNumber: "09375080734",
				Payload:     "Hello World",
			},
			providers:      []string{"fake1"},
			providerHeader: "fake2",
			status:         http.StatusOK,
		},
		{
			name: "no provider is registered",
			request: request.SMS{
				PhoneNumber: "09375080734",
				Payload:     "Hello World",
			},
			providers: []string{},
			status:    http.StatusInternalServerError,
		},
		{
			name: "rabbit error",
			request: request.SMS{
				PhoneNumber: "09375080734",
				Payload:     "the producer must fail to produce this message because of error",
			},
			providers: []string{"fake"},
			status:    http.StatusInternalServerError,
		},
		{
			name: "bad number format 1",
			request: request.SMS{
				PhoneNumber: "19390909540",
				Payload:     "Hello World",
			},
			providers: []string{"fake"},
			status:    http.StatusBadRequest,
		},
		{
			name: "bad number format 2",
			request: request.SMS{
				PhoneNumber: "9390909540",
				Payload:     "Hello World",
			},
			providers: []string{"fake"},
			status:    http.StatusBadRequest,
		},
		{
			name: "bad number format 3",
			request: request.SMS{
				PhoneNumber: "093909095401",
				Payload:     "Hello World",
			},
			providers: []string{"fake"},
			status:    http.StatusBadRequest,
		},
		{
			name: "empty payload 1",
			request: request.SMS{
				PhoneNumber: "09390909540",
				Payload:     "",
			},
			providers: []string{"fake"},
			status:    http.StatusBadRequest,
		},
		{
			name: "empty payload 2",
			request: request.SMS{
				PhoneNumber: "09375080734",
			},
			providers: []string{"fake"},
			status:    http.StatusBadRequest,
		},
		{
			name: "empty recipient 1",
			request: request.SMS{
				PhoneNumber: "",
				Payload:     "Hello World",
			},
			providers: []string{"fake"},
			status:    http.StatusBadRequest,
		},
		{
			name: "empty recipient 2",
			request: request.SMS{
				Payload: "Hello World",
			},
			providers: []string{"fake"},
			status:    http.StatusBadRequest,
		},
		{
			name:      "empty payload and recipient",
			request:   request.SMS{},
			providers: []string{"fake"},
			status:    http.StatusBadRequest,
		},
		{
			name: "parameterized payload",
			request: request.SMS{
				PhoneNumber: "09390909540",
				Payload:     "hello",
			},
			providers:       []string{"fake"},
			status:          http.StatusOK,
			expectedPayload: "Hello raha",
		},
		{
			name: "otp provider",
			request: request.SMS{
				PhoneNumber: "09390909540",
				Payload:     "Critical Hello",
			},
			providers: []string{"rahyab"},
			status:    http.StatusOK,
			isOTP:     true,
		},
	}

	for i := range tcs {
		tc := tcs[i]

		suite.Run(tc.name, func() {
			if tc.status == http.StatusInternalServerError {
				suite.Producer.rabbitFail = true
			}

			testSms := func(url string) {
				data, err := json.Marshal(tc.request)
				suite.NoError(err)

				w := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

				suite.NoError(err)

				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				if tc.providerHeader != "" {
					req.Header.Set("X-Provider", tc.providerHeader)
				}

				suite.engine.ServeHTTP(w, req)
				suite.Equal(tc.status, w.Code, map[string]string{"url": url})

				// nolint:nestif
				if tc.status == http.StatusOK {
					if tc.expectedPayload == "" {
						suite.Equal(tc.request.Payload, "")
					} else {
						suite.Equal(tc.expectedPayload, "")
					}
				}
			}

			if tc.isOTP {
				testSms("/api/otp")
			} else {
				testSms("/api/sms")

				testSms("/api/sms/phone")
			}

			suite.Producer.rabbitFail = false
		})
	}
}

// nolint:funlen
func (suite *SMSTestSuite) TestSMSPhoneNumberWhiteList() {
	engine := echo.New()
	g := engine.Group("api")

	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// c.Set("client", suite.Client)

			return next(c)
		}
	})

	defaultReq := request.SMS{
		PhoneNumber: "+989193041055",
		Payload:     "Hi",
	}

	defaultErrMsg := "recipient is not whitelisted: +989193041055"

	// The fake user client returns the exact input as the phone number.
	// We need to pass a valid phone number as user ID in related cases.
	// So the white list process won't be disturbed.

	tcs := []struct {
		status    int
		hasError  bool
		name      string
		errMsg    string
		whiteList []string
		req       request.SMS
	}{
		{
			name:      "empty whitelist in test environment",
			whiteList: []string{},
			req:       defaultReq,
			status:    http.StatusForbidden,
			hasError:  true,
			errMsg:    defaultErrMsg,
		},
		{
			name:      "empty whitelist in non-test environment",
			whiteList: []string{},
			req:       defaultReq,
			status:    http.StatusOK,
			hasError:  false,
		},
		{
			name:      "forbidden in non-test environment",
			whiteList: []string{"+989121111111"},
			req:       defaultReq,
			status:    http.StatusOK,
			hasError:  false,
		},
		{
			name:      "whitelisted in test environment 1",
			whiteList: []string{"+989193041055"},
			req:       defaultReq,
			status:    http.StatusOK,
			hasError:  false,
		},
		{
			name:      "whitelisted in test environment 2",
			whiteList: []string{"09193041055"},
			req:       defaultReq,
			status:    http.StatusOK,
			hasError:  false,
		},
		{
			name:      "whitelisted in test environment with user id",
			whiteList: []string{"+989193041055"},
			req: request.SMS{
				PhoneNumber: "09193041055",
				Payload:     "Hi",
			},
			status:   http.StatusOK,
			hasError: false,
		},
		{
			name:      "forbidden in test environment",
			whiteList: []string{"+989121111111"},
			req:       defaultReq,
			status:    http.StatusForbidden,
			hasError:  true,
			errMsg:    defaultErrMsg,
		},
		{
			name:      "forbidden in test environment with user id",
			whiteList: []string{"+989121111111"},
			req: request.SMS{
				PhoneNumber: "09193041055",
				Payload:     "Hi",
			},
			status:   http.StatusForbidden,
			hasError: true,
			errMsg:   "recipient is not whitelisted: 09193041055",
		},
	}

	for i := range tcs {
		tc := tcs[i]

		suite.Run(tc.name, func() {
			smsHandler := NewSMSHandler(nil, i18n.Turkey, nil, suite.reqValidator)

			g.POST("/sms/phone", smsHandler.Sms)

			testFunc := func(url string) {
				data, err := json.Marshal(tc.req)
				suite.NoError(err)

				w := httptest.NewRecorder()
				httpReq := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(data))

				httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				engine.ServeHTTP(w, httpReq)
				suite.Equal(tc.status, w.Code)

				if tc.hasError {
					var resp map[string]interface{}

					body, _ := io.ReadAll(w.Body)

					suite.NoError(json.Unmarshal(body, &resp))
					suite.Equal(tc.errMsg, resp["message"])
				}
			}

			testFunc("/api/sms")

			testFunc("/api/sms/phone")
		})
	}
}

func TestSMS(t *testing.T) {
	suite.Run(t, new(SMSTestSuite))
}
