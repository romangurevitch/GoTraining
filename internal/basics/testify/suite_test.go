package testify

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExampleTestSuite struct {
	suite.Suite
	SharedValue                   int
	SharedValueDefaultAcrossTests int
}

// Run at the beginning of every suite test functions
func (s *ExampleTestSuite) SetupTest() {
	s.T().Log("SetupTest called")
	s.SharedValue = 100
}

func (s *ExampleTestSuite) SetupSuite() {
	s.T().Log("SetupSuite called")
	s.SharedValue = 100
	s.SharedValueDefaultAcrossTests = 150
}

// Run at the end of every suite test functions (for clean up)
func (s *ExampleTestSuite) TearDownTest() {
	s.T().Log("TearDownTest called")
	s.SharedValue = 0
}

func (s *ExampleTestSuite) TearDownSuite() {
	s.T().Log("TearDownSuite called")
	s.SharedValue = 0
	s.SharedValueDefaultAcrossTests = 0
}

// Run at the end of every suite test functions (for clean up)
func (s *ExampleTestSuite) BeforeTest(suiteName, testName string) {
	s.T().Log("BeforeTest called")
}

// Run at the end of every suite test functions (for clean up)
func (s *ExampleTestSuite) AfterTest(suiteName, testName string) {
	s.T().Log("AfterTest called")
}

// This function runs all the test functions in the suite
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}

// One of the test function in the suite
func (s *ExampleTestSuite) TestExample() {
	s.T().Log("Test 1 running")
	assert.Equal(s.T(), 100, s.SharedValue)
	assert.Equal(s.T(), 150, s.SharedValueDefaultAcrossTests)
}

// Another One of the test function in the suite
func (s *ExampleTestSuite) TestExample1() {
	s.T().Log("Test 2 running")
	assert.Equal(s.T(), 100, s.SharedValue)
	assert.Equal(s.T(), 150, s.SharedValueDefaultAcrossTests)
}
