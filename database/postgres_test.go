package database

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestBuildDSN(t *testing.T) {
	dsn := buildDSN("localhost", "user", "pass", "db", "5432")
	want := "host=localhost user=user password=pass dbname=db port=5432 sslmode=disable"
	if dsn != want {
		t.Fatalf("buildDSN() got %q, want %q", dsn, want)
	}
}

func TestNewPostgresConn_Smoke(t *testing.T) {
	err := godotenv.Load("../.env.test")
	if err != nil {
		t.Error(err)
	}
	_, err = NewPostgresConn()
	if err != nil {
		t.Fatalf("NewPostgresConn() error: %v", err)
	}
}
