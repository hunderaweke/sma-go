package main

import (
	"encoding/base64"
	"log"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
)

func main() {
	pgp := crypto.PGPWithProfile(profile.RFC9580())
	key, err := generateKey(pgp)
	if err != nil {
		log.Fatal(err)
	}
	msg := `


Lottery Quick Pick is perhaps the Internet's most popular with over 280 lotteries
Keno Quick Pick for the popular game played in many countries
Coin Flipper will give you heads or tails in many currencies
Dice Roller does exactly what it says on the tin
Playing Card Shuffler will draw cards from multiple shuffled decks
Birdie Fund Generator will create birdie holes for golf courses
	`
	encrpted, err := encrypt(msg, key, pgp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(encrpted)
	decrypted, err := decrypt(encrpted, key, pgp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(decrypted)
}

func encrypt(msg string, key *crypto.Key, pgp *crypto.PGPHandle) (string, error) {
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

func decrypt(msg string, key *crypto.Key, pgp *crypto.PGPHandle) (string, error) {
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
func generateKey(pgp *crypto.PGPHandle) (*crypto.Key, error) {
	key, err := pgp.KeyGeneration().New().GenerateKey()
	if err != nil {
		return nil, err
	}
	return key, nil
}
