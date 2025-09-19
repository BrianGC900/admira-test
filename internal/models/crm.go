package models

import "time"

type CrmData struct {
	External struct {
		Crm struct {
			Opportunities []CrmOpportunity `json:"opportunities"`
		} `json:"crm"`
	} `json:"external"`
}

type CrmOpportunity struct {
	OpportunityID string    `json:"opportunity_id"`
	ContactEmail  string    `json:"contact_email"`
	Stage         string    `json:"stage"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	UtmCampaign   string    `json:"utm_campaign"`
	UtmSource     string    `json:"utm_source"`
	UtmMedium     string    `json:"utm_medium"`
}
