package main

import (
	"database/sql"
)

type domain struct {
	ID        int    `json:"id"`
	Domain    string `json:"domain"`
	Operation string `json:"operation"`
}

func (d *domain) getDomain(db *sql.DB) error {
	return db.QueryRow("SELECT domain FROM domains WHERE domain=$1 RETURNING id,domain,operation",
		d.Domain).Scan(&d.Domain, &d.Operation)
}

func (d *domain) getOperation(db *sql.DB) ([]domain, error) {
	rows, err := db.Query("SELECT domain FROM domains WHERE operation = $1",
		d.Operation)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	domains := []domain{}

	for rows.Next() {
		var d domain
		if err := rows.Scan(&d.Domain); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}

	return domains, nil
}

func (d *domain) updateDomain(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE domains SET operation=$2 WHERE domain=$1",
			d.Domain, d.Operation)

	return err
}

func (d *domain) deleteDomain(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM domains WHERE domain=$1", d.Domain)

	return err
}

func (d *domain) createDomain(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO domains(domain, operation) VALUES($1, $2) RETURNING id",
		d.Domain, d.Operation).Scan(d.ID)

	if err != nil {
		return err
	}

	return nil
}

func getDomains(db *sql.DB, count int) ([]domain, error) {
	rows, err := db.Query(
		"SELECT domain, operation FROM domains WHERE operation IS NULL ORDER BY random() LIMIT $1",
		count)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	domains := []domain{}

	for rows.Next() {
		var d domain
		if err := rows.Scan(&d.ID, &d.Domain, &d.Operation); err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}

	return domains, nil
}
