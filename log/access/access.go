package access

import (
	"io"
	"os"
	"time"

	"arvanch/config"
	"arvanch/log"
	"arvanch/security"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

type (
	SMSLog struct {
		Priorities      []string
		RegionWhiteList []string
		ForcedProvider  string
		Environment     string
		Client          string
		Organization    string
		Venture         string
		UUID            string
		Payload         string
		Recipient       string
		RecipientType   string
		XForwardedFor   string
		XRealIP         string
		RemoteAddress   string
		Language        string
		Error           string
		TemplateName    string
		MessageLength   int
		MessageBytes    int
	}

	Logger struct {
		securePayload    bool
		logger           *logrus.Logger
		payloadEncryptor security.PayloadTransformer
		payloadHMAC      security.PayloadTransformer
	}
)

func NewAccessLogger(cfg config.CustomAccessLogger) (*Logger, error) {
	logrusLogger := logrus.New()

	if cfg.StdoutEnable {
		logrusLogger.SetOutput(os.Stdout)
	} else {
		logrusLogger.SetOutput(io.Discard)
	}

	if cfg.HookEnable {
		rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
			Filename:   cfg.Path,
			MaxSize:    log.MaxSize,
			MaxBackups: log.MaxBackups,
			MaxAge:     log.MaxAge,
			Level:      logrus.InfoLevel,
			Formatter: &logrus.JSONFormatter{
				TimestampFormat:  time.RFC3339,
				DisableTimestamp: false,
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyMsg:  "message",
					logrus.FieldKeyTime: "timestamp",
				},
			},
		})
		if err != nil {
			return nil, err
		}

		logrusLogger.AddHook(rotateFileHook)
	}

	tr, err := security.NewAESTransformer(cfg.EncryptionKey)
	if err != nil {
		return nil, err
	}

	return &Logger{
		logger:           logrusLogger,
		securePayload:    cfg.SecurePayload,
		payloadEncryptor: tr,
		payloadHMAC:      security.NewHMACTransformer(cfg.HMACKey),
	}, nil
}

func (l *Logger) LogSMS(smsLog *SMSLog) {
	payloadEnc, payloadHMAC := "", ""
	recipientEnc, recipientHMAC := "", ""
	payload, recipient := smsLog.Payload, smsLog.Recipient

	if l.securePayload {
		payloadEnc, payloadHMAC = l.secure(smsLog.Payload)
		recipientEnc, recipientHMAC = l.secure(smsLog.Recipient)
		payload, recipient = "", ""
	}

	l.logger.WithFields(logrus.Fields{
		"environment":      smsLog.Environment,
		"client":           smsLog.Client,
		"organization":     smsLog.Organization,
		"venture":          smsLog.Venture,
		"uuid":             smsLog.UUID,
		"priority_list":    smsLog.Priorities,
		"forced_provider":  smsLog.ForcedProvider,
		"payload_enc":      payloadEnc,
		"payload_hmac":     payloadHMAC,
		"recipient_enc":    recipientEnc,
		"recipient_hmac":   recipientHMAC,
		"payload":          payload,
		"recipient":        recipient,
		"recipient_type":   smsLog.RecipientType,
		"secured":          l.securePayload,
		"x_forwarded_for":  smsLog.XForwardedFor,
		"x_real_ip":        smsLog.XRealIP,
		"remote_address":   smsLog.RemoteAddress,
		"error":            smsLog.Error,
		"template_name":    smsLog.TemplateName,
		"message_length":   smsLog.MessageLength,
		"message_bytes":    smsLog.MessageBytes,
		"language":         smsLog.Language,
		"region_whitelist": smsLog.RegionWhiteList,
		"media":            "sms",
	}).Info("sms request received")
}

func (l *Logger) encrypt(field string) string {
	encryptedField, err := l.payloadEncryptor.Transform(field)
	if err != nil {
		logrus.Errorf("failed to encrypt field for logs: %s", err.Error())

		encryptedField = ""
	}

	return encryptedField
}

func (l *Logger) secure(field string) (string, string) {
	encryptedField := l.encrypt(field)

	fieldMac, err := l.payloadHMAC.Transform(field)
	if err != nil {
		logrus.Errorf("failed to calculate mac for logs: %s", err.Error())

		fieldMac = ""
	}

	return encryptedField, fieldMac
}
