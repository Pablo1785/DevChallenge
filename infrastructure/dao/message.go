package dao

import (
	"DevChallenge/model"
	"context"
)

type MessageDao interface {
	FindMessageRecipients(ctx context.Context, message *model.Message) (*model.MessageRecipients, error)
}
