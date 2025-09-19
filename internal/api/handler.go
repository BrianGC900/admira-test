package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/admira-project/backend/internal/etl"
	"github.com/admira-project/backend/internal/models"
	"github.com/admira-project/backend/internal/storage"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	extractor   *etl.Extractor
	transformer *etl.Transformer
	storage     storage.Storage
	logger      *logrus.Logger
}

func NewHandler(extractor *etl.Extractor, transformer *etl.Transformer, storage storage.Storage, logger *logrus.Logger) *Handler {
	return &Handler{
		extractor:   extractor,
		transformer: transformer,
		storage:     storage,
		logger:      logger,
	}
}

func (h *Handler) IngestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sinceParam := r.URL.Query().Get("since")
	var since time.Time
	var err error

	if sinceParam != "" {
		since, err = time.Parse("2006-01-02", sinceParam)
		if err != nil {
			h.logger.Warnf("Invalid since parameter: %s", sinceParam)
			http.Error(w, "Invalid since parameter format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	adsData, err := h.extractor.ExtractAdsData(ctx)
	if err != nil {
		h.logger.Errorf("Failed to extract ads data: %v", err)
		http.Error(w, "Failed to extract ads data", http.StatusInternalServerError)
		return
	}

	crmData, err := h.extractor.ExtractCrmData(ctx)
	if err != nil {
		h.logger.Errorf("Failed to extract CRM data: %v", err)
		http.Error(w, "Failed to extract CRM data", http.StatusInternalServerError)
		return
	}

	metrics, err := h.transformer.Transform(adsData, crmData, since)
	if err != nil {
		h.logger.Errorf("Failed to transform data: %v", err)
		http.Error(w, "Failed to transform data", http.StatusInternalServerError)
		return
	}

	if err := h.storage.SaveMetrics(metrics); err != nil {
		h.logger.Errorf("Failed to save metrics: %v", err)
		http.Error(w, "Failed to save metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "Data ingestion completed successfully",
		"count":   len(metrics),
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) MetricsChannelHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	request := models.MetricsRequest{
		From:    params.Get("from"),
		To:      params.Get("to"),
		Channel: params.Get("channel"),
	}

	if limitStr := params.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			request.Limit = limit
		}
	}

	if offsetStr := params.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			request.Offset = offset
		}
	}

	metrics, err := h.storage.GetMetricsByChannel(request)
	if err != nil {
		h.logger.Errorf("Failed to get metrics by channel: %v", err)
		http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

func (h *Handler) MetricsFunnelHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	request := models.MetricsRequest{
		From:        params.Get("from"),
		To:          params.Get("to"),
		UtmCampaign: params.Get("utm_campaign"),
	}

	if limitStr := params.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			request.Limit = limit
		}
	}

	if offsetStr := params.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			request.Offset = offset
		}
	}

	metrics, err := h.storage.GetMetricsByFunnel(request)
	if err != nil {
		h.logger.Errorf("Failed to get metrics by funnel: %v", err)
		http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *Handler) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}
