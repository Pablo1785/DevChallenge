package dao

import (
	"DevChallenge/model"
	"context"
)

type PersonDao interface {
	Create(ctx context.Context, person *model.Person) (*model.Person, error)
	CreateOrUpdateTrustConnections(ctx context.Context, personId string, trustConnections []model.TrustConnection) error
}
