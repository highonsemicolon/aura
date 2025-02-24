package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/highonsemicolon/aura/config"
	"github.com/highonsemicolon/aura/src/utils"
)

func InitDB(config config.MySQL) *sql.DB {
	if config.CAPath != "" {
		tlsConfig, err := utils.InitTLS(config.CAPath)
		if err != nil {
			panic(fmt.Errorf("TLS setup failed for CA Cert '%s': %w", config.CAPath, err))
		}

		mysql.RegisterTLSConfig("custom", tlsConfig)
	}

	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	duration, err := time.ParseDuration(config.ConnMaxLifetime)
	if err != nil {
		log.Fatal("Invalid ConnMaxLifetime:", err)
	}
	db.SetConnMaxLifetime(duration)

	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	return db
}

func (r *MySQLRepository[T]) withTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
