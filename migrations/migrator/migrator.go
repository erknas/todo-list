package migrator

import (
	"errors"
	"fmt"
	"time"

	"github.com/erknas/todo-list/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func New(cfg *config.Config) error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)

	time.Sleep(time.Second * 3)

	m, err := migrate.New(cfg.Postgres.MigrationPath, connString)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no new migrations")
			return nil
		}
		return err
	}

	fmt.Println("successful migration")

	return nil
}
