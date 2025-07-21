package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go-backend-skeleton/app/internal/db/dynamodb"
	"go-backend-skeleton/app/internal/db/none"
	"go-backend-skeleton/app/internal/logging"
	"go-backend-skeleton/app/internal/logging/loggingmsg"
	"go-backend-skeleton/app/internal/logging/loggingnone"
	"go-backend-skeleton/app/internal/svc/svcmsg"
	"go-backend-skeleton/app/internal/svc/svcnone"
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
		dbLogger := logger.With("layer", "database")
		svcLogger := logger.With("layer", "service")

		dynamodbClient, err := makeDynamoDBClient(dbLogger)
		if err != nil {
			return fmt.Errorf("failed to create dynamodb client: %w", err)
		}

		// TODO: refactor
		// create repos and services and connect them with logging
		noneRepo := none.NewNoneRepository()
		loggedNoneRepo := loggingnone.NewLoggingRepo(noneRepo, dbLogger)
		noneSvc := svcnone.New(&svcnone.NoneSvcConfig{
			NoneRepo: loggedNoneRepo,
		})
		loggedNoneSvc := loggingnone.NewLoggingSvc(noneSvc, svcLogger)

		msgRepo := dynamodb.NewMsgRepository(dynamodbClient, os.Getenv("DATABASE_AWS_DYNAMODB_MSG_TABLE"), dbLogger)
		loggedMsgRepo := loggingmsg.NewLoggingRepo(msgRepo, dbLogger)
		msgSvc := svcmsg.New(&svcmsg.MsgSvcConfig{
			MsgRepo: loggedMsgRepo,
		})
		loggedMsgSvc := loggingmsg.NewLoggingSvc(msgSvc, svcLogger)

		// create the main http handler with the services
		handler := internalhttp.NewHandler(internalhttp.HandlerConfig{
			NoneSvc: loggedNoneSvc,
			MsgSvc:  loggedMsgSvc,
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

// TODO: refactor
func makeDynamoDBClient(l *slog.Logger) (*dynamodb.DynamoDBClient, error) {
	return dynamodb.NewDynamoDBClient(
		os.Getenv("DATABASE_AWS_DYNAMODB_ENDPOINT"),
		os.Getenv("AWS_REGION"),
		l,
	)
}
