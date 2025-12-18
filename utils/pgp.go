package utils

import (
	"encoding/base64"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
)

type PGPHandler struct {
	pgp *crypto.PGPHandle
}

func NewPGPHandler() *PGPHandler {
	pgp := crypto.PGPWithProfile(profile.RFC9580())
	return &PGPHandler{
		pgp: pgp,
	}
}

func (h *PGPHandler) Encrypt(msg string, publicKey *crypto.Key) (string, error) {
	handle, err := h.pgp.Encryption().Recipient(publicKey).New()
	if err != nil {
		return "", err
	}
	pgpMessage, err := handle.Encrypt([]byte(msg))
	if err != nil {
		return "", err
	}
	bytes, err := pgpMessage.ArmorBytes()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (h *PGPHandler) Decrypt(msg string, key *crypto.Key) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}
	decHandle, err := h.pgp.Decryption().DecryptionKey(key).New()
	if err != nil {
		return "", err
	}
	decMsg, err := decHandle.Decrypt(raw, crypto.Armor)
	if err != nil {
		return "", err
	}
	return string(decMsg.Bytes()), nil
}
func (h *PGPHandler) GenerateKey(pgp *crypto.PGPHandle) (*crypto.Key, error) {
	key, err := pgp.KeyGeneration().New().GenerateKey()
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (h *PGPHandler) ParsePublicKey(armoredPubKey string) (*crypto.Key, error) {
	pubKey, err := crypto.NewKeyFromArmored(armoredPubKey)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}
