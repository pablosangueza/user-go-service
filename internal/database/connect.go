package database

import (
	"fmt"
	"net/url"

	"github.com/dom/user/internal/config"
	"github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4"

	"github.com/jmoiron/sqlx"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
)

func OpenDB(cfg *config.AppConf) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", url.PathEscape(cfg.DBUser), url.PathEscape(cfg.DBPassword), cfg.DBHost, cfg.DBPort, cfg.DBName)
	_, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	connFunc := sqlx.Connect

	if cfg.Tracing.Enabled {
		sqltrace.Register("pgx", stdlib.GetDefaultDriver(), sqltrace.WithServiceName("gouser-db"))
		connFunc = sqlxtrace.Connect
	}

	db, err := connFunc("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping database: %w", err)
	}

	return db, nil
}
