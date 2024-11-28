package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabaseContainer struct {
	dbContainer  testcontainers.Container
	DatabaseName string
	Username     string
	Password     string
	Port         string
	Host         string
}

func NewTestDatabaseContainer() (*TestDatabaseContainer, error) {
	databaseName := "test_content_git"
	username := "postgres"
	password := "password"

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	ctx := context.Background()
	postgresContainer, err := postgres.Run(
		ctx,
		"bettercode2016/pg16-pgmq",
		postgres.WithInitScripts(filepath.Join(currentDir, "create_database.sql")),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	return &TestDatabaseContainer{
		dbContainer:  postgresContainer,
		DatabaseName: databaseName,
		Username:     username,
		Password:     password,
		Port:         port.Port(),
		Host:         "127.0.0.1",
	}, nil
}

func (c TestDatabaseContainer) Terminate() error {
	return c.dbContainer.Terminate(context.Background())
}

func (c *TestDatabaseContainer) ResetContentEventsQueue() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host,
		c.Username,
		c.Password,
		c.DatabaseName,
		c.Port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		DO $$
		BEGIN
			IF EXISTS (SELECT 1 FROM pgmq.list_queues() WHERE queue_name = 'content') THEN
				PERFORM pgmq.drop_queue('content');
			END IF;
		END $$;

		SELECT pgmq.create('content');
	`)
	if err != nil {
		return fmt.Errorf("failed to reset queue: %w", err)
	}

	return nil
}
