package dal

import (
	"fmt"
	"reflect"
)

type MySQLRepository[T any] struct {
	db        Database
	tableName string
}

func NewMySQLRepository[T any](db Database, tableName string) *MySQLRepository[T] {
	if !isValidRepositoryName(tableName) {
		tType := reflect.TypeOf((*T)(nil)).Elem()
		panic(fmt.Sprintf("invalid table name: `%s`, for type: `%s`", tableName, tType))
	}
	return &MySQLRepository[T]{db: db, tableName: tableName}
}

func (r *MySQLRepository[T]) GetByID(id string) (*T, error) {
	query := fmt.Sprintf("SELECT * FROM %v WHERE id = ?", r.tableName)
	row := r.db.QueryRow(query, id)

	var entity T
	err := row.Scan(query, entity)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *MySQLRepository[T]) GetAll() ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %v", r.tableName)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []T
	for rows.Next() {
		var entity T
		err := rows.Scan(entity)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *MySQLRepository[T]) Create(entity *T) error {
	query := fmt.Sprintf("INSERT INTO %v VALUES (?)", r.tableName)
	_, err := r.db.Exec(query, entity)
	return err
}

func (r *MySQLRepository[T]) Update(entity *T) error {
	query := fmt.Sprintf("UPDATE %v SET ? WHERE id = ?", r.tableName)
	_, err := r.db.Exec(query, entity)
	return err
}

func (r *MySQLRepository[T]) Delete(id string) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE id = ?", r.tableName)
	_, err := r.db.Exec(query, id)
	return err
}
