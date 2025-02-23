package repository

import (
	"database/sql"
	"fmt"
	"log"

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

	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	return db
}
