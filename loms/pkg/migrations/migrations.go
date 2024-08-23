package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(connStr string) error {
	m, err := migrate.New(
		"file://migrations",
		connStr,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Up()
}
