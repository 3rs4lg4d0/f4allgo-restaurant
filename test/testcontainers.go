package test

import (
	"context"
	"path/filepath"
	"time"

	"github.com/integralist/go-findroot/find"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// InitApplicationDatabase initializes a Postgres container with the application database.
// The database is migrated using the proper creational sql scripts and populated with test
// data.
func InitApplicationDatabase(ctx context.Context) (*postgres.PostgresContainer, error) {
	root, _ := find.Repo()
	return postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithInitScripts(
			filepath.Join(root.Path, "sql/000001_create_schema.up.sql"),
			filepath.Join(root.Path, "sql/000002_add_outbox.up.sql"),
			filepath.Join(root.Path, "test/test_data.sql"),
		),
		postgres.WithDatabase("dbname"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
}

func TerminateApplicationDatabase(ctx context.Context, c *postgres.PostgresContainer) error {
	return c.Terminate(ctx)
}
