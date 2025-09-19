package etl

import (
	"fmt"
	"time"

	"github.com/admira-project/backend/internal/models"
	"github.com/sirupsen/logrus"
)

type Transformer struct {
	logger *logrus.Logger
}

func NewTransformer(logger *logrus.Logger) *Transformer {
	return &Transformer{logger: logger}
}

func (t *Transformer) Transform(adsData *models.AdsData, crmData *models.CrmData, since time.Time) ([]models.Metric, error) {
	t.logger.Info("Transforming data")

	cleanedAds := t.cleanAdsData(adsData, since)
	cleanedCrm := t.cleanCrmData(crmData, since)

	adsByUtm := t.groupAdsByUtm(cleanedAds)
	crmByUtm := t.groupCrmByUtm(cleanedCrm)

	metrics := t.joinAndCalculateMetrics(adsByUtm, crmByUtm)

	t.logger.Infof("Transformed data into %d metric records", len(metrics))
	return metrics, nil
}

func (t *Transformer) cleanAdsData(adsData *models.AdsData, since time.Time) []models.AdsPerformance {
	var cleaned []models.AdsPerformance

	for _, record := range adsData.External.Ads.Performance {
		// Validar fecha
		recordDate, err := time.Parse("2006-01-02", record.Date)
		if err != nil {
			t.logger.Warnf("Invalid date in ads record: %s", record.Date)
			continue
		}

		if !since.IsZero() && recordDate.Before(since) {
			continue
		}

		if record.CampaignID == "" || record.Channel == "" {
			t.logger.Warnf("Missing required fields in ads record: %+v", record)
			continue
		}

		if record.UtmCampaign == "" {
			record.UtmCampaign = "unknown"
		}
		if record.UtmSource == "" {
			record.UtmSource = "unknown"
		}
		if record.UtmMedium == "" {
			record.UtmMedium = "unknown"
		}

		cleaned = append(cleaned, record)
	}

	return cleaned
}

func (t *Transformer) cleanCrmData(crmData *models.CrmData, since time.Time) []models.CrmOpportunity {
	var cleaned []models.CrmOpportunity

	for _, record := range crmData.External.Crm.Opportunities {
		if !since.IsZero() && record.CreatedAt.Before(since) {
			continue
		}

		if record.OpportunityID == "" || record.Stage == "" {
			t.logger.Warnf("Missing required fields in CRM record: %+v", record)
			continue
		}

		// Establecer valores por defecto para UTMs si faltan
		if record.UtmCampaign == "" {
			record.UtmCampaign = "unknown"
		}
		if record.UtmSource == "" {
			record.UtmSource = "unknown"
		}
		if record.UtmMedium == "" {
			record.UtmMedium = "unknown"
		}

		cleaned = append(cleaned, record)
	}

	return cleaned
}

func (t *Transformer) groupAdsByUtm(ads []models.AdsPerformance) map[string][]models.AdsPerformance {
	groups := make(map[string][]models.AdsPerformance)

	for _, record := range ads {
		key := fmt.Sprintf("%s|%s|%s", record.UtmCampaign, record.UtmSource, record.UtmMedium)
		groups[key] = append(groups[key], record)
	}

	return groups
}

func (t *Transformer) groupCrmByUtm(crm []models.CrmOpportunity) map[string][]models.CrmOpportunity {
	groups := make(map[string][]models.CrmOpportunity)

	for _, record := range crm {
		key := fmt.Sprintf("%s|%s|%s", record.UtmCampaign, record.UtmSource, record.UtmMedium)
		groups[key] = append(groups[key], record)
	}

	return groups
}

func (t *Transformer) joinAndCalculateMetrics(
	adsByUtm map[string][]models.AdsPerformance,
	crmByUtm map[string][]models.CrmOpportunity,
) []models.Metric {
	var metrics []models.Metric

	for utmKey, adsRecords := range adsByUtm {
		crmRecords, exists := crmByUtm[utmKey]
		if !exists {
			crmRecords = []models.CrmOpportunity{}
		}

		// Calcular métricas agregadas para este UTM
		metric := t.calculateMetricsForUtm(adsRecords, crmRecords)
		metrics = append(metrics, metric)
	}

	return metrics
}

func (t *Transformer) calculateMetricsForUtm(
	adsRecords []models.AdsPerformance,
	crmRecords []models.CrmOpportunity,
) models.Metric {
	firstAd := adsRecords[0]
	metric := models.Metric{
		Date:        firstAd.Date,
		Channel:     firstAd.Channel,
		CampaignID:  firstAd.CampaignID,
		UtmCampaign: firstAd.UtmCampaign,
		UtmSource:   firstAd.UtmSource,
		UtmMedium:   firstAd.UtmMedium,
	}

	// Sumar métricas de Ads
	for _, ad := range adsRecords {
		metric.Clicks += ad.Clicks
		metric.Impressions += ad.Impressions
		metric.Cost += ad.Cost
	}

	// Calcular métricas de CRM
	for _, crm := range crmRecords {
		metric.Leads++

		if crm.Stage != "lead" {
			metric.Opportunities++
		}

		if crm.Stage == "closed_won" {
			metric.ClosedWon++
			metric.Revenue += crm.Amount
		}
	}

	// Calcular métricas derivadas
	if metric.Clicks > 0 {
		metric.CPC = metric.Cost / float64(metric.Clicks)
	}

	if metric.Leads > 0 {
		metric.CPA = metric.Cost / float64(metric.Leads)
	}

	if metric.Leads > 0 {
		metric.CvrLeadToOpp = float64(metric.Opportunities) / float64(metric.Leads)
	}

	if metric.Opportunities > 0 {
		metric.CvrOppToWon = float64(metric.ClosedWon) / float64(metric.Opportunities)
	}

	if metric.Cost > 0 {
		metric.Roas = metric.Revenue / metric.Cost
	}

	return metric
}
