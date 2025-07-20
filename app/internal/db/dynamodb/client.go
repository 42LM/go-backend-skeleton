package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	// FIXME
	// Xray
	// causes problems locally, for simplification this is omitted
	// "github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"github.com/aws/smithy-go/logging"
)

// DynamoDBClient encapsulates low-level AWS DynamoDB calls.
type DynamoDBClient struct {
	Client *dynamodb.Client
}

type CreateTableOptions struct {
	AttributeDefinitions   []types.AttributeDefinition
	GlobalSecondaryIndexes []types.GlobalSecondaryIndex
	KeySchemas             []types.KeySchemaElement
	ProvisionedThroughput  *types.ProvisionedThroughput
}

type CreateTableOption func(*CreateTableOptions)

func SetAttributeDefinitions(attributeDefinitions []types.AttributeDefinition) CreateTableOption {
	return func(o *CreateTableOptions) {
		o.AttributeDefinitions = attributeDefinitions
	}
}

func SetGlobalSecondaryIndexes(gsIndexes []types.GlobalSecondaryIndex) CreateTableOption {
	return func(o *CreateTableOptions) {
		o.GlobalSecondaryIndexes = gsIndexes
	}
}

func SetKeySchemas(keySchemas []types.KeySchemaElement) CreateTableOption {
	return func(o *CreateTableOptions) {
		o.KeySchemas = keySchemas
	}
}

func SetProvisionedThroughput(readCapacityUnits, writeCapacityUnits int64) CreateTableOption {
	return func(o *CreateTableOptions) {
		o.ProvisionedThroughput = &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(readCapacityUnits),
			WriteCapacityUnits: aws.Int64(writeCapacityUnits),
		}
	}
}

func SetDefault() CreateTableOption {
	return func(o *CreateTableOptions) {
		o.AttributeDefinitions = []types.AttributeDefinition{
			{
				AttributeName: aws.String("hash_key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("range_key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		}
		o.KeySchemas = []types.KeySchemaElement{
			{
				AttributeName: aws.String("hash_key"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("range_key"),
				KeyType:       types.KeyTypeRange,
			},
		}
		o.ProvisionedThroughput = &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		}
	}
}

// NewDynamoDBClient returns a new DynamoDB client.
func NewDynamoDBClient(baseEndpoint, region string, logger *slog.Logger) (*DynamoDBClient, error) {
	dynamoDBLogger := logger.With("database", "dynamodb")

	// Wrap the dynamoDBLogger in a function to satisfy the logger interface.
	wl := logging.LoggerFunc(func(classification logging.Classification, format string, v ...any) {
		dynamoDBLogger.Info("["+string(classification)+"] "+format, v...)
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithLogger(wl),
		config.WithLogConfigurationWarnings(true),
	)
	if err != nil {
		return nil, err
	}

	// Xray
	// awsv2.AWSV2Instrumentor(&cfg.APIOptions)

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if baseEndpoint != "" {
			o.BaseEndpoint = aws.String(baseEndpoint)
		}
		o.Region = region
	})

	return &DynamoDBClient{client}, nil
}

// CreateTable creates a dynamodb table.
func (c *DynamoDBClient) CreateTable(ctx context.Context, tableName string, opts ...CreateTableOption) error {
	options := &CreateTableOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Default initialization fallback when opts are empty
	if len(opts) < 1 {
		SetDefault()(options)
	}

	_, err := c.Client.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions:   options.AttributeDefinitions,
		GlobalSecondaryIndexes: options.GlobalSecondaryIndexes,
		KeySchema:              options.KeySchemas,
		ProvisionedThroughput:  options.ProvisionedThroughput,
		TableName:              aws.String(tableName),
	})
	if err != nil {
		var existsErr *types.ResourceInUseException
		if errors.As(err, &existsErr) { // okay if table already exists
			return nil
		}
		return fmt.Errorf("failed to create table %q: %v", tableName, err)
	}
	return nil
}

// DeleteItem deletes a dynamodb item.
func (c *DynamoDBClient) DeleteItem(ctx context.Context, tableName string, compositePrimaryKey map[string]any) error {
	pkMap := make(map[string]types.AttributeValue, len(compositePrimaryKey))
	for k, v := range compositePrimaryKey {
		marshaledKey, err := attributevalue.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal key %v: %v", v, err)
		}
		pkMap[k] = marshaledKey
	}

	_, err := c.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       pkMap,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete dynamodb item %q: %w", compositePrimaryKey, err)
	}

	return nil
}

// GetItem fetches a dynamodb item.
func (c *DynamoDBClient) GetItem(ctx context.Context, tableName string, compositePrimaryKey map[string]any, data any) error {
	pkMap := make(map[string]types.AttributeValue, len(compositePrimaryKey))
	for k, v := range compositePrimaryKey {
		marshaledKey, err := attributevalue.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal key %v: %v", v, err)
		}
		pkMap[k] = marshaledKey
	}

	resp, err := c.Client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:            pkMap,
		TableName:      aws.String(tableName),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to get dynamodb item %q: %w", compositePrimaryKey, err)
	}

	return attributevalue.UnmarshalMap(resp.Item, &data)
}

// PutItem creates or updates a dynamodb item.
func (c *DynamoDBClient) PutItem(ctx context.Context, tableName string, data any, conditionExpression *string) error {
	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		return fmt.Errorf("failed to marshal item %v: %v", item, err)
	}

	_, err = c.Client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String(tableName),
		ConditionExpression: conditionExpression,
	})
	if err != nil {
		return fmt.Errorf("failed to add item table %q: %w", tableName, err)
	}

	return err
}

// UpdateItem updates a dynamodb item.
func (c *DynamoDBClient) UpdateItem(ctx context.Context, tableName string, compositePrimaryKey map[string]any, updateExpression string, expressionAttributeValues any) error {
	pkMap := make(map[string]types.AttributeValue, len(compositePrimaryKey))
	for k, v := range compositePrimaryKey {
		marshaledKey, err := attributevalue.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal key %v: %v", v, err)
		}
		pkMap[k] = marshaledKey
	}

	marshaledValues, err := attributevalue.MarshalMap(expressionAttributeValues)
	if err != nil {
		return fmt.Errorf("failed to marshal expression attribute values %v: %v", expressionAttributeValues, err)
	}

	input := &dynamodb.UpdateItemInput{
		Key:                       pkMap,
		TableName:                 aws.String(tableName),
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: marshaledValues,
	}

	_, err = c.Client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update DynamoDB item with keys %v: %w", compositePrimaryKey, err)
	}

	return nil
}

// QueryItems queries dyanmodb items.
func (c *DynamoDBClient) QueryItems(ctx context.Context, tableName string, keyConditionExpression string, expressionAttributeValues map[string]any, config *QueryConfig) ([]map[string]types.AttributeValue, error) {
	if config == nil {
		config = &QueryConfig{}
	}
	config.SetDefaults()

	marshaledAttrValues, err := attributevalue.MarshalMap(expressionAttributeValues)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal expression attribute values: %w", err)
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyConditionExpression),
		ExpressionAttributeValues: marshaledAttrValues,
		ScanIndexForward:          aws.Bool(config.ScanIndexForward),
		Limit:                     aws.Int32(config.Limit),
	}

	if config.IndexName != "" {
		queryInput.IndexName = aws.String(config.IndexName)
	}

	resp, err := c.Client.Query(ctx, queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}

	return resp.Items, nil
}

type QueryConfig struct {
	Limit            int32  `json:"limit,omitempty"`
	ScanIndexForward bool   `json:"scanIndexForward"`
	IndexName        string `json:"indexName,omitempty"`
}

func (qc *QueryConfig) SetDefaults() {
	if qc.Limit == 0 {
		qc.Limit = 100
	}
}
