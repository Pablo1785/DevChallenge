package service

import (
	"DevChallenge/model"
	"context"
)

type MessageService interface {
	Send(ctx context.Context, message *model.Message) (*model.MessageRecipients, error)
}
