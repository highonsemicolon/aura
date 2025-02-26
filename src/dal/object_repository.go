package dal

import (
	"github.com/highonsemicolon/aura/src/models"
)

type ObjectRepository struct {
	*MySQLRepository[models.Object]
}

func NewObjectRepository(dal Database, tableName string) *ObjectRepository {
	return &ObjectRepository{NewMySQLRepository[models.Object](dal, tableName)}
}
