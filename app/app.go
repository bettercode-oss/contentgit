package app

import (
	"contentgit/app/datasource"
	"contentgit/config"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	gormDB            *gorm.DB
	gin               *gin.Engine
	router            GinRoute
	dbConnector       datasource.DatabaseConnector
	componentRegistry *ComponentRegistry
	logger            *zap.Logger
}

func NewApp(router GinRoute, dbConnector datasource.DatabaseConnector, registry *ComponentRegistry) *App {
	if strings.EqualFold(os.Getenv("CONFIGOR_ENV"), "production") {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.Default()
	if g.SetTrustedProxies(nil) != nil {
		log.Fatal("Failed to set trusted proxies")
	}
	return &App{gin: g, router: router, dbConnector: dbConnector, componentRegistry: registry}
}

func (a *App) SetUp() error {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
		return err
	}
	a.logger = logger

	db, err := a.dbConnector.Connect()
	if err != nil {
		return err
	}
	a.gormDB = db

	if err := a.migrateDatabase(); err != nil {
		return err
	}

	err = a.registerComponents()
	if err != nil {
		return err
	}

	a.subscribeToEvents()

	a.addGinMiddlewares()
	a.router.MapRoutes(a.componentRegistry, a.gin.Group("/api"))

	return nil
}

func (a *App) Run() error {
	if err := a.SetUp(); err != nil {
		return err
	}
	sqlDB, err := a.gormDB.DB()
	if err != nil {
		return err
	}

	defer func() {
		fmt.Println("Closing database connection")
		sqlDB.Close()
	}()

	return a.gin.Run(fmt.Sprintf(":%s", config.Config.HttpPort))
}

func (a *App) GetGin() *gin.Engine {
	return a.gin
}

func (a *App) GetDB() *gorm.DB {
	return a.gormDB
}
