package db

import (
	"github.com/pranesh/bitespeed/home"
)

func InsertPrimaryContact(email, phone *string) (*home.DSContact, error) {
	row := DB.QueryRow(
		`INSERT INTO contacts (phone_number, email, link_precedence) VALUES ($1, $2, $3) RETURNING `+selectCols,
		phone, email, home.LinkPrimary,
	)
	return scanRow(row)
}

func InsertSecondaryContact(email, phone *string, primaryID int) (*home.DSContact, error) {
	row := DB.QueryRow(
		`INSERT INTO contacts (phone_number, email, linked_id, link_precedence) VALUES ($1, $2, $3, $4) RETURNING `+selectCols,
		phone, email, primaryID, home.LinkSecondary,
	)
	return scanRow(row)
}
