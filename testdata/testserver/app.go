package testserver

import (
	"contentgit/app"
	"contentgit/config"
	"contentgit/testdata/testdb"
	"fmt"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestAppServer struct {
	router      app.GinRoute
	registry    *app.ComponentRegistry
	dbContainer *testdb.TestDatabaseContainer
	internalApp *app.App
}

func (t *TestAppServer) run() *TestAppServer {
	if err := config.InitConfig("../../../config"); err != nil {
		panic(err)
	}

	t.internalApp = app.NewApp(t.router, TestDbConnector{
		Host:     t.dbContainer.Host,
		Port:     t.dbContainer.Port,
		Database: t.dbContainer.DatabaseName,
		User:     t.dbContainer.Username,
		Password: t.dbContainer.Password,
	}, t.registry)
	err := t.internalApp.SetUp()
	if err != nil {
		panic(err)
	}

	return t
}

func (t *TestAppServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t.internalApp.GetGin().ServeHTTP(w, req)
}

func (t *TestAppServer) getDB() *gorm.DB {
	return t.internalApp.GetDB()
}

func (t *TestAppServer) setMockComponent(name string, mock any) {
	t.registry.Register(name, mock)
}

func newTestAppServer(router app.GinRoute, dbContainer *testdb.TestDatabaseContainer) *TestAppServer {
	return &TestAppServer{
		router:      router,
		dbContainer: dbContainer,
		registry:    app.NewComponentRegistry(),
	}
}

type TestDbConnector struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func (c TestDbConnector) Connect() (*gorm.DB, error) {
	fmt.Println("Set up database")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		c.Host,
		c.User,
		c.Password,
		c.Database,
		c.Port)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

type TestAppServerBuilder struct {
	appServer *TestAppServer
	dbFixture bool
}

func NewTestAppServerBuilder(router app.GinRoute, dbContainer *testdb.TestDatabaseContainer) *TestAppServerBuilder {
	return &TestAppServerBuilder{
		appServer: newTestAppServer(router, dbContainer),
		dbFixture: false,
	}
}

func (builder *TestAppServerBuilder) WithDatabaseFixture() *TestAppServerBuilder {
	builder.dbFixture = true
	return builder
}

func (builder *TestAppServerBuilder) SetMockComponent(name string, component any) *TestAppServerBuilder {
	builder.appServer.setMockComponent(name, component)
	return builder
}

func (builder *TestAppServerBuilder) Build() *TestAppServer {
	builder.appServer.run()
	if builder.dbFixture {
		testdb.DatabaseFixture{}.SetUpDefault(builder.appServer.getDB())
	}
	return builder.appServer
}
