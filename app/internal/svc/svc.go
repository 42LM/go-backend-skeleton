// Package svc holds all services in packages underneath.
// It provides functionality to connect repos with service and logging.
package svc

import (
	"fmt"
	"log/slog"
	"os"

	"go-backend-skeleton/app/internal/db/dynamodb"
	"go-backend-skeleton/app/internal/db/none"
	"go-backend-skeleton/app/internal/logging"
	"go-backend-skeleton/app/internal/logging/loggingmsg"
	"go-backend-skeleton/app/internal/logging/loggingnone"
	"go-backend-skeleton/app/internal/svc/svcmsg"
	"go-backend-skeleton/app/internal/svc/svcnone"
	"go-backend-skeleton/app/internal/transport"
	"go-backend-skeleton/app/internal/transport/http/httpnone"
)

// MakeService creates service wrapped with logging.
func MakeService(
	logger *logging.LoggerWrapper,
) (
	httpnone.NoneSvc,
	transport.MsgSvc,
	error,
) {
	dbRepoLogger := logger.Logger.With("layer", "database")
	dbLogger := logging.New(dbRepoLogger)
	svcLogger := logging.New(logger.Logger.With("layer", "service"))

	dynamodbClient, err := makeDynamoDBClient(dbRepoLogger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create dynamodb client: %w", err)
	}

	noneRepo := none.NewNoneRepository()
	loggedNoneRepo := loggingnone.NewLoggingRepo(noneRepo, dbLogger)
	noneSvc := svcnone.New(&svcnone.NoneSvcConfig{
		NoneRepo: loggedNoneRepo,
	})
	loggedNoneSvc := loggingnone.NewLoggingSvc(noneSvc, svcLogger)

	msgRepo := dynamodb.NewMsgRepository(dynamodbClient, os.Getenv("DATABASE_AWS_DYNAMODB_MSG_TABLE"))
	loggedMsgRepo := loggingmsg.NewLoggingRepo(msgRepo, dbLogger)
	msgSvc := svcmsg.New(&svcmsg.MsgSvcConfig{
		MsgRepo: loggedMsgRepo,
	})
	loggedMsgSvc := loggingmsg.NewLoggingSvc(msgSvc, svcLogger)

	return loggedNoneSvc, loggedMsgSvc, nil
}

func makeDynamoDBClient(l *slog.Logger) (*dynamodb.DynamoDBClient, error) {
	return dynamodb.NewDynamoDBClient(
		os.Getenv("DATABASE_AWS_DYNAMODB_ENDPOINT"),
		os.Getenv("AWS_REGION"),
		l,
	)
}
