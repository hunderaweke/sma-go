package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PGPUtilsTestSuite struct {
	suite.Suite
	pgpHandler *PGPHandler
}

func (s *PGPUtilsTestSuite) SetupSuite() {
	s.pgpHandler = NewPGPHandler()
}
func TestPGPUtils(t *testing.T) {
	suite.Run(t, new(PGPUtilsTestSuite))
}
func (s *PGPUtilsTestSuite) TestEncrypt() {
	key, err := s.pgpHandler.GenerateKey()
	s.NoError(err)
	msg := "Random message for testing"
	publicKey, err := key.ToPublic()
	s.NoError(err)
	encryptedMsg, err := s.pgpHandler.Encrypt(msg, publicKey)
	s.NoError(err)
	s.NotEqual(msg, encryptedMsg)
	decryptedMsg, err := s.pgpHandler.Decrypt(encryptedMsg, key)
	s.NoError(err)
	s.Equal(msg, decryptedMsg)
}
