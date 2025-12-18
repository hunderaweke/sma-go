package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PGPUtilsTestSuite struct {
	suite.Suite
}

func TestPGPUtiles(t *testing.T) {
	suite.Run(t, new(PGPUtilsTestSuite))
}
func (s *PGPUtilsTestSuite) TestEncrypt(){

}
