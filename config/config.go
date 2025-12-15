package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
}

var (
	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")

	sampleEnv = map[string]string{
		"DB_PORT":     "your_db_port",
		"DB_NAME":     "your_db_name",
		"DB_PASSWORD": "your_db_password",
		"DB_USERNAME": "your_db_username",
	}
)

func GenerateSampleEnv() error {
	data, err := godotenv.Marshal(sampleEnv)
	if err != nil {
		return err
	}
	file, err := os.Create(".env.sample")
	if err != nil {
		return err
	}
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}
