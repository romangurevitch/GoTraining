package testify

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ResourceTestSuite shows how to manage a shared external resource across all
// tests in a suite. The temp file is created once in SetupSuite and deleted
// once in TearDownSuite — not re-created for every test.
//
// Use this pattern for expensive resources: database connections, network
// listeners, or any setup that should be shared (not duplicated) across tests.
type ResourceTestSuite struct {
	suite.Suite
	tmpFile *os.File
}

func (s *ResourceTestSuite) SetupSuite() {
	var err error
	s.tmpFile, err = os.CreateTemp("", "testify-suite-*.txt")
	require.NoError(s.T(), err, "shared temp file must be created before suite tests run")
}

func (s *ResourceTestSuite) TearDownSuite() {
	if s.tmpFile != nil {
		name := s.tmpFile.Name()
		s.tmpFile.Close()
		os.Remove(name)
	}
}

func (s *ResourceTestSuite) TestWriteToSharedFile() {
	_, err := s.tmpFile.WriteString("hello")
	assert.NoError(s.T(), err)
}

func (s *ResourceTestSuite) TestFileIsStillOpen() {
	// Because SetupSuite runs once, the same file handle is available across tests.
	assert.NotNil(s.T(), s.tmpFile)
	info, err := s.tmpFile.Stat()
	require.NoError(s.T(), err)
	assert.False(s.T(), info.IsDir())
}

func TestResourceTestSuite(t *testing.T) {
	suite.Run(t, new(ResourceTestSuite))
}
