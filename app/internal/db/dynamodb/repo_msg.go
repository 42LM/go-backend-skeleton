package dynamodb

import (
	"context"
	"log/slog"
)

// MsgRepository persists messages in DynamoDB.
type MsgRepository struct {
	Client    *DynamoDBClient
	TableName string
}

// NewMsgRepository constructs a message repository.
func NewMsgRepository(
	client *DynamoDBClient,
	tableName string,
	logger *slog.Logger,
) *MsgRepository {
	return &MsgRepository{client, tableName}
}

// Delete removes a message.
func (r *MsgRepository) Delete(ctx context.Context, id string) error {
	err := r.Client.DeleteItem(ctx, r.TableName, map[string]any{"id": id})
	return handleError(err)
}

// Find loads the message and returns it.
func (r *MsgRepository) Find(
	ctx context.Context,
	id string,
) string {
	var dynamoItem dynamoMsg
	err := r.Client.GetItem(ctx, r.TableName, map[string]any{"id": id}, &dynamoItem)
	if err != nil {
		return ""
	}

	return dynamoItem.Msg
}

// Put creates or updates a message.
func (r *MsgRepository) Put(
	ctx context.Context,
	id,
	msg string,
) error {
	item := &dynamoMsg{
		ID:  id,
		Msg: msg,
	}
	err := r.Client.PutItem(ctx, r.TableName, item, nil)
	if err != nil {
		return err
	}
	return nil
}

type dynamoMsg struct {
	ID  string `dynamodbav:"id"`
	Msg string `dynamodbav:"msg"`
}
