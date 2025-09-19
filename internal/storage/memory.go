package storage

import (
	"time"

	"github.com/admira-project/backend/internal/models"
)

type Storage interface {
	SaveMetrics(metrics []models.Metric) error
	GetMetricsByChannel(request models.MetricsRequest) ([]models.Metric, error)
	GetMetricsByFunnel(request models.MetricsRequest) ([]models.Metric, error)
}

type MemoryStorage struct {
	metrics []models.Metric
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		metrics: make([]models.Metric, 0),
	}
}

func (s *MemoryStorage) SaveMetrics(metrics []models.Metric) error {
	s.metrics = append(s.metrics, metrics...)
	return nil
}

func (s *MemoryStorage) GetMetricsByChannel(request models.MetricsRequest) ([]models.Metric, error) {
	var filtered []models.Metric

	for _, metric := range s.metrics {
		if !s.filterByDate(metric, request.From, request.To) {
			continue
		}

		if request.Channel != "" && metric.Channel != request.Channel {
			continue
		}

		filtered = append(filtered, metric)
	}

	start, end := s.applyPagination(filtered, request.Limit, request.Offset)
	return filtered[start:end], nil
}

func (s *MemoryStorage) GetMetricsByFunnel(request models.MetricsRequest) ([]models.Metric, error) {
	var filtered []models.Metric

	for _, metric := range s.metrics {
		if !s.filterByDate(metric, request.From, request.To) {
			continue
		}

		if request.UtmCampaign != "" && metric.UtmCampaign != request.UtmCampaign {
			continue
		}

		filtered = append(filtered, metric)
	}

	start, end := s.applyPagination(filtered, request.Limit, request.Offset)
	return filtered[start:end], nil
}

func (s *MemoryStorage) filterByDate(metric models.Metric, from, to string) bool {
	metricDate, err := time.Parse("2006-01-02", metric.Date)
	if err != nil {
		return false
	}

	if from != "" {
		fromDate, err := time.Parse("2006-01-02", from)
		if err == nil && metricDate.Before(fromDate) {
			return false
		}
	}

	if to != "" {
		toDate, err := time.Parse("2006-01-02", to)
		if err == nil && metricDate.After(toDate) {
			return false
		}
	}

	return true
}

func (s *MemoryStorage) applyPagination(metrics []models.Metric, limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 50 // Valor por defecto
	}

	if offset < 0 {
		offset = 0
	}

	start := offset
	if start > len(metrics) {
		start = len(metrics)
	}

	end := start + limit
	if end > len(metrics) {
		end = len(metrics)
	}

	return start, end
}
