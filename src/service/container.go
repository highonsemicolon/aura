package service

import (
	"github.com/highonsemicolon/aura/src/dal"
)

type ServiceContainer struct {
	ObjectService ObjectService
}

func NewServiceContainer(repo *dal.DalContainer) *ServiceContainer {
	return &ServiceContainer{
		ObjectService: NewObjectService(repo.Objects),
	}
}
