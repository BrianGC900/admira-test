package etl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/admira-project/backend/internal/models"
	"github.com/admira-project/backend/internal/utils"
	"github.com/sirupsen/logrus"
)

type Extractor struct {
	httpClient utils.HTTPClient
	adsAPIURL  string
	crmAPIURL  string
	logger     *logrus.Logger
}

func NewExtractor(httpClient utils.HTTPClient, adsAPIURL, crmAPIURL string, logger *logrus.Logger) *Extractor {
	if adsAPIURL == "" {
		adsAPIURL = "http://mock-ads:3001"
		logger.Warn("ADS_API_URL is empty, using default: http://mock-ads:3001")
	}
	if crmAPIURL == "" {
		crmAPIURL = "http://mock-crm:3002"
		logger.Warn("CRM_API_URL is empty, using default: http://mock-crm:3002")
	}

	logger.Infof("ADS_API_URL: %s", adsAPIURL)
	logger.Infof("CRM_API_URL: %s", crmAPIURL)

	return &Extractor{
		httpClient: httpClient,
		adsAPIURL:  adsAPIURL,
		crmAPIURL:  crmAPIURL,
		logger:     logger,
	}
}

func (e *Extractor) ExtractAdsData(ctx context.Context) (*models.AdsData, error) {
	e.logger.Info("Extracting Ads data")
	e.logger.Infof("Fetching from: %s", e.adsAPIURL)

	body, err := utils.FetchData(ctx, e.httpClient, e.adsAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ads data: %v", err)
	}

	var adsData models.AdsData
	if err := json.Unmarshal(body, &adsData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ads data: %v", err)
	}

	e.logger.Infof("Extracted %d ads performance records", len(adsData.External.Ads.Performance))
	return &adsData, nil
}

func (e *Extractor) ExtractCrmData(ctx context.Context) (*models.CrmData, error) {
	e.logger.Info("Extracting CRM data")
	e.logger.Infof("Fetching from: %s", e.crmAPIURL)

	body, err := utils.FetchData(ctx, e.httpClient, e.crmAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch crm data: %v", err)
	}

	var crmData models.CrmData
	if err := json.Unmarshal(body, &crmData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal crm data: %v", err)
	}

	e.logger.Infof("Extracted %d CRM opportunities", len(crmData.External.Crm.Opportunities))
	return &crmData, nil
}
