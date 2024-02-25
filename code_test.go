package fserr

import (
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type TestCodeSuite struct {
	suite.Suite
}

func (s *TestCodeSuite) SetupTest() {
	NewInternalError(ErrBasic, "basic error")
	NewBadRequest(100102, "customer error")
	SetDefault(http.StatusInternalServerError, 100101, "default")
}

func (s *TestCodeSuite) TestBuiltinErrCode() {
	errCode := getCode(ErrBasic)
	s.Equal(http.StatusInternalServerError, errCode.HttpCode)
	s.Equal("basic error", errCode.Message)
}

func (s *TestCodeSuite) TestDefaultErrCode() {
	s.Equal(http.StatusInternalServerError, defaultErrCode.HttpCode)
	s.Equal("default", defaultErrCode.Message)
}

func (s *TestCodeSuite) TestCustomerErrCode() {
	errCode := getCode(100102)
	s.Equal(http.StatusBadRequest, errCode.HttpCode)
	s.Equal("customer error", errCode.Message)
}

func TestCode(t *testing.T) {
	suite.Run(t, &TestCodeSuite{})
}
