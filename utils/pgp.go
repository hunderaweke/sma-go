package utils

import (
	"encoding/base64"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func Encrypt(msg string, key *crypto.Key, pgp *crypto.PGPHandle) (string, error) {
	publicKey, err := key.ToPublic()
	if err != nil {
		return "", err
	}
	handle, err := pgp.Encryption().Recipient(publicKey).New()
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

func Decrypt(msg string, key *crypto.Key, pgp *crypto.PGPHandle) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}
	decHandle, err := pgp.Decryption().DecryptionKey(key).New()
	if err != nil {
		return "", err
	}
	decMsg, err := decHandle.Decrypt(raw, crypto.Armor)
	if err != nil {
		return "", err
	}
	return string(decMsg.Bytes()), nil
}
func GenerateKey(pgp *crypto.PGPHandle) (*crypto.Key, error) {
	key, err := pgp.KeyGeneration().New().GenerateKey()
	if err != nil {
		return nil, err
	}
	return key, nil
}

func ParsePublicKey(armoredPubKey string) (*crypto.Key, error) {
	pubKey, err := crypto.NewKeyFromArmored(armoredPubKey)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}
