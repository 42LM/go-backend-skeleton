package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go-backend-skeleton/app/internal/db/none"
	"go-backend-skeleton/app/internal/svc"
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

		svc := svc.New(&svc.ServiceConfig{
			NoneRepo: none.NewNoneRepository(),
			Logger:   slog.Default().With("layer", "service"),
		})
		handler := internalhttp.NewHandler(internalhttp.HandlerConfig{
			Svc: svc,
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
