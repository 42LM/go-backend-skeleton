package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go-backend-skeleton/app/internal/db/none"
	"go-backend-skeleton/app/internal/svc/svcnone"
	internalhttp "go-backend-skeleton/app/internal/transport/http"
	"go-backend-skeleton/app/internal/transport/http/httpnone"

	"github.com/spf13/cobra"
)

func init() {
	serverCmd.Flags().StringP("port", "p", "8080", "set the port")

	rootCmd.AddCommand(serverCmd)
}

// TODO: remove and put into logging pkg
// ===== REPOSITORY LOGGING WRAPPER =====

type loggingRepo struct {
	next   svcnone.NoneRepo
	logger *slog.Logger
}

func NewLoggingRepo(next svcnone.NoneRepo, logger *slog.Logger) svcnone.NoneRepo {
	return &loggingRepo{next: next, logger: logger}
}

func (l *loggingRepo) Find(ctx context.Context) string {
	defer func(begin time.Time) {
		l.logger.Info(
			"Find",
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return l.next.Find(ctx)
}

// ===== SERVICE LOGGING WRAPPER =====

type loggingService struct {
	next   httpnone.NoneSvc
	logger *slog.Logger
}

func NewLoggingService(next httpnone.NoneSvc, logger *slog.Logger) httpnone.NoneSvc {
	return &loggingService{next: next, logger: logger}
}

func (l *loggingService) FindNone(ctx context.Context) string {
	defer func(begin time.Time) {
		l.logger.Info(
			"FindNone",
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return l.next.FindNone(ctx)
}

// ===== MAIN APPLICATION =====

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
		loggedRepo := NewLoggingRepo(noneRepo, logger.With("layer", "repo"))

		noneSvc := svcnone.New(&svcnone.NoneSvcConfig{
			NoneRepo: loggedRepo,
		})
		loggedService := NewLoggingService(noneSvc, logger.With("layer", "service"))

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
