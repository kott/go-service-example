package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/kott/go-service-example/pkg/db"
	articles "github.com/kott/go-service-example/pkg/services/articles/transport"
	"github.com/kott/go-service-example/pkg/utils/log"
	"github.com/kott/go-service-example/pkg/utils/middleware"
)

// Config defines what the API requires to run
type Config struct {
	DBHost       string
	DBPort       int
	DBUser       string
	DBPassword   string
	DBName       string
	RunMigration bool

	AppHost string
	AppPort int
}

// Start initializes the API server, adding the reuired middleware and dependent services
func Start(cfg *Config) {
	ctx := context.Background()
	conn, err := db.GetConnection(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName)
	if err != nil {
		log.Error(ctx, "unable to establish a database connection: %s", err.Error())
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	if cfg.RunMigration && conn != nil {
		if err := db.Migrate(conn, cfg.DBName); err != nil {
			log.Error(ctx, "unable to complete auto migration", err.Error())
		}
	}

	router := gin.New()

	router.Use(middleware.PersistContext())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.ForceJSON())
	router.Use(middleware.Recover())

	router.NoRoute(middleware.NoRoute())
	router.NoMethod(middleware.NoMethod())

	articles.Activate(router, conn)

	if err := router.Run(fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)); err != nil {
		log.Fatal(context.Background(), err.Error())
	}
}
