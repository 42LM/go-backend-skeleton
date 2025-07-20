package dynamodb_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"go-backend-skeleton/app/internal/db/dynamodb"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MsgRepo_Bootstrap(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.New().String()

	r := buildMsgRepo(t)

	{ // Lookup - nothing to find
		msg := r.Find(ctx, id)
		assert.Empty(t, msg)
	}
	{ // Put msg into db
		err := r.Put(ctx, id, "test-msg-1")
		require.NoError(t, err)
	}
	{ // Find
		msg := r.Find(ctx, id)
		assert.Equal(t, "test-msg-1", msg)
	}
	{ // Delete
		err := r.Delete(ctx, id)
		require.NoError(t, err)
	}
	{ // Find again -> find nothing
		msg := r.Find(ctx, id)
		assert.Empty(t, msg)
	}
}

func Test_MsgRepo_Data_init(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	r := buildMsgRepo(t)

	fixtures := []struct {
		id  string
		msg string
	}{
		{
			"1",
			"a wonderful message",
		},
		{
			"2",
			"a beautiful message",
		},
		{
			"3",
			"a truly meaningful message",
		},
	}

	{ // Put msg into db
		for _, x := range fixtures {
			err := r.Put(ctx, x.id, x.msg)
			require.NoError(t, err)
		}
	}
}

func buildMsgRepo(t *testing.T) *dynamodb.MsgRepository {
	tableName := os.Getenv("DATABASE_AWS_DYNAMODB_MSG_TABLE")
	repo := dynamodb.NewMsgRepository(dbClient, tableName, dbLogger)

	err := dbClient.CreateTable(
		context.Background(),
		tableName,
		dynamodb.SetAttributeDefinitions([]types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		}),
		dynamodb.SetKeySchemas([]types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		}),
		dynamodb.SetProvisionedThroughput(5, 5),
	)
	if err != nil {
		var existsErr *types.ResourceInUseException
		if !errors.As(err, &existsErr) { // okay if the table already exists
			t.Fatalf("failed to create table %q (base endpoint: %v): %v",
				tableName,
				os.Getenv("DATABASE_AWS_DYNAMODB_TOKEN_TABLE"),
				err,
			)
		}
	}

	return repo
}
