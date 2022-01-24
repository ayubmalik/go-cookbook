package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Pharmacy struct {
	Code      string
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
	FindByCode(code string) (*Pharmacy, error)
	FindByPostcode(partial string) ([]*Pharmacy, error)
	Insert(p Pharmacy) (int, error)
}

type PSQLPharmacyRepo struct {
	DB *sql.DB
}

func (r PSQLPharmacyRepo) FindByCode(code string) (*Pharmacy, error) {
	row := r.DB.QueryRow(`select code, name, addr_line_1, addr_line_2, addr_line_3, addr_line_4, postcode
		from pharmacy where code = $1`, code)
	p := &Pharmacy{}
	err := row.Scan(&p.Code, &p.Name, &p.AddrLine1, &p.AddrLine2, &p.AddrLine3, &p.AddrLine4, &p.Postcode)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r PSQLPharmacyRepo) FindByPostcode(partial string) ([]*Pharmacy, error) {
	rows, err := r.DB.Query(`select code, name, addr_line_1, addr_line_2, addr_line_3, addr_line_4, postcode
		from pharmacy where postcode like concat($1::text, '%')`, partial)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pharmacies := make([]*Pharmacy, 0)
	for rows.Next() {
		p := &Pharmacy{}
		err := rows.Scan(&p.Code, &p.Name, &p.AddrLine1, &p.AddrLine2, &p.AddrLine3, &p.AddrLine4, &p.Postcode)

		if err != nil {
			return nil, err
		}
		pharmacies = append(pharmacies, p)
	}

	return pharmacies, nil
}

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("could not get DB connection %s\n", err)
	}
	defer db.Close()

	repo := PSQLPharmacyRepo{
		DB: db,
	}

	// find single
	code := "FA512"
	p, err := repo.FindByCode(code)
	if err != nil {
		log.Fatalf("could not find pharmacy with code %s %v", code, err)
	}
	log.Println(p.Code, p.Name, p.Postcode)

	// find multiple
	postcode := "LE"
	pharmacies, err := repo.FindByPostcode("LE")
	if err != nil {
		log.Fatalf("could not find pharmacies with postcode %s %v", postcode, err)
	}
	log.Println(len(pharmacies))
	log.Println(pharmacies[0].Code, pharmacies[0].Name, pharmacies[0].Postcode)
	log.Println(pharmacies[1].Code, pharmacies[1].Name, pharmacies[1].Postcode)

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
