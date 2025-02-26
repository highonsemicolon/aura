package dal

import "github.com/highonsemicolon/aura/src/models"

type DalContainer struct {
	Objects       DAL[models.Object]
	relationships DAL[models.Relationship]
}

func NewDalContainer(db Database, tables map[string]string) *DalContainer {
	return &DalContainer{
		Objects:       NewObjectRepository(db, tables["objects"]),
		relationships: NewRelationshipRepository(db, tables["relationships"]),
	}
}
