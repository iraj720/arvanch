package handler

import (
	"errors"
	"fmt"
	"net/http"

	"arvanch/i18n"
	"arvanch/log/access"
	"arvanch/model"
	"arvanch/pkg/locale"
	"arvanch/repository"
	"arvanch/request"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	xUserIDHeader = "X-USER-ID"

	SmsPrice = 100
)

type (
	SMSHandler struct {
		msgRepo      repository.MessageRepository
		Region       i18n.Region
		AccessLogger *access.Logger
		reqValidator *validator.Validate
	}
)

func NewSMSHandler(
	msgRepo repository.MessageRepository,
	region i18n.Region,
	accessLogger *access.Logger,
	reqValidator *validator.Validate,
) SMSHandler {
	return SMSHandler{
		msgRepo:      msgRepo,
		Region:       region,
		AccessLogger: accessLogger,
		reqValidator: reqValidator,
	}
}

// nolint:funlen,gocognit,gocyclo
func (s SMSHandler) Sms(c echo.Context) error {
	smsLog := s.setupSMSLog(c)

	defer func() {
		if s.AccessLogger != nil {
			s.AccessLogger.LogSMS(smsLog)
		}
	}()

	// Get the userID header
	userID := c.Request().Header.Get(xUserIDHeader)

	if userID == "" {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Missing %v header", xUserIDHeader))
	}

	msgID := uuid.New().String()
	smsLog.UUID = msgID

	var req request.SMS
	req.Locale = locale.Default

	if err := c.Bind(&req); err != nil {
		smsLog.Payload = request.MarshalRawRequest(req)
		smsLog.Error = fmt.Sprintf("sms handler: parsing body failed: %s", err.Error())

		return c.JSON(http.StatusBadRequest, echo.Map{"message": errors.New("request's body is not valid")})
	}

	smsLog.Payload = request.MarshalRawRequest(req)

	if err := req.Validate(s.reqValidator); err != nil {
		smsLog.Error = fmt.Sprintf("sms handler: validation failed: %s", err.Error())

		return c.JSON(http.StatusBadRequest, echo.Map{"message": errors.New("request's body is not valid")})
	}

	smsLog.Recipient = req.PhoneNumber

	// read from cache
	userProfile, err := s.msgRepo.GetUserProfile(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	err = s.msgRepo.InsertMessage(&model.Message{
		ID:       msgID,
		UserID:   userID,
		Payload:  req.Payload,
		Language: string(req.Locale),
	})

	// TODO : use more specific errors
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	// TODO : use more specific errors
	err = s.msgRepo.IncrementAccountBalance(userProfile.AccountID, -SmsPrice)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

// nolint:funlen,gocognit,gocyclo
func (s SMSHandler) ChargeAccount(c echo.Context) error {
	smsLog := s.setupSMSLog(c)

	defer func() {
		if s.AccessLogger != nil {
			s.AccessLogger.LogSMS(smsLog)
		}
	}()

	// Get the userID header
	userID := c.Request().Header.Get(xUserIDHeader)

	if userID == "" {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Missing %v header", xUserIDHeader))
	}

	var req request.Charge
	if err := c.Bind(&req); err != nil {
		smsLog.Payload = request.MarshalRawRequest(req)
		smsLog.Error = fmt.Sprintf("sms handler: parsing body failed: %s", err.Error())

		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	smsLog.Payload = request.MarshalRawRequest(req)

	if err := req.Validate(s.reqValidator); err != nil {
		smsLog.Error = fmt.Sprintf("sms handler: validation failed: %s", err.Error())

		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	userProfile, err := s.msgRepo.GetUserProfile(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	err = s.msgRepo.IncrementAccountBalance(userProfile.AccountID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// nolint:funlen,gocognit,gocyclo
func (s SMSHandler) GetUserMessages(c echo.Context) error {
	smsLog := s.setupSMSLog(c)

	defer func() {
		if s.AccessLogger != nil {
			s.AccessLogger.LogSMS(smsLog)
		}
	}()

	// Get the userID header
	userID := c.Request().Header.Get(xUserIDHeader)

	if userID == "" {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Missing %v header", xUserIDHeader))
	}

	fmt.Println("hello")

	msgs, err := s.msgRepo.GetUserMessages(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, msgs)
}

// nolint:funlen,gocognit,gocyclo
func (s SMSHandler) CreateAccount(c echo.Context) error {
	smsLog := s.setupSMSLog(c)

	defer func() {
		if s.AccessLogger != nil {
			s.AccessLogger.LogSMS(smsLog)
		}
	}()

	var req request.Account
	if err := c.Bind(&req); err != nil {
		smsLog.Payload = request.MarshalRawRequest(req)
		smsLog.Error = fmt.Sprintf("sms handler: parsing body failed: %s", err.Error())

		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	smsLog.Payload = request.MarshalRawRequest(req)

	if err := req.Validate(s.reqValidator); err != nil {
		smsLog.Error = fmt.Sprintf("sms handler: validation failed: %s", err.Error())

		return c.JSON(http.StatusBadRequest, echo.Map{"message": err})
	}

	userID := uuid.New().String()

	if err := s.msgRepo.InsertUserWithAccount(userID, req.Name); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, &model.User{
		ID: userID,
	})
}

// nolint:funlen,gocognit,gocyclo
func (s SMSHandler) GetProfile(c echo.Context) error {
	smsLog := s.setupSMSLog(c)

	defer func() {
		if s.AccessLogger != nil {
			s.AccessLogger.LogSMS(smsLog)
		}
	}()

	// Get the userID header
	userID := c.Request().Header.Get(xUserIDHeader)

	if userID == "" {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Missing %v header", xUserIDHeader))
	}

	fmt.Println("user ID : ", userID)

	profile, err := s.msgRepo.GetUserProfile(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, profile)
}

func (s SMSHandler) setupSMSLog(c echo.Context) *access.SMSLog {
	return &access.SMSLog{
		XForwardedFor: c.Request().Header.Get(echo.HeaderXForwardedFor),
		XRealIP:       c.Request().Header.Get(echo.HeaderXRealIP),
		RemoteAddress: c.Request().RemoteAddr,
	}
}
