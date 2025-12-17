package main

import (
	"log"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
	"github.com/hunderaweke/sma-go/config"
	"github.com/hunderaweke/sma-go/utils"
)

func main() {
	pgp := crypto.PGPWithProfile(profile.RFC9580())
	key, err := utils.GenerateKey(pgp)
	if err != nil {
		log.Fatal(err)
	}
	pubKey, _ := key.GetArmoredPublicKey()

	publicKey, _ := utils.ParsePublicKey(pubKey)
	msg := `

	Lottery Quick Pick is perhaps the Internet's most popular with over 280 lotteries
	Keno Quick Pick for the popular game played in many countries
	Coin Flipper will give you heads or tails in many currencies
	Dice Roller does exactly what it says on the tin
	Playing Card Shuffler will draw cards from multiple shuffled decks
	Birdie Fund Generator will create birdie holes for golf courses
		`
	encrpted, err := utils.Encrypt(msg, publicKey, pgp)
	if err != nil {
		log.Fatal(err)
	}
	decrypted, err := utils.Decrypt(encrpted, key, pgp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(decrypted)
	config.GenerateSampleEnv()
}
