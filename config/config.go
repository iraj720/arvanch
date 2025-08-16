package config

import (
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"

	"github.com/sirupsen/logrus"
)

const (
	// Namespace is used by prometheus metrics.
	Namespace = "arvanch"

	// Prefix indicates environment variables prefix.
	Prefix = "arvanch_"
)

//nolint:maligned
type (
	Config struct {
		Token         string `koanf:"token"`
		ReporterToken string `koanf:"reporter-token"`
		Secret        string `koanf:"secret"`

		Logger Logger `koanf:"logger"`
		// AccessLogger              log.AccessLogger   `koanf:"access-logger"`
		DPNLogger                 DPNLogger          `koanf:"dpn-logger"`
		CustomAccessLogger        CustomAccessLogger `koanf:"custom-access-logger"`
		InboundCustomAccessLogger CustomAccessLogger `koanf:"inbound-custom-access-logger"`
		NATS                      NATS               `koanf:"nats"`
		BaseAPI                   BaseAPI            `koanf:"base-api"`

		I18N I18N `koanf:"i18n"`

		Cache    Cache    `koanf:"cache"`
		Postgres Postgres `koanf:"postgres"`

		Rate          float64               `koanf:"rate"`
		SMSRate       float64               `koanf:"sms-rate"`
		SMSBucket     int                   `koanf:"sms-bucket"`
		MaskRecipient bool                  `koanf:"mask-recipient"`
		PersiaFava    map[string]PersiaFava `koanf:"persiafava"`

		Vonage map[string]Vonage `koanf:"vonage"`

		Monitoring Monitoring `koanf:"monitoring"`
		JTIForOTP  []string   `koanf:"jti-for-otp"`
		JTIForBulk []string   `koanf:"jti-for-bulk"`

		RateLimits    RateLimits `koanf:"rate-limits"`
		UserWhiteList WhiteList  `koanf:"white-list"`
	}

	I18N struct {
		Region    string    `koanf:"region"`
		WhiteList WhiteList `koanf:"white-list"`
	}

	InboundWebhook struct {
		Name       string        `koanf:"name"`
		APIKey     string        `koanf:"api-key"`
		URL        string        `koanf:"url"`
		Codes      []string      `koanf:"codes"`
		MaxRetries int           `koanf:"max-retries"`
		Timeout    time.Duration `koanf:"timeout"`
	}

	WhiteList struct {
		SMS      []string `koanf:"sms"`
		Email    []string `koanf:"email"`
		Voice    []string `koanf:"voice"`
		Whatsapp []string `koanf:"whatsapp"`
	}

	Logger struct {
		Level string `koanf:"level"`
	}

	DPNLogger struct {
		HookEnable   bool   `koanf:"hook-enable"`
		StdoutEnable bool   `koanf:"stdout-enable"`
		Enable       bool   `koanf:"enable"`
		SMSPath      string `koanf:"sms-path"`
		EmailPath    string `koanf:"email-path"`
		VoicePath    string `koanf:"voice-path"`
		WhatsappPath string `koanf:"whatsapp-path"`
	}

	CustomAccessLogger struct {
		HookEnable    bool   `koanf:"hook-enable"`
		StdoutEnable  bool   `koanf:"stdout-enable"`
		Path          string `koanf:"path"`
		SecurePayload bool   `koanf:"secure-payload"`
		EncryptionKey string `koanf:"encryption-key"`
		HMACKey       string `koanf:"hmac-key"`
	}

	// BaseAPI represents base-api client configurations.
	BaseAPI struct {
		BaseURL string        `koanf:"base-url"`
		Timeout time.Duration `koanf:"timeout"`
	}

	NATS struct {
		URL            string        `koanf:"url"`
		ReconnectWait  time.Duration `koanf:"reconnect-wait"`
		MaxReconnect   int           `koanf:"max-reconnect"`
		PublishEnabled bool          `koanf:"publish-enabled"`
	}

	Cache struct {
		CronPattern string `koanf:"cron-pattern"`
	}

	Postgres struct {
		Host               string        `koanf:"host"`
		Port               int           `koanf:"port"`
		Username           string        `koanf:"user"`
		Password           string        `koanf:"pass"`
		DBName             string        `koanf:"dbname"`
		ConnectTimeout     time.Duration `koanf:"connect-timeout"`
		ConnectionLifetime time.Duration `koanf:"connection-lifetime"`
		MaxOpenConnections int           `koanf:"max-open-connections"`
		MaxIdleConnections int           `koanf:"max-idle-connections"`
	}

	Rabbitmq struct {
		Host           string `koanf:"host"`
		Port           int    `koanf:"port"`
		User           string `koanf:"user"`
		Pass           string `koanf:"pass"`
		RetryThreshold int    `koanf:"retry-threshold"`
	}

	Rahyab struct {
		Number   string `koanf:"number"`
		Username string `koanf:"username"`
		Password string `koanf:"password"`
		Token    string `koanf:"token"`
		Company  string `koanf:"company"`
		URL      string `koanf:"url"`
	}

	RahyabVoice struct {
		APIKey string `koanf:"api-key"`
		URL    string `koanf:"url"`
	}

	Magfa struct {
		Number   string `koanf:"number"`
		Username string `koanf:"username"`
		Password string `koanf:"password"`
		Domain   string `koanf:"domain"`
		URL      string `koanf:"url"`
	}

	AtiyePardaz struct {
		Number string `koanf:"number"`
		APIKey string `koanf:"api-key"`
		URL    string `koanf:"url"`
	}

	Matrix struct {
		URL      string `koanf:"url"`
		UserID   string `koanf:"user-id"`
		RoomID   string `koanf:"room-id"`
		Password string `koanf:"password"`
	}

	Alibaba struct {
		URL             string `koanf:"url"`
		From            string `koanf:"from"`
		RegionID        string `koanf:"region-id"`
		AccessKeyID     string `koanf:"access-key-id"`
		AccessKeySecret string `koanf:"access-key-secret"`
	}

	PersiaFava struct {
		Number   string        `koanf:"number"`
		Username string        `koanf:"username"`
		Password string        `koanf:"password"`
		URL      string        `koanf:"url"`
		Timeout  time.Duration `koanf:"timeout"`
	}

	Irancell struct {
		Address  string        `koanf:"address"`
		Username string        `koanf:"username"`
		Password string        `koanf:"password"`
		URL      string        `koanf:"url"`
		Timeout  time.Duration `koanf:"timeout"`
	}

	Sendinblue struct {
		APIKey         string        `koanf:"api-key"`
		URL            string        `koanf:"url"`
		Domain         string        `koanf:"domain"`
		ConnectTimeout time.Duration `koanf:"connect-timeout"`
	}

	UserAuth struct {
		BaseURL string        `koanf:"base-url"`
		APIKey  string        `koanf:"api-key"`
		Timeout time.Duration `koanf:"timeout"`
	}

	Gringotts struct {
		BaseURL string        `koanf:"base-url"`
		APIKey  string        `koanf:"api-key"`
		Timeout time.Duration `koanf:"timeout"`
	}

	Sparkpost struct {
		Domain string `koanf:"domain"`
		APIKey string `koanf:"api-key"`
	}

	Vonage struct {
		Username  string `koanf:"username"`
		Password  string `koanf:"password"`
		APIKey    string `koanf:"api-key"`
		APISecret string `koanf:"api-secret"`
		BrandName string `koanf:"brand-name"`
	}

	Unifonic struct {
		BaseURL     string        `koanf:"base-url"`
		AppsID      string        `koanf:"appsid"`
		Sender      string        `koanf:"sender"`
		Timeout     time.Duration `koanf:"timeout"`
		DumpEnabled bool          `koanf:"dump-enabled"`
	}

	Cequens struct {
		URL           string        `koanf:"url"`
		DeliveryToken string        `koanf:"delivery-token"`
		Token         string        `koanf:"token"`
		Timeout       time.Duration `koanf:"timeout"`
	}

	Monitoring struct {
		Prometheus Prometheus `koanf:"prometheus"`
	}

	Prometheus struct {
		Enabled bool   `koanf:"enabled"`
		Address string `koanf:"address"`
	}

	RateLimits struct {
		BulkClientsRPS     RateLimitRule `koanf:"bulk-clients-rps"`
		ReporterClientsRPS RateLimitRule `koanf:"reporter-clients-rps"`
		RahyabBatch        RateLimitRule `koanf:"rahyab-batch"`
		FakeBatch          RateLimitRule `koanf:"fake-batch"`

		SMSRules      []TokenBasedRateLimitRule `koanf:"sms-rules"`
		EmailRules    []TokenBasedRateLimitRule `koanf:"email-rules"`
		VoiceRules    []TokenBasedRateLimitRule `koanf:"voice-rules"`
		WhatsappRules []TokenBasedRateLimitRule `koanf:"whatsapp-rules"`
	}

	TokenBasedRateLimitRule struct {
		// mineCfg.RateLimitRule `koanf:",squash"`
		JTI string `koanf:"jti"`
	}

	// RateLimitRule defines a ratelimit rule in gubernator.
	// Name should be unique for each rule record.
	// For Algorithm and Behaviour definition refer to
	// https://github.com/mailgun/gubernator/blob/master/proto/gubernator.proto
	RateLimitRule struct {
		Name      string        `koanf:"name"`
		Duration  time.Duration `koanf:"duration"`
		Limit     int64         `koanf:"limit"`
		Algorithm int32         `koanf:"algorithm"`
		Behaviour int32         `koanf:"behaviour"`
	}
)

func Init() Config {
	var cfg Config

	k := koanf.New(".")

	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		logrus.Fatalf("error loading default: %s", err)
	}

	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		logrus.Errorf("error loading config.yml: %s", err)
	}

	// since koanf does not convert - separated tags in mapstructure to _ automatically.
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		parsedEnv := strings.Replace(strings.ToLower(strings.TrimPrefix(s, Prefix)), "__", "-", -1)
		return strings.Replace(parsedEnv, "_", ".", -1)
	}), nil); err != nil {
		logrus.Errorf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &cfg); err != nil {
		logrus.Fatalf("error unmarshalling config: %s", err)
	}

	return cfg
}
