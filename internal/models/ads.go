package models

type AdsData struct {
	External struct {
		Ads struct {
			Performance []AdsPerformance `json:"performance"`
		} `json:"ads"`
	} `json:"external"`
}

type AdsPerformance struct {
	Date        string  `json:"date"`
	CampaignID  string  `json:"campaign_id"`
	Channel     string  `json:"channel"`
	Clicks      int     `json:"clicks"`
	Impressions int     `json:"impressions"`
	Cost        float64 `json:"cost"`
	UtmCampaign string  `json:"utm_campaign"`
	UtmSource   string  `json:"utm_source"`
	UtmMedium   string  `json:"utm_medium"`
}
