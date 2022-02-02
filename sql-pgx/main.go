package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
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
	Insert(p Pharmacy) error
}

type PSQLPharmacyRepo struct {
	Conn *pgx.Conn
	PharmacyRepo
}

func (r PSQLPharmacyRepo) FindByCode(code string) (*Pharmacy, error) {
	sql := `select code, name, addr_line_1, addr_line_2, addr_line_3, addr_line_4, postcode
	from pharmacy where code = $1`

	p := &Pharmacy{}
	row := r.Conn.QueryRow(context.Background(), sql, code)
	err := row.Scan(&p.Code, &p.Name, &p.AddrLine1, &p.AddrLine2, &p.AddrLine3, &p.AddrLine4, &p.Postcode)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r PSQLPharmacyRepo) FindByPostcode(postcode string) ([]*Pharmacy, error) {
	sql := `select code, name, addr_line_1, addr_line_2, addr_line_3, addr_line_4, postcode
	from pharmacy where postcode like concat($1::text, '%')`

	rows, err := r.Conn.Query(context.Background(), sql, postcode)
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

func (r PSQLPharmacyRepo) Insert(p Pharmacy) error {
	sql := `insert into pharmacy(code, name, addr_line_1, addr_line_2, addr_line_3, addr_line_4,
		postcode, lat, lng) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.Conn.Exec(
		context.Background(),
		sql,
		p.Code,
		p.Name,
		p.AddrLine1,
		p.AddrLine2,
		p.AddrLine3,
		p.AddrLine4,
		p.Postcode,
		p.LatLng.Lat,
		p.LatLng.Lng,
	)

	if err != nil {
		return err
	}

	return nil
}

func main() {
	conn, err := openConn()
	if err != nil {
		log.Fatalf("could not get DB connection %s\n", err)
	}
	defer conn.Close(context.Background())

	repo := PSQLPharmacyRepo{
		Conn: conn,
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

	pharmacy := Pharmacy{
		Code:      "NEW1",
		Name:      "A test pharmacy",
		AddrLine1: "line1",
		AddrLine2: "line2",
		AddrLine3: "line3",
		AddrLine4: "line4",
		Postcode:  "postccode",
		LatLng:    LatLng{Lat: 1.1, Lng: 2.2},
	}

	// will error on second run due to PK violation :)
	err = repo.Insert(pharmacy)
	if err != nil {
		log.Fatalf("could not insert pharmacy %v", err)
	}

}

func openConn() (*pgx.Conn, error) {
	url := getEnvOrDefault("DB_URL", "postgres://root:password@localhost:5432/pharmacy")
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func getEnvOrDefault(key string, defValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defValue
}
