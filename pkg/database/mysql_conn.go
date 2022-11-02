package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"path/filepath"
	"runtime"
	"urls/pkg/config"
)

var connection *sql.DB

func GetConnection() *sql.DB {
	if connection == nil {
		connection = InitConnection()
	}

	return connection
}

func InitConnection() *sql.DB {
	cnf := config.GetConfig()
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=skip-verify&autocommit=true",
		cnf.Database.User,
		cnf.Database.Password,
		cnf.Database.Host,
		cnf.Database.Port,
		cnf.Database.Database,
	)

	conn, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	driver, err := mysql.WithInstance(conn, &mysql.Config{})
	if err != nil {
		log.Fatalf("failed get db driver: %e\n", err)
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		getMigrationsPath(),
		cnf.Database.Database,
		driver,
	)

	if err != nil {
		log.Fatalf("failed to load migrations: %e\n", err)
	}

	if err = migrations.Up(); err != nil {
		if err.Error() != "no change" {
			log.Fatalf("failed to apply migrations: %e\n", err)
		}
	}

	return conn
}

func getMigrationsPath() string {
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../..")

	return fmt.Sprintf("file:///%s/migrations", rootPath)
}
