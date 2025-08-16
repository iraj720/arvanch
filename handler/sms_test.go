package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"arvanch/config"
	"arvanch/db"
	"arvanch/i18n"
	"arvanch/pkg/locale"
	"arvanch/repository"
	"arvanch/request"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

const (
	DefaultUserID = "39880004-467e-479d-a1fd-37dce2e76704"
)

type SMSTestSuite struct {
	suite.Suite
	engine       *echo.Echo
	reqValidator *validator.Validate
}

func (suite *SMSTestSuite) SetupSuite() {
	suite.engine = echo.New()

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

	// TODO : should be implemented with mock repo
	repo := repository.NewMessageRepo(db.WithRetry(db.Create, config.Init().Postgres))

	repo.InsertUserWithAccount(DefaultUserID, "test")

	prof, err := repo.GetUserProfile(DefaultUserID)

	repo.IncrementAccountBalance(prof.AccountID, 10000000)

	g.POST("/sms/phone", NewSMSHandler(repo, i18n.Arvan, nil, suite.reqValidator).Sms)
}

// nolint:funlen,gocognit
func (suite *SMSTestSuite) TestSMSHandler() {
	tcs := []struct {
		name            string
		request         request.SMS
		status          int
		expectedPayload string
		errExpected     bool
	}{
		{
			name: "successful send",
			request: request.SMS{
				PhoneNumber: "09375080734",
				Payload:     "Hello World",
				Locale:      locale.EN,
			},
			status: http.StatusCreated,
		},
		{
			name: "bad number format 1",
			request: request.SMS{
				PhoneNumber: "19390909540",
				Payload:     "Hello World",
				Locale:      locale.EN,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "bad number format 2",
			request: request.SMS{
				PhoneNumber: "9390909540",
				Payload:     "Hello World",
				Locale:      locale.EN,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "bad number format 3",
			request: request.SMS{
				PhoneNumber: "093909095401",
				Payload:     "Hello World",
				Locale:      locale.EN,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "empty payload 1",
			request: request.SMS{
				PhoneNumber: "09390909540",
				Payload:     "",
				Locale:      locale.EN,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "empty payload 2",
			request: request.SMS{
				PhoneNumber: "09375080734",
				Locale:      locale.EN,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "empty recipient 1",
			request: request.SMS{
				PhoneNumber: "",
				Payload:     "Hello World",
				Locale:      locale.EN,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "empty recipient 2",
			request: request.SMS{
				Payload: "Hello World",
				Locale:  locale.EN,
			},
			status: http.StatusBadRequest,
		},
		{
			name:    "empty payload and recipient",
			request: request.SMS{},
			status:  http.StatusBadRequest,
		},
	}

	for i := range tcs {
		tc := tcs[i]

		suite.Run(tc.name, func() {
			data, err := json.Marshal(tc.request)
			suite.NoError(err)

			url := "/api/sms/phone"

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))

			req.Header.Set("X-USER-ID", DefaultUserID)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.NoError(err)

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
		})
	}
}

func TestSMS(t *testing.T) {
	suite.Run(t, new(SMSTestSuite))
}
