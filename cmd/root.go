package cmd

import (
	"os"

	"arvanch/cmd/accounting"
	"arvanch/cmd/messanger"
	"arvanch/cmd/migrate"
	"arvanch/config"
	"arvanch/log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const exitFailure = 1

func Execute() {
	cfg := config.Init()

	var cmd = &cobra.Command{
		Use:   "arvanch",
		Short: "arvanch sends sms and email messages to users",
	}

	log.SetupLogger(cfg.Logger)

	logrus.Debugf("config loaded: %+v", cfg)

	messanger.Register(cmd, cfg)
	accounting.Register(cmd, cfg)
	migrate.Register(cmd, cfg)

	if err := cmd.Execute(); err != nil {
		logrus.Error(err.Error())
		os.Exit(exitFailure)
	}
}
