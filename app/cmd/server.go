package cmd

import (
	"fmt"
	"net/http"
	"time"

	"go-backend-skeleton/app/internal/logging"
	"go-backend-skeleton/app/internal/svc"
	internalhttp "go-backend-skeleton/app/internal/transport/http"

	"github.com/spf13/cobra"
)

func init() {
	serverCmd.Flags().StringP("port", "p", "8080", "set the port")

	rootCmd.AddCommand(serverCmd)
}

// serverCmd starts a web server, optionally by given port.
// If no port is given the server will run on the default port `:8080`.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Spin up HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			return fmt.Errorf("could not parse `port` flag: %w", err)
		}

		// create a smol schlogger
		logger := logging.NewSlogger()
		// setup sub loggers for database and service layer
		transportLogger := logger.With("layer", "http")

		// create repos and services and connect them with logging
		noneSvc, msgSvc, err := svc.MakeService(logging.New(logger))
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}

		// create the main http handler with the services
		handler := internalhttp.NewHandler(&internalhttp.HandlerConfig{
			NoneSvc: noneSvc,
			MsgSvc:  msgSvc,
			Logger:  transportLogger,
		})

		httpServer := http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		// spin up the server
		err = httpServer.ListenAndServe()
		if err != nil {
			return fmt.Errorf("failed to spin up HTTP server: %w", err)
		}

		return nil
	},
}
