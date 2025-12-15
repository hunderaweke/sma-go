package database

import (
	"fmt"

	"github.com/hunderaweke/sma-go/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func buildDSN(host, user, password, dbname, port string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
}
func NewPostgresConn() (*gorm.DB, error) {
	dsn := buildDSN(config.DBHost, config.DBUsername, config.DBPassword, config.DBName, config.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{FullSaveAssociations: true})
	if err != nil {
		return nil, err
	}
	return db, nil
}
