package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/amitesh/dbaaas-monitoring/models"
	"github.com/amitesh/dbaaas-monitoring/services"
)

type Handler struct {
	promService *services.PrometheusService
}

func NewHandler(promService *services.PrometheusService) *Handler {
	return &Handler{
		promService: promService,
	}
}

// HealthCheck handler
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:    "success",
		Message:   "DBaaS Monitoring API is running",
		Timestamp: time.Now().Unix(),
	}
	sendJSON(w, http.StatusOK, response)
}

// QueryMetrics - Custom Prometheus query
func (h *Handler) QueryMetrics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		sendError(w, http.StatusBadRequest, "query parameter is required")
		return
	}

	ctx := r.Context()
	result, err := h.promService.Query(ctx, query)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	sendJSON(w, http.StatusOK, response)
}

// QueryRangeMetrics - Time series query
func (h *Handler) QueryRangeMetrics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		sendError(w, http.StatusBadRequest, "query parameter is required")
		return
	}

	// Default time range: last 1 hour
	end := time.Now()
	start := end.Add(-1 * time.Hour)
	step := 30 * time.Second

	// Parse optional parameters
	if startParam := r.URL.Query().Get("start"); startParam != "" {
		if t, err := time.Parse(time.RFC3339, startParam); err == nil {
			start = t
		}
	}
	if endParam := r.URL.Query().Get("end"); endParam != "" {
		if t, err := time.Parse(time.RFC3339, endParam); err == nil {
			end = t
		}
	}

	ctx := r.Context()
	result, err := h.promService.QueryRange(ctx, query, start, end, step)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	sendJSON(w, http.StatusOK, response)
}

// GetCPUMetrics handler
func (h *Handler) GetCPUMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.promService.GetCPUMetrics(ctx)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	sendJSON(w, http.StatusOK, response)
}

// GetMemoryMetrics handler
func (h *Handler) GetMemoryMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.promService.GetMemoryMetrics(ctx)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	sendJSON(w, http.StatusOK, response)
}

// GetDiskMetrics handler
func (h *Handler) GetDiskMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.promService.GetDiskMetrics(ctx)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	sendJSON(w, http.StatusOK, response)
}

// GetPodMetrics handler
func (h *Handler) GetPodMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.promService.GetPodMetrics(ctx)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	sendJSON(w, http.StatusOK, response)
}

// GetNodeMetrics handler
func (h *Handler) GetNodeMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.promService.GetNodeMetrics(ctx)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	sendJSON(w, http.StatusOK, response)
}

// GetContainerCPU handler
func (h *Handler) GetContainerCPU(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.promService.GetContainerCPU(ctx)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	sendJSON(w, http.StatusOK, response)
}

// GetContainerMemory handler
func (h *Handler) GetContainerMemory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.promService.GetContainerMemory(ctx)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	sendJSON(w, http.StatusOK, response)
}

// GetAllTargets handler
func (h *Handler) GetAllTargets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	result, err := h.promService.GetAllTargets(ctx)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.MetricsResponse{
		Status: "success",
		Data:   result,
	}
	sendJSON(w, http.StatusOK, response)
}

// Helper functions
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func sendError(w http.ResponseWriter, status int, message string) {
	response := models.ErrorResponse{
		Status: "error",
		Error:  message,
	}
	sendJSON(w, status, response)
}

