package main

import (
	"log"
	"net/http"

	"github.com/amitesh/dbaaas-monitoring/config"
	"github.com/amitesh/dbaaas-monitoring/handlers"
	"github.com/amitesh/dbaaas-monitoring/router"
	"github.com/amitesh/dbaaas-monitoring/services"
	"github.com/rs/cors"
)

func main() {
	log.Println("========================================")
	log.Println("🚀 DBaaS Monitoring API")
	log.Println("========================================")
	log.Println()

	// Load configuration with auto-discovery
	cfg := config.Load()
	log.Println()

	// Initialize Prometheus service
	log.Println("📡 Initializing Prometheus service...")
	promService, err := services.NewPrometheusService(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize Prometheus service: %v", err)
	}
	log.Println("✅ Prometheus service initialized")
	log.Println()

	// Initialize handlers
	h := handlers.NewHandler(promService)

	// Setup routes
	r := router.SetupRoutes(h)

	// CORS setup
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Server info
	log.Println("========================================")
	log.Println("🌐 API Server Started Successfully!")
	log.Println("========================================")
	log.Printf("📊 Prometheus: %s", cfg.PrometheusURL)
	log.Printf("🔌 API Port: %s", cfg.Port)
	log.Printf("🏥 Health Check: http://localhost:%s/api/health", cfg.Port)
	log.Println()
	log.Println("📋 Available Endpoints:")
	log.Println("  - GET /api/health")
	log.Println("  - GET /api/metrics/query?query=<prom-query>")
	log.Println("  - GET /api/metrics/cpu")
	log.Println("  - GET /api/metrics/memory")
	log.Println("  - GET /api/metrics/disk")
	log.Println("  - GET /api/metrics/pods")
	log.Println("  - GET /api/metrics/nodes")
	log.Println("  - GET /api/metrics/container/cpu")
	log.Println("  - GET /api/metrics/container/memory")
	log.Println("========================================")
	log.Println()
	
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

