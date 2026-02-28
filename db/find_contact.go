package db

import (
	"database/sql"

	"github.com/pranesh/bitespeed/home"
)

const selectCols = `id, phone_number, email, linked_id, link_precedence, created_at, updated_at, deleted_at`

func FindContactsByEmailOrPhone(email, phone *string) ([]home.DSContact, error) {
	var rows *sql.Rows
	var err error

	switch {
	case email != nil && phone != nil:
		rows, err = DB.Query(
			`SELECT `+selectCols+` FROM contacts WHERE deleted_at IS NULL AND (email = $1 OR phone_number = $2) ORDER BY created_at ASC`,
			*email, *phone,
		)
	case email != nil:
		rows, err = DB.Query(
			`SELECT `+selectCols+` FROM contacts WHERE deleted_at IS NULL AND email = $1 ORDER BY created_at ASC`,
			*email,
		)
	case phone != nil:
		rows, err = DB.Query(
			`SELECT `+selectCols+` FROM contacts WHERE deleted_at IS NULL AND phone_number = $1 ORDER BY created_at ASC`,
			*phone,
		)
	default:
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func FindContactByID(id int) (*home.DSContact, error) {
	row := DB.QueryRow(`SELECT `+selectCols+` FROM contacts WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanRow(row)
}

func FindAllLinkedContacts(primaryID int) ([]home.DSContact, error) {
	rows, err := DB.Query(
		`SELECT `+selectCols+` FROM contacts WHERE deleted_at IS NULL AND (id = $1 OR linked_id = $1) ORDER BY link_precedence ASC, created_at ASC`,
		primaryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func scanRows(rows *sql.Rows) ([]home.DSContact, error) {
	var out []home.DSContact
	for rows.Next() {
		c, err := scanFromRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}

func scanRow(row *sql.Row) (*home.DSContact, error) {
	var c home.DSContact
	var phone, email sql.NullString
	var linkedID sql.NullInt64
	var deletedAt sql.NullTime

	err := row.Scan(&c.ID, &phone, &email, &linkedID, &c.LinkPrecedence, &c.CreatedAt, &c.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}

	if phone.Valid {
		c.PhoneNumber = &phone.String
	}
	if email.Valid {
		c.Email = &email.String
	}
	if linkedID.Valid {
		v := int(linkedID.Int64)
		c.LinkedID = &v
	}
	if deletedAt.Valid {
		c.DeletedAt = &deletedAt.Time
	}
	return &c, nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanFromRow(s scannable) (*home.DSContact, error) {
	var c home.DSContact
	var phone, email sql.NullString
	var linkedID sql.NullInt64
	var deletedAt sql.NullTime

	err := s.Scan(&c.ID, &phone, &email, &linkedID, &c.LinkPrecedence, &c.CreatedAt, &c.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}

	if phone.Valid {
		c.PhoneNumber = &phone.String
	}
	if email.Valid {
		c.Email = &email.String
	}
	if linkedID.Valid {
		v := int(linkedID.Int64)
		c.LinkedID = &v
	}
	if deletedAt.Valid {
		c.DeletedAt = &deletedAt.Time
	}
	return &c, nil
}
