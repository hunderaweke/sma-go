package main

import (
	"github.com/hunderaweke/sma-go/config"
)

func main() {
	// 	pgp := crypto.PGPWithProfile(profile.RFC9580())
	// 	key, err := generateKey(pgp)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	msg := `

	// Lottery Quick Pick is perhaps the Internet's most popular with over 280 lotteries
	// Keno Quick Pick for the popular game played in many countries
	// Coin Flipper will give you heads or tails in many currencies
	// Dice Roller does exactly what it says on the tin
	// Playing Card Shuffler will draw cards from multiple shuffled decks
	// Birdie Fund Generator will create birdie holes for golf courses
	// 	`
	// 	encrpted, err := encrypt(msg, key, pgp)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println(encrpted)
	// 	decrypted, err := decrypt(encrpted, key, pgp)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println(decrypted)
	config.GenerateSampleEnv()
}
