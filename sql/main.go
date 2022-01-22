package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Pharmacy struct {
	Name      string
	AddrLine1 string
	AddrLine2 string
	AddrLine3 string
	AddrLine4 string
	Postcode  string
	LatLng    LatLng
}

type LatLng struct {
	Lat float32
	Lng float32
}

type PharmacyRepo interface {
	FindByID(id string) (Pharmacy, error)
	FindByPostcode(partial string) ([]Pharmacy, error)
	Insert(p Pharmacy) (int, error)
}

type PSQLPharmacyRepo struct {
	DB *sql.DB
}

func (r PSQLPharmacyRepo) FindByID(id string) (Pharmacy, error) {

}

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("could not get DB connection %s\n", err)
	}

	repo := PSQLPharmacyRepo{
		DB: db,
	}

}

func openDB() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefault("DB_PORT", "5432"),
		getEnvOrDefault("DB_USER", "root"),
		getEnvOrDefault("DB_PASSWORD", "password"),
		getEnvOrDefault("DB_NAME", "pharmacy"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getEnvOrDefault(key string, defValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defValue
}
