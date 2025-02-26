package dal

import (
	"github.com/highonsemicolon/aura/src/models"
)

type RelationshipRepository struct {
	*MySQLRepository[models.Relationship]
}

func NewRelationshipRepository(dal *MySQLDAL, tableName string) *RelationshipRepository {
	return &RelationshipRepository{NewMySQLRepository[models.Relationship](dal, tableName)}
}
