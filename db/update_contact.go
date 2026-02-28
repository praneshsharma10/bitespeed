package db

import (
	"github.com/pranesh/bitespeed/home"
)

func UpdateContactToSecondary(contactID, newPrimaryID int) error {
	_, err := DB.Exec(
		`UPDATE contacts SET linked_id = $1, link_precedence = $2, updated_at = NOW() WHERE id = $3`,
		newPrimaryID, home.LinkSecondary, contactID,
	)
	return err
}

func UpdateLinkedContactsToNewPrimary(oldPrimaryID, newPrimaryID int) error {
	_, err := DB.Exec(
		`UPDATE contacts SET linked_id = $1, updated_at = NOW() WHERE linked_id = $2`,
		newPrimaryID, oldPrimaryID,
	)
	return err
}
