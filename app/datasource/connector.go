package datasource

import (
	"contentgit/config"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConnector interface {
	Connect() (*gorm.DB, error)
}

type ProductionDbConnector struct {
}

func (c ProductionDbConnector) Connect() (*gorm.DB, error) {
	dialector := c.createDialector()
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, errors.Wrap(err, "Database Connection Error")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "Database Connection Error")
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}

func (ProductionDbConnector) createDialector() gorm.Dialector {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		config.Config.DataSource.Host,
		config.Config.DataSource.UserName,
		config.Config.DataSource.Password,
		config.Config.DataSource.DatabaseName,
		config.Config.DataSource.Port)
	return postgres.Open(dsn)
}
