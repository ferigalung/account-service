package database

import (
	"context"
	"fmt"

	"github.com/ferigalung/account-service/config"
	"github.com/ferigalung/account-service/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnectionPool return postgres db pool
func NewConnectionPool(dbConfig *config.DBConfig) *pgxpool.Pool {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		logger.Log("fatal", "Unable to connect to database", fiber.Map{"error": err.Error()})
		return nil
	}

	return dbPool
}
