package router

import (
	"github.com/amitesh/dbaaas-monitoring/handlers"
	"github.com/gorilla/mux"
)

func SetupRoutes(h *handlers.Handler) *mux.Router {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Health check
	api.HandleFunc("/health", h.HealthCheck).Methods("GET")

	// Custom queries
	api.HandleFunc("/metrics/query", h.QueryMetrics).Methods("GET")
	api.HandleFunc("/metrics/query_range", h.QueryRangeMetrics).Methods("GET")

	// System metrics
	api.HandleFunc("/metrics/cpu", h.GetCPUMetrics).Methods("GET")
	api.HandleFunc("/metrics/memory", h.GetMemoryMetrics).Methods("GET")
	api.HandleFunc("/metrics/disk", h.GetDiskMetrics).Methods("GET")

	// Kubernetes metrics
	api.HandleFunc("/metrics/pods", h.GetPodMetrics).Methods("GET")
	api.HandleFunc("/metrics/nodes", h.GetNodeMetrics).Methods("GET")
	api.HandleFunc("/metrics/container/cpu", h.GetContainerCPU).Methods("GET")
	api.HandleFunc("/metrics/container/memory", h.GetContainerMemory).Methods("GET")

	// Targets
	api.HandleFunc("/metrics/targets", h.GetAllTargets).Methods("GET")

	return r
}

