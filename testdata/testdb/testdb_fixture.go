package testdb

import (
	"fmt"
	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/gorm"
	"path/filepath"
	"runtime"
)

type DatabaseFixture struct {
}

func (DatabaseFixture) SetUpDefault(gormDB *gorm.DB) {
	fmt.Println("Set up database test fixture")
	sqlDB, err := gormDB.DB()
	if err != nil {
		panic(err)
	}

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDB),                                       // You database connection
		testfixtures.Dialect("postgresql"),                                 // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory(filepath.Join(currentDir, "data_fixtures")), // the directory containing the YAML files
	)

	if err != nil {
		panic(err)
	}

	if err := fixtures.Load(); err != nil {
		panic(err)
	}
	fmt.Println("End of database test fixture")
}
