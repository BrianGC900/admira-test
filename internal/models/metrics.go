package models

type Metric struct {
	Date          string  `json:"date"`
	Channel       string  `json:"channel"`
	CampaignID    string  `json:"campaign_id"`
	UtmCampaign   string  `json:"utm_campaign"`
	UtmSource     string  `json:"utm_source"`
	UtmMedium     string  `json:"utm_medium"`
	Clicks        int     `json:"clicks"`
	Impressions   int     `json:"impressions"`
	Cost          float64 `json:"cost"`
	Leads         int     `json:"leads"`
	Opportunities int     `json:"opportunities"`
	ClosedWon     int     `json:"closed_won"`
	Revenue       float64 `json:"revenue"`
	CPC           float64 `json:"cpc"`
	CPA           float64 `json:"cpa"`
	CvrLeadToOpp  float64 `json:"cvr_lead_to_opp"`
	CvrOppToWon   float64 `json:"cvr_opp_to_won"`
	Roas          float64 `json:"roas"`
}

type MetricsRequest struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Channel     string `json:"channel"`
	UtmCampaign string `json:"utm_campaign"`
	Limit       int    `json:"limit"`
	Offset      int    `json:"offset"`
}
