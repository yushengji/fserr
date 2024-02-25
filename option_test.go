package fserr

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestOptionSuite struct {
	suite.Suite
	opt Option
}

func (s *TestOptionSuite) SetupTest() {
	s.opt = WithMessage("cover message")
}

func (s *TestOptionSuite) TestCode() {
	w := &withCode{}
	s.opt(w)
	s.Equal("cover message", w.Msg)
}

func TestOption(t *testing.T) {
	suite.Run(t, &TestOptionSuite{})
}
