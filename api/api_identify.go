package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranesh/bitespeed/db"
	"github.com/pranesh/bitespeed/home"
)

func HandleIdentify(c *gin.Context) {
	var req home.DSIdentifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	req.Email = home.TrimPtr(req.Email)
	req.PhoneNumber = home.TrimPtr(req.PhoneNumber)

	if req.Email == nil && req.PhoneNumber == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": home.ErrInvalidInput.Error()})
		return
	}

	resp, err := processIdentify(req.Email, req.PhoneNumber)
	if err != nil {
		log.Println("api identify error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func processIdentify(email, phone *string) (*home.DSIdentifyResponse, error) {
	contacts, err := db.FindContactsByEmailOrPhone(email, phone)
	if err != nil {
		return nil, err
	}

	// case 1 no existing contacts, new primary
	if len(contacts) == 0 {
		primary, err := db.InsertPrimaryContact(email, phone)
		if err != nil {
			return nil, err
		}
		return buildResponse(primary, []home.DSContact{*primary}), nil
	}

	primaryIDs := collectPrimaryIDs(contacts)

	// case 2 single primary
	if len(primaryIDs) == 1 {
		primaryID := primaryIDs[0]

		if shouldCreateSecondary(contacts, email, phone) {
			if _, err := db.InsertSecondaryContact(email, phone, primaryID); err != nil {
				return nil, err
			}
		}

		primary, err := db.FindContactByID(primaryID)
		if err != nil {
			return nil, err
		}
		all, err := db.FindAllLinkedContacts(primaryID)
		if err != nil {
			return nil, err
		}
		return buildResponse(primary, all), nil
	}

	// case 3 multiple primary groups, merge
	oldest, err := mergeGroups(primaryIDs)
	if err != nil {
		return nil, err
	}

	all, err := db.FindAllLinkedContacts(oldest.ID)
	if err != nil {
		return nil, err
	}
	return buildResponse(oldest, all), nil
}

func collectPrimaryIDs(contacts []home.DSContact) []int {
	seen := map[int]bool{}
	var ids []int
	for _, c := range contacts {
		pid := c.ID
		if c.LinkPrecedence == home.LinkSecondary && c.LinkedID != nil {
			pid = *c.LinkedID
		}
		if !seen[pid] {
			seen[pid] = true
			ids = append(ids, pid)
		}
	}
	return ids
}

func shouldCreateSecondary(contacts []home.DSContact, email, phone *string) bool {
	if email == nil || phone == nil {
		return false
	}
	for _, c := range contacts {
		eMatch := c.Email != nil && *c.Email == *email
		pMatch := c.PhoneNumber != nil && *c.PhoneNumber == *phone
		if eMatch && pMatch {
			return false
		}
	}
	return true
}

// keep the oldest primary n link others as secondary
func mergeGroups(primaryIDs []int) (*home.DSContact, error) {
	var primaries []*home.DSContact
	for _, id := range primaryIDs {
		p, err := db.FindContactByID(id)
		if err != nil {
			return nil, err
		}
		primaries = append(primaries, p)
	}

	oldest := primaries[0]
	for _, p := range primaries[1:] {
		if p.CreatedAt.Before(oldest.CreatedAt) {
			oldest = p
		}
	}

	for _, p := range primaries {
		if p.ID == oldest.ID {
			continue
		}
		if err := db.UpdateLinkedContactsToNewPrimary(p.ID, oldest.ID); err != nil {
			return nil, err
		}
		if err := db.UpdateContactToSecondary(p.ID, oldest.ID); err != nil {
			return nil, err
		}
	}

	return oldest, nil
}

func buildResponse(primary *home.DSContact, all []home.DSContact) *home.DSIdentifyResponse {
	var emails, phones []string
	var secondaryIDs []int

	if primary.Email != nil {
		emails = append(emails, *primary.Email)
	}
	if primary.PhoneNumber != nil {
		phones = append(phones, *primary.PhoneNumber)
	}

	for _, c := range all {
		if c.ID == primary.ID {
			continue
		}
		if c.Email != nil {
			emails = append(emails, *c.Email)
		}
		if c.PhoneNumber != nil {
			phones = append(phones, *c.PhoneNumber)
		}
		secondaryIDs = append(secondaryIDs, c.ID)
	}

	return &home.DSIdentifyResponse{
		Contact: home.DSContactPayload{
			PrimaryContactID:    primary.ID,
			Emails:              home.UniqueStrings(emails),
			PhoneNumbers:        home.UniqueStrings(phones),
			SecondaryContactIDs: home.UniqueInts(secondaryIDs),
		},
	}
}
