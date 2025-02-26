package dal

import (
	"database/sql"

	"github.com/highonsemicolon/aura/src/models"
)

type RelationshipRepository struct {
	*MySQLRepository[models.Relationship]
}

func NewRelationshipRepository(db *sql.DB, tableName string) *RelationshipRepository {
	return &RelationshipRepository{NewMySQLRepository[models.Relationship](db, tableName)}
}
