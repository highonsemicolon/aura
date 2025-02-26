package service

import (
	"context"

	"github.com/highonsemicolon/aura/src/dal"
	"github.com/highonsemicolon/aura/src/models"
)

type ObjectServiceInterface interface {
	Create(ctx context.Context, user, object string) error
}

type ObjectService struct {
	dal dal.DAL[models.Object]
}

func NewObjectService(dal dal.DAL[models.Object]) *ObjectService {
	return &ObjectService{dal: dal}
}

func (o *ObjectService) Create(ctx context.Context, user, object string) error {
	return o.dal.Create(&models.Object{ID: 123456543, ObjectID: object})
}
