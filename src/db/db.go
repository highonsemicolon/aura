package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

type DB interface {
	AssignRole(userID, role, resourceID string) error
	RemoveRole(userID, resourceID string) error
	GetRole(userID, resourceID string) (string, error)
	Close() error
}

type sqlDB struct {
	conn      *sql.DB
	validator *validator.Validate
}

func NewSqlDB(dsn, CACertPath string) *sqlDB {

	if CACertPath != "" {
		if err := setupTLSConfig(CACertPath); err != nil {
			panic(fmt.Errorf("TLS setup failed for CA Cert '%s': %w", CACertPath, err))
		}
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to open DB connection: %w", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping DB: %w", err))
	}
	return &sqlDB{conn: db, validator: validator.New()}
}

func (db *sqlDB) Close() error {
	return db.conn.Close()
}

func (db *sqlDB) validateInput(input interface{}) error {
	if err := db.validator.Struct(input); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	return nil
}

func setupTLSConfig(CACertPath string) error {
	caCert, err := os.ReadFile(CACertPath)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate file at '%s': %w", CACertPath, err)
	}

	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(caCert); !ok {
		return fmt.Errorf("failed to append CA certificate from '%s'", CACertPath)
	}

	tlsConfig := &tls.Config{
		RootCAs: rootCertPool,
	}
	if err := mysql.RegisterTLSConfig("custom", tlsConfig); err != nil {
		return fmt.Errorf("failed to register TLS config: %w", err)
	}
	return nil
}
