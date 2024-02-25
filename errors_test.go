package fserr

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type TestErrorsSuite struct {
	suite.Suite
	originFundamental *fundamental
	originStack       *withStack
	originMessage     *withMessage
	originCode        *withCode
}

func (s *TestErrorsSuite) SetupTest() {
	s.originFundamental = &fundamental{
		msg:   "origin fundamental",
		stack: callers(),
	}
	s.originStack = &withStack{
		error: s.originFundamental,
		stack: callers(),
	}
	s.originMessage = &withMessage{
		cause: s.originStack,
		msg:   "origin message",
	}
	s.originCode = &withCode{
		cause:        s.originMessage,
		Msg:          "with code",
		HttpCode:     http.StatusOK,
		BusinessCode: ErrBasic,
	}
}

func (s *TestErrorsSuite) TestFundamental() {
	s.Equal("origin fundamental", s.originFundamental.msg)
	s.Equal("origin fundamental", s.originFundamental.Error())
	s.Equal("origin fundamental",
		fmt.Sprintf("%s", s.originFundamental))
	s.Equal(`"origin fundamental"`,
		fmt.Sprintf("%q", s.originFundamental))
}

func (s *TestErrorsSuite) TestStack() {
	s.Equal(s.originFundamental, s.originStack.Cause())
	s.Equal(s.originFundamental, s.originStack.Unwrap())
	s.Equal("origin fundamental", s.originStack.Error())
	s.Equal("origin fundamental", fmt.Sprintf("%s", s.originStack))
	s.Equal(`"origin fundamental"`, fmt.Sprintf("%q", s.originStack))
}

func (s *TestErrorsSuite) TestMessage() {
	s.Equal(s.originStack, s.originMessage.Cause())
	s.Equal(s.originStack, s.originMessage.Unwrap())
	s.Equal("origin message", s.originMessage.Error())
	s.Equal("origin message",
		fmt.Sprintf("%s", s.originMessage))
	s.Equal(`"origin message"`,
		fmt.Sprintf("%q", s.originMessage))
}

func (s *TestErrorsSuite) TestCode() {
	s.Equal(s.originMessage, s.originCode.Cause())
	s.Equal(s.originMessage, s.originCode.Unwrap())
	s.Equal("with code", s.originCode.Error())
	s.Equal("with code", fmt.Sprintf("%s", s.originCode))
	s.Equal(`"with code"`, fmt.Sprintf("%q", s.originCode))
}

func TestErrors(t *testing.T) {
	suite.Run(t, &TestErrorsSuite{})
}
