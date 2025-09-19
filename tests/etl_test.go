package tests

import (
	"testing"
	"time"

	"github.com/admira-project/backend/internal/etl"
	"github.com/admira-project/backend/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTransformer(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Silenciar logs en tests

	transformer := etl.NewTransformer(logger)

	// Crear datos de prueba
	adsData := &models.AdsData{}
	ad := models.AdsPerformance{
		Date:        "2023-01-01",
		CampaignID:  "TEST-001",
		Channel:     "google_ads",
		Clicks:      100,
		Impressions: 10000,
		Cost:        50.0,
		UtmCampaign: "test_campaign",
		UtmSource:   "google",
		UtmMedium:   "cpc",
	}
	adsData.External.Ads.Performance = append(adsData.External.Ads.Performance, ad)

	crmData := &models.CrmData{}
	opp := models.CrmOpportunity{
		OpportunityID: "OPP-001",
		ContactEmail:  "test@example.com",
		Stage:         "closed_won",
		Amount:        500.0,
		CreatedAt:     time.Now(),
		UtmCampaign:   "test_campaign",
		UtmSource:     "google",
		UtmMedium:     "cpc",
	}
	crmData.External.Crm.Opportunities = append(crmData.External.Crm.Opportunities, opp)

	// Probar transformación
	metrics, err := transformer.Transform(adsData, crmData, time.Time{})
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)

	metric := metrics[0]
	assert.Equal(t, "google_ads", metric.Channel)
	assert.Equal(t, "TEST-001", metric.CampaignID)
	assert.Equal(t, 100, metric.Clicks)
	assert.Equal(t, 10000, metric.Impressions)
	assert.Equal(t, 50.0, metric.Cost)
	assert.Equal(t, 1, metric.Leads)
	assert.Equal(t, 1, metric.Opportunities)
	assert.Equal(t, 1, metric.ClosedWon)
	assert.Equal(t, 500.0, metric.Revenue)
	assert.Equal(t, 0.5, metric.CPC)          // 50 / 100
	assert.Equal(t, 50.0, metric.CPA)         // 50 / 1
	assert.Equal(t, 1.0, metric.CvrLeadToOpp) // 1 / 1
	assert.Equal(t, 1.0, metric.CvrOppToWon)  // 1 / 1
	assert.Equal(t, 10.0, metric.Roas)        // 500 / 50
}

func TestTransformerWithMissingUtm(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	transformer := etl.NewTransformer(logger)

	// Crear datos con UTM faltante
	adsData := &models.AdsData{}
	ad := models.AdsPerformance{
		Date:        "2023-01-01",
		CampaignID:  "TEST-001",
		Channel:     "google_ads",
		Clicks:      100,
		Impressions: 10000,
		Cost:        50.0,
		// UTM fields intentionally missing
	}
	adsData.External.Ads.Performance = append(adsData.External.Ads.Performance, ad)

	crmData := &models.CrmData{}

	// Probar transformación
	metrics, err := transformer.Transform(adsData, crmData, time.Time{})
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)

	metric := metrics[0]
	assert.Equal(t, "unknown", metric.UtmCampaign)
	assert.Equal(t, "unknown", metric.UtmSource)
	assert.Equal(t, "unknown", metric.UtmMedium)
}

func TestTransformerWithDateFilter(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	transformer := etl.NewTransformer(logger)

	adsData := &models.AdsData{}

	oldAd := models.AdsPerformance{
		Date:        "2023-01-01",
		CampaignID:  "OLD-001",
		Channel:     "google_ads",
		Clicks:      50,
		Impressions: 5000,
		Cost:        25.0,
		UtmCampaign: "old_campaign",
	}

	recentAd := models.AdsPerformance{
		Date:        "2023-02-01",
		CampaignID:  "RECENT-001",
		Channel:     "google_ads",
		Clicks:      50,
		Impressions: 5000,
		Cost:        25.0,
		UtmCampaign: "recent_campaign",
	}

	adsData.External.Ads.Performance = append(adsData.External.Ads.Performance, oldAd, recentAd)

	crmData := &models.CrmData{}

	// Filtrar desde 2023-01-15
	since, _ := time.Parse("2006-01-02", "2023-01-15")
	metrics, err := transformer.Transform(adsData, crmData, since)

	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "recent_campaign", metrics[0].UtmCampaign)
}
