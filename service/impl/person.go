package impl

import (
	"DevChallenge/infrastructure/dao"
	"DevChallenge/model"
	"DevChallenge/service"
	"context"
)

type personService struct {
	personDao dao.PersonDao
}

func (ps personService) Add(ctx context.Context, person *model.Person) (*model.Person, error) {
	return ps.personDao.Create(ctx, person)
}

func (ps personService) CreateOrUpdateTrustConnections(ctx context.Context, personId string, trustConnections []model.TrustConnection) error {
	//TODO implement me
	panic("implement me")
}

func NewPersonService(personDao dao.PersonDao) service.PersonService {
	return &personService{personDao: personDao}
}
