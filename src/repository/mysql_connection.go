package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/highonsemicolon/aura/utils"
)

func InitDB(dsn, CACertPath string) *sql.DB {
	if CACertPath != "" {
		tlsConfig, err := utils.InitTLS(CACertPath)
		if err != nil {
			panic(fmt.Errorf("TLS setup failed for CA Cert '%s': %w", CACertPath, err))
		}

		mysql.RegisterTLSConfig("custom", tlsConfig)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	return db
}
