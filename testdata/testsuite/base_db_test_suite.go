package testsuite

import (
	"contentgit/testdata/testdb"

	"github.com/stretchr/testify/suite"
)

type BaseDatabaseTestSuite struct {
	suite.Suite
	TestDbContainer *testdb.TestDatabaseContainer
}

func (suite *BaseDatabaseTestSuite) SetupSuite() {
	testDbContainer, err := testdb.NewTestDatabaseContainer()
	if err != nil {
		panic(err)
	}
	suite.TestDbContainer = testDbContainer
}

func (suite *BaseDatabaseTestSuite) SetupTest() {
	if err := suite.TestDbContainer.ResetContentEventsQueue(); err != nil {
		panic(err)
	}
}

func (suite *BaseDatabaseTestSuite) TearDownSuite() {
	if err := suite.TestDbContainer.Terminate(); err != nil {
		panic(err)
	}
}
