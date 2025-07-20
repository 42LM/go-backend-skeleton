package dynamodb_test

import (
	"log/slog"
	"os"

	"go-backend-skeleton/app/internal/db/dynamodb"
)

var (
	dbClient *dynamodb.DynamoDBClient
	dbLogger *slog.Logger
)

func init() {
	dbLogger = slog.Default()
	// dbLogger = zerolog.New(os.Stdout)

	var err error
	dbClient, err = dynamodb.NewDynamoDBClient(
		os.Getenv("DATABASE_AWS_DYNAMODB_ENDPOINT"),
		os.Getenv("AWS_REGION"),
		dbLogger,
	)
	if err != nil {
		panic(err)
	}
}
