package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go-backend-skeleton/app/internal/db/none"
	"go-backend-skeleton/app/internal/logging/loggingnone"
	"go-backend-skeleton/app/internal/svc/svcnone"
	internalhttp "go-backend-skeleton/app/internal/transport/http"

	"github.com/spf13/cobra"
)

func init() {
	serverCmd.Flags().StringP("port", "p", "8080", "set the port")

	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts an HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			return fmt.Errorf("could not parse `port` flag: %w", err)
		}

		logger := slog.Default()

		noneRepo := none.NewNoneRepository()
		loggedRepo := loggingnone.NewLoggingRepo(noneRepo, logger.With("layer", "repo"))

		noneSvc := svcnone.New(&svcnone.NoneSvcConfig{
			NoneRepo: loggedRepo,
		})
		loggedService := loggingnone.NewLoggingService(noneSvc, logger.With("layer", "service"))

		handler := internalhttp.NewHandler(internalhttp.HandlerConfig{
			NoneSvc: loggedService,
		})

		httpServer := http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		fmt.Println("listening on port " + port)
		err = httpServer.ListenAndServe()
		if err != nil {
			return fmt.Errorf("HTTP server failed: %w", err)
		}

		return nil
	},
}
