package messanger

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"arvanch/config"
	"arvanch/db"
	"arvanch/handler"
	"arvanch/i18n"
	"arvanch/log/access"
	"arvanch/repository"
	"arvanch/request"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const exitTimeout = 5 * time.Second

func Register(root *cobra.Command, cfg config.Config) {
	var port int

	cmd := &cobra.Command{
		Use:   "messanger",
		Short: "Start a new arvanch messanger for accepting sms requests",
		Run: func(cmd *cobra.Command, args []string) {
			main(port, cfg)
		},
	}

	cmd.Flags().IntVar(&port, "port", 0, "port on which messanger will listen to requests")

	if err := cmd.MarkFlagRequired("port"); err != nil {
		logrus.Fatal(err.Error())
	}

	root.AddCommand(cmd)
}

// nolint:funlen
func main(port int, cfg config.Config) {
	region, err := i18n.ToRegion(cfg.I18N.Region)
	if err != nil {
		logrus.Fatalf("messanger : invalid region: %s", err)
	}

	database := db.WithRetry(db.Create, cfg.Postgres)

	defer func() {
		if err := database.Close(); err != nil {
			logrus.Error(err.Error())
		}
	}()

	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/healthz", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	api := e.Group("/api")

	accessLogger, err := access.NewAccessLogger(cfg.CustomAccessLogger)
	if err != nil {
		logrus.Fatalf("messanger : failed to init access logger: %s", err.Error())
	}

	reqValidator, err := request.NewValidator()
	if err != nil {
		logrus.Fatalf("messanger : failed to create validator : %s", err.Error())
	}

	msgRepo := repository.NewMessageRepo(database)

	smsHandler := handler.NewSMSHandler(
		msgRepo,
		region,
		accessLogger,
		reqValidator,
	)

	api.POST("/sms/phone", smsHandler.Sms)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); !errors.Is(err, http.ErrServerClosed) && err != nil {
			e.Logger.Fatal(err.Error())
		}
	}()

	s := <-sig
	logrus.Infof("signal %s received\n", s)

	ctx, cancel := context.WithTimeout(context.Background(), exitTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logrus.Error(err.Error())
	}
}
