package service

import (
	"DevChallenge/model"
	"context"
)

type PersonService interface {
	Add(ctx context.Context, person *model.Person) (*model.Person, error)
	CreateOrUpdateTrustConnections(ctx context.Context, personId string, trustConnections model.TrustConnections) error
}
