package fserr

import (
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type TestCenterSuite struct {
	suite.Suite
	expectedBusinessCode int
}

func (s *TestCenterSuite) SetupTest() {
	s.expectedBusinessCode = 1001
	register(NewOK(s.expectedBusinessCode, "ok"))
}

func (s *TestCenterSuite) TestGetCenterErrCode() {
	s.Equal(s.expectedBusinessCode, getCode(s.expectedBusinessCode).BusinessCode)
}

func (s *TestCenterSuite) TestDefaultErrCode() {
	notExistErrCode := getCode(2001)
	s.Equal("", notExistErrCode.Message)
	s.Equal(http.StatusOK, defaultErrCode.HttpCode)
	s.Equal(2001, notExistErrCode.BusinessCode)
}

func TestCenter(t *testing.T) {
	suite.Run(t, &TestCenterSuite{})
}
