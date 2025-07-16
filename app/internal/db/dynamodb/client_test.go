package dynamodb_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"

	dynamodbclient "go-backend-skeleton/app/internal/db/dynamodb"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	dbClient *dynamodbclient.DynamoDBClient
	dbLogger *slog.Logger
)

func init() {
	dbLogger = slog.New(slog.NewTextHandler(io.Discard, nil))
	// dbLogger = slog.Default()

	var err error
	dbClient, err = dynamodbclient.NewDynamoDBClient(os.Getenv("DATABASE_AWS_DYNAMODB_ENDPOINT"), os.Getenv("AWS_REGION"), dbLogger)
	if err != nil {
		panic(err)
	}
}

type TestItem1 struct {
	HashKey     string `dynamodbav:"hash_key"`
	RangeKey    string `dynamodbav:"range_key"`
	Title       string `dynamodbav:"title"`
	Description string `dynamodbav:"description"`
}

func Test_Bootstrap(t *testing.T) {
	tableName := "test-table"
	ctx := context.Background()

	err := dbClient.CreateTable(context.Background(), tableName)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	testItem := TestItem1{
		HashKey:  "<id>|<entityName>|<entityID>", // example hash key structure
		RangeKey: "<updated_at>|<UID>",           // example range key structure
		Title:    "007",
	}
	err = dbClient.PutItem(ctx, tableName, &testItem, nil)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	updateExpression := "SET description = :description"
	updateAttributes := map[string]any{":description": "test description"}
	err = dbClient.UpdateItem(ctx, tableName, map[string]any{
		"hash_key":  "<id>|<entityName>|<entityID>",
		"range_key": "<updated_at>|<UID>",
	}, updateExpression, updateAttributes)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	data := TestItem1{}
	err = dbClient.GetItem(ctx, tableName, map[string]any{
		"hash_key":  "<id>|<entityName>|<entityID>",
		"range_key": "<updated_at>|<UID>",
	}, &data)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	if data.Title != testItem.Title {
		t.Errorf("%v failed:\nexpected: %v\nactual: %v\n", t.Name(), testItem.Title, data.Title)
	}
	if data.Description != "test description" {
		t.Errorf("%v failed:\nexpected: test description \nactual: %v\n", t.Name(), data.Description)
	}
}

type TestItem2 struct {
	HashKey string `dynamodbav:"hash_key"`
	Number  int    `dynamodbav:"number"`
	Title   string `dynamodbav:"title"`
}

func Test_DeleteItem(t *testing.T) {
	tableName := "test-delete-item2"
	ctx := context.Background()

	err := createCustomTable(
		tableName,
		dynamodbclient.SetAttributeDefinitions([]types.AttributeDefinition{
			{
				AttributeName: aws.String("hash_key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("number"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		}),
		dynamodbclient.SetKeySchemas([]types.KeySchemaElement{
			{
				AttributeName: aws.String("hash_key"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("number"),
				KeyType:       types.KeyTypeRange,
			},
		}),
		dynamodbclient.SetProvisionedThroughput(5, 5),
	)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	{ // Delete item - nothing there - no error
		err = dbClient.DeleteItem(ctx, tableName, map[string]any{
			"hash_key": "<id>|<entityName>|<entityID>",
			"number":   777,
		})
		if err != nil {
			t.Errorf("%v failed:\nexpected returned error to not be nil: %v", t.Name(), err)
		}
	}

	{ // put item
		testItem := TestItem2{
			HashKey: "<id>|<entityName>|<entityID>",
			Number:  777,
			Title:   "007",
		}
		err = dbClient.PutItem(ctx, tableName, &testItem, nil)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}
	}

	{ // find item
		data := TestItem2{}
		err = dbClient.GetItem(ctx, tableName, map[string]any{
			"hash_key": "<id>|<entityName>|<entityID>",
			"number":   777,
		}, &data)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}

		if data.Title != "007" {
			t.Errorf("%v failed:\nexpected: \"007\"\nactual: %v\n", t.Name(), data.Title)
		}
	}

	{ // delete item
		err = dbClient.DeleteItem(ctx, tableName, map[string]any{
			"hash_key": "<id>|<entityName>|<entityID>",
			"number":   777,
		})
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}
	}

	{ // find item again - not found
		data := TestItem2{}
		err = dbClient.GetItem(ctx, tableName, map[string]any{
			"hash_key": "<id>|<entityName>|<entityID>",
			"number":   777,
		}, &data)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}
		if data.Number != 0 {
			t.Errorf("%v failed:\nexpected: 0\nactual: %v\n", t.Name(), data.Title)
		}
	}
}

type TestItem3 struct {
	Binary []byte `dynamodbav:"binary"`
	Title  string `dynamodbav:"title"`
}

func Test_GetItem(t *testing.T) {
	tableName := "test-delete-item3"
	ctx := context.Background()

	err := createCustomTable(
		tableName,
		dynamodbclient.SetAttributeDefinitions([]types.AttributeDefinition{
			{
				AttributeName: aws.String("binary"),
				AttributeType: types.ScalarAttributeTypeB,
			},
		}),
		dynamodbclient.SetKeySchemas([]types.KeySchemaElement{
			{
				AttributeName: aws.String("binary"),
				KeyType:       types.KeyTypeHash,
			},
		}),
		dynamodbclient.SetProvisionedThroughput(5, 5),
	)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	{ // put item
		testItem := TestItem3{
			Binary: []byte{'A'},
			Title:  "007",
		}
		err = dbClient.PutItem(ctx, tableName, &testItem, nil)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}
	}

	{ // find item
		data := TestItem3{}
		err = dbClient.GetItem(ctx, tableName, map[string]any{
			"binary": []byte{'A'},
		}, &data)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}

		if data.Title != "007" {
			t.Errorf("%v failed:\nexpected: \"007\"\nactual: %v\n", t.Name(), data.Title)
		}
	}
}

func Test_PutItem(t *testing.T) {
	type TestPutItem struct {
		HashKey string `dynamodbav:"hash_key"`
		Title   string `dynamodbav:"title"`
	}

	tableName := "test-put-item"
	ctx := context.Background()

	err := createCustomTable(
		tableName,
		dynamodbclient.SetAttributeDefinitions([]types.AttributeDefinition{
			{
				AttributeName: aws.String("hash_key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		}),
		dynamodbclient.SetKeySchemas([]types.KeySchemaElement{
			{
				AttributeName: aws.String("hash_key"),
				KeyType:       types.KeyTypeHash,
			},
		}),
		dynamodbclient.SetProvisionedThroughput(5, 5),
	)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	{ // put item
		testItem := TestPutItem{
			HashKey: "test_key",
			Title:   "Test title",
		}
		err = dbClient.PutItem(ctx, tableName, &testItem, nil)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}
	}
	{
		// put item with condition expression
		conExpr := "attribute_not_exists(hash_key)"
		testItem := TestPutItem{
			HashKey: "test_key",
			Title:   "New title",
		}
		err = dbClient.PutItem(ctx, tableName, &testItem, aws.String(conExpr))
		var apiErr *types.ConditionalCheckFailedException
		if !errors.As(err, &apiErr) {
			t.Errorf("%v failed:\nexpected: \"ConditionalCheckFailedException\"\nactual: %v\n", t.Name(), err)
		}
	}

	{ // find item - should be the first one
		data := TestPutItem{}
		err = dbClient.GetItem(ctx, tableName, map[string]any{
			"hash_key": "test_key",
		}, &data)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}

		if data.Title != "Test title" {
			t.Errorf("%v failed:\nexpected: \"Test title\"\nactual: %v\n", t.Name(), data.Title)
		}
	}
}

func Test_UpdateItem(t *testing.T) {
	type TestUpdateItem struct {
		HashKey     string `dynamodbav:"hash_key"`
		Title       string `dynamodbav:"title"`
		Description string `dynamodbav:"description"`
	}

	tableName := "test-update-item"
	ctx := context.Background()

	err := createCustomTable(
		tableName,
		dynamodbclient.SetAttributeDefinitions([]types.AttributeDefinition{
			{
				AttributeName: aws.String("hash_key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		}),
		dynamodbclient.SetKeySchemas([]types.KeySchemaElement{
			{
				AttributeName: aws.String("hash_key"),
				KeyType:       types.KeyTypeHash,
			},
		}),
		dynamodbclient.SetProvisionedThroughput(5, 5),
	)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	{ // put item
		testItem := TestUpdateItem{
			HashKey:     "test_key",
			Title:       "Test title",
			Description: "Test description",
		}
		err = dbClient.PutItem(ctx, tableName, &testItem, nil)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}
	}
	{
		// update item
		updateExpression := "SET description = :description"
		updateAttributes := map[string]any{":description": "Updated description"}
		err = dbClient.UpdateItem(ctx, tableName, map[string]any{
			"hash_key": "test_key",
		}, updateExpression, updateAttributes)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}
	}

	{ // find item
		data := TestUpdateItem{}
		err = dbClient.GetItem(ctx, tableName, map[string]any{
			"hash_key": "test_key",
		}, &data)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}

		if data.Title != "Test title" {
			t.Errorf("%v failed:\nexpected: \"Test titles\"\nactual: %v\n", t.Name(), data.Title)
		}
		if data.Description != "Updated description" {
			t.Errorf("%v failed:\nexpected: \"updated description\"\nactual: %v\n", t.Name(), data.Description)
		}
	}
}

func Test_QueryItems(t *testing.T) {
	type TestQueryItem struct {
		PartitionKey string `dynamodbav:"partition_key"`
		SortKey      string `dynamodbav:"sort_key"`
		GSIKey       string `dynamodbav:"gsi_key"`
		Data         string `dynamodbav:"data"`
	}

	tableName := "test-query-items-gsi"
	ctx := context.Background()

	err := createCustomTable(
		tableName,
		dynamodbclient.SetAttributeDefinitions([]types.AttributeDefinition{
			{AttributeName: aws.String("partition_key"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("sort_key"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("gsi_key"), AttributeType: types.ScalarAttributeTypeS},
		}),
		dynamodbclient.SetKeySchemas([]types.KeySchemaElement{
			{AttributeName: aws.String("partition_key"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("sort_key"), KeyType: types.KeyTypeRange},
		}),
		dynamodbclient.SetGlobalSecondaryIndexes([]types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("GSI1"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("gsi_key"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("sort_key"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits: aws.Int64(5), WriteCapacityUnits: aws.Int64(5),
				},
			},
		}),
		dynamodbclient.SetProvisionedThroughput(5, 5),
	)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	// Insert items into the table
	items := []TestQueryItem{
		{"key1", "001", "gsi1", "2024-07-30T15:00:00Z"},
		{"key1", "002", "gsi1", "2024-07-30T16:00:00Z"},
		{"key1", "003", "gsi1", "2024-07-30T17:00:00Z"},
		{"key2", "001", "gsi2", "2024-07-30T18:00:00Z"},
	}
	for _, item := range items {
		err = dbClient.PutItem(ctx, tableName, &item, nil)
		if err != nil {
			t.Errorf("%v failed:\n%v", t.Name(), err)
		}
	}

	keyConditionExpression := "gsi_key = :gk"
	expressionAttributeValues := map[string]any{":gk": "gsi1"}
	queryConfig := &dynamodbclient.QueryConfig{
		IndexName: "GSI1",
		Limit:     2,
	}

	// test query by GSI / order by sort key descending / limit 2
	queriedItems, err := dbClient.QueryItems(ctx, tableName, keyConditionExpression, expressionAttributeValues, queryConfig)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	if len(queriedItems) != 2 {
		t.Errorf("%v failed: expected 2 items, got %d", t.Name(), len(queriedItems))
	}
	if queriedItems[0]["data"].(*types.AttributeValueMemberS).Value != items[2].Data {
		t.Errorf("%v failed: expected %s, got %s", t.Name(), items[2].Data, queriedItems[0]["data"].(*types.AttributeValueMemberS).Value)
	}

	// test query by partition key / order by sort key ascending
	keyConditionExpression = "partition_key = :pk"
	expressionAttributeValues = map[string]any{":pk": "key1"}
	queryConfig = &dynamodbclient.QueryConfig{
		ScanIndexForward: true,
	}

	queriedItems, err = dbClient.QueryItems(ctx, tableName, keyConditionExpression, expressionAttributeValues, queryConfig)
	if err != nil {
		t.Errorf("%v failed:\n%v", t.Name(), err)
	}

	if len(queriedItems) != 3 {
		t.Errorf("%v failed: expected 2 items, got %d", t.Name(), len(queriedItems))
	}
	if queriedItems[0]["data"].(*types.AttributeValueMemberS).Value != items[0].Data {
		t.Errorf("%v failed: expected %s, got %s", t.Name(), items[0].Data, queriedItems[0]["data"].(*types.AttributeValueMemberS).Value)
	}
}

func createCustomTable(tableName string, options ...dynamodbclient.CreateTableOption) error {
	// Clean up the table if it already exists
	_, err := dbClient.Client.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		var nfErr *types.ResourceNotFoundException
		if !errors.As(err, &nfErr) {
			return fmt.Errorf("failed to delete table %q: %v", tableName, err)
		}
	}

	err = dbClient.CreateTable(context.Background(), tableName,
		options...,
	)
	if err != nil {
		return err
	}

	return nil
}
