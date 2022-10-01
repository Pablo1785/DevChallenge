package impl

import (
	"DevChallenge/infrastructure/dao"
	"DevChallenge/model"
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type personDao struct {
	DbSession *neo4j.Session
}

func (pd *personDao) Create(ctx context.Context, person *model.Person) (*model.Person, error) {
	result, err := (*pd.DbSession).WriteTransaction(func(transaction neo4j.Transaction) (any, error) {
		result, err := transaction.Run(
			"CREATE (p:Person) SET p.id = $id, p.topics = $topics RETURN p",
			map[string]any{"id": person.Id, "topics": person.Topics})
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
	props, isKey := result.(*neo4j.Record).Get("properties")
	if !isKey {
		return nil, errors.New("somehow created a record with no properties")
	}
	return &model.Person{
		Id:               props.(map[string]any)["id"].(string),
		Topics:           props.(map[string]any)["topics"].([]string),
		TrustConnections: nil,
	}, nil
}

func (pd *personDao) CreateOrUpdateTrustConnections(ctx context.Context, personId string, trustConnections []model.TrustConnection) error {
	result, err := (*pd.DbSession).WriteTransaction(func(transaction neo4j.Transaction) (any, error) {
		result, err := transaction.Run(
			`UNWIND
						$trustConnections as trustConnections
					MATCH
						(p:Person {id:$pId}),
						(o:Person {id:trustConnections.id})
					CREATE (p)-[rel:TRUSTS {trustConnections.trust_level}]->(o)
					RETURN p, rel, o`,
			map[string]any{"trustConnections": [{id}]})
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
	props, isKey := result.(*neo4j.Record).Get("properties")
	if !isKey {
		return nil, errors.New("somehow created a record with no properties")
	}
	return &model.Person{
		Id:               props.(map[string]any)["id"].(string),
		Topics:           props.(map[string]any)["topics"].([]string),
		TrustConnections: nil,
	}, nil
}

func NewPersonDao(dbSession *neo4j.Session) dao.PersonDao {
	return &personDao{DbSession: dbSession}
}
