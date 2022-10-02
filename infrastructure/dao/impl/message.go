package impl

import (
	"DevChallenge/infrastructure/dao"
	"DevChallenge/model"
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type messageDao struct {
	DbSession *neo4j.Session
}

func (md *messageDao) FindMessageRecipients(ctx context.Context, message *model.Message) (*model.MessageRecipients, error) {
	result, err := (*md.DbSession).WriteTransaction(func(transaction neo4j.Transaction) (any, error) {
		result, err := transaction.Run(
			`
					match path = (p:Person {id:$fromPersonId})-[rel:TRUSTS*]->(o:Person)
					WHERE ALL(r in rel where r.trust_level > 5)
					unwind nodes(path) as n 
					with n 
					match(p2:Person {id:n.id}) 
					where n.id <> "Garry"
						and all(requiredTopic in $topics where requiredTopic in n.topics)
					return collect(distinct n.id)`,
			map[string]any{"fromPersonId": message.FromPersonId, "topics": message.Topics})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record(), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", result.(*neo4j.Record))
	messageRecipientIds, isKey := result.(*neo4j.Record).Get("messageRecipientIds")
	if !isKey {
		return nil, errors.New("somehow created a record with no properties")
	}

	var messageRecipients model.MessageRecipients
	messageRecipients[message.FromPersonId] = messageRecipientIds.([]string)

	return &messageRecipients, nil
}

func NewMessageDao(dbSession *neo4j.Session) dao.MessageDao {
	return &messageDao{DbSession: dbSession}
}
