package home

import "time"

type DSContact struct {
	ID             int
	PhoneNumber    *string
	Email          *string
	LinkedID       *int
	LinkPrecedence string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type DSIdentifyRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
}

type DSIdentifyResponse struct {
	Contact DSContactPayload `json:"contact"`
}

type DSContactPayload struct {
	PrimaryContactID    int      `json:"primaryContactId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNumbers"`
	SecondaryContactIDs []int    `json:"secondaryContactIds"`
}
