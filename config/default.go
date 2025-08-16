package config

import (
	"time"
)

// nolint:gomnd,mnd,funlen
func Default() Config {
	return Config{
		Token:         "secret",
		ReporterToken: "secret",
		Secret:        "secret",
		Rate:          100,
		SMSRate:       1.0,
		SMSBucket:     5,
		MaskRecipient: true,
		Cache: Cache{
			CronPattern: "0 0/5 * * * *",
		},
		I18N: I18N{
			Region:    "turkey",
			WhiteList: WhiteList{SMS: []string{"arvan"}},
		},
		Postgres: Postgres{
			Host: "localhost",
			Port: 5432,
			// Username:           "postgres",
			Password:           "postgres",
			DBName:             "arvanch",
			ConnectTimeout:     30 * time.Second,
			ConnectionLifetime: 30 * time.Minute,
			MaxOpenConnections: 10,
			MaxIdleConnections: 5,
		},
		Logger: Logger{
			Level: "debug",
		},
		// AccessLogger: log.AccessLogger{
		// 	Enabled: false,
		// },
		DPNLogger: DPNLogger{
			HookEnable:   false,
			StdoutEnable: false,
			Enable:       false,
			SMSPath:      "/logs/sms.log",
			EmailPath:    "/logs/email.log",
			VoicePath:    "/logs/voice.log",
			WhatsappPath: "/logs/whatsapp.log",
		},
		// Redis: redis.Config{
		// 	Address: "redis:6379",
		// },
		CustomAccessLogger: CustomAccessLogger{
			HookEnable:    false,
			StdoutEnable:  false,
			Path:          "/logs/custom-access.log",
			SecurePayload: false,
			EncryptionKey: "secret",
			HMACKey:       "secret",
		},
		NATS: NATS{
			URL:            "127.0.0.1:4222",
			ReconnectWait:  1 * time.Second,
			MaxReconnect:   120,
			PublishEnabled: false,
		},
		Monitoring: Monitoring{
			Prometheus: Prometheus{
				Enabled: true,
				Address: ":9001",
			},
		},
		JTIForOTP:  []string{},
		JTIForBulk: []string{},

		// Gubernator: mineCfg.Gubernator{
		// 	Enabled:     true,
		// 	Timeout:     200 * time.Millisecond,
		// 	GRPCAddress: "localhost:8081",
		// 	AppName:     "arvanch",
		// },
		RateLimits: RateLimits{
			BulkClientsRPS: RateLimitRule{
				Name:      "bulk",
				Duration:  time.Second,
				Limit:     5,
				Algorithm: 1,
				Behaviour: 0,
			},
			ReporterClientsRPS: RateLimitRule{
				Name:      "reporter",
				Duration:  time.Second,
				Limit:     110,
				Algorithm: 1,
				Behaviour: 0,
			},
			RahyabBatch: RateLimitRule{
				Name:      "rahyab",
				Duration:  time.Second,
				Limit:     300,
				Algorithm: 0,
				Behaviour: 0,
			},

			FakeBatch: RateLimitRule{
				Name:      "fake",
				Duration:  time.Second,
				Limit:     300,
				Algorithm: 0,
				Behaviour: 0,
			},
		},
		UserWhiteList: WhiteList{
			SMS: []string{},
		},
	}
}
