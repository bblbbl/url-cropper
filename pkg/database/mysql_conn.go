package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"path/filepath"
	"runtime"
	"sync"
	"urls/pkg/etc"
)

const driver = "mysql"

var (
	connection *sqlx.DB
	mysqlOnce  sync.Once
)

func GetConnection() *sqlx.DB {
	mysqlOnce.Do(func() {
		connection = initialise()
	})

	return connection
}

func initialise() *sqlx.DB {
	cnf := etc.GetConfig()
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?tls=skip-verify&autocommit=true",
		cnf.Database.User,
		cnf.Database.Password,
		cnf.Database.Host,
		cnf.Database.Port,
		cnf.Database.Database,
	)

	conn, err := sqlx.Connect(driver, connStr)
	if err != nil {
		etc.GetLogger().Fatalf("failed get db connection: %e\n", err)
	}

	if err = conn.Ping(); err != nil {
		etc.GetLogger().Fatalf("failed get make mysql ping: %e\n", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	d, err := mysql.WithInstance(conn.DB, &mysql.Config{})
	if err != nil {
		etc.GetLogger().Fatalf("failed get db driver: %e\n", err)
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		getMigrationsPath(),
		cnf.Database.Database,
		d,
	)

	if err != nil {
		etc.GetLogger().Fatalf("failed to load migrations: %e\n", err)
	}

	if err = migrations.Up(); err != nil {
		if err.Error() != "no change" {
			etc.GetLogger().Fatalf("failed to apply migrations: %e\n", err)
		}
	}

	return conn
}

func CloseMysqlConnection() {
	err := connection.Close()
	if err != nil {
		etc.GetLogger().Fatalf("failed to close mysql connection: %e\n", err)
	}
}

func getMigrationsPath() string {
	_, f, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(f), "../..")

	return fmt.Sprintf("file:///%s/migrations", rootPath)
}
