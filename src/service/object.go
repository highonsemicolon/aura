package service

import (
	"context"

	"github.com/highonsemicolon/aura/src/dal"
	"github.com/highonsemicolon/aura/src/models"
)

type ObjectService interface {
	Create(ctx context.Context, user, object string) error
}

type objectService struct {
	dal dal.DAL[models.Object]
}

func NewObjectService(dal dal.DAL[models.Object]) *objectService {
	return &objectService{dal: dal}
}

func (o *objectService) Create(ctx context.Context, user, object string) error {
	return o.dal.Create(&models.Object{ID: 123456543, ObjectID: object})
}
