package impl

import (
	"DevChallenge/infrastructure/dao"
	"DevChallenge/model"
	"DevChallenge/service"
	"context"
)

type messageService struct {
	messageDao dao.MessageDao
}

func (ms *messageService) Send(ctx context.Context, message *model.Message) (*model.MessageRecipients, error) {
	return ms.messageDao.FindMessageRecipients(ctx, message)
}

func NewMessageService(messageDao dao.MessageDao) service.MessageService {
	return &messageService{messageDao: messageDao}
}
