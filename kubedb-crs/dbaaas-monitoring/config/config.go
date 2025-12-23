package config

import (
	"log"
	"os"

	"github.com/amitesh/dbaaas-monitoring/utils"
)

type Config struct {
	PrometheusURL string
	Port          string
	Environment   string
	NodeIP        string
	NodePort      string
}

func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Dynamically discover Prometheus endpoint
	log.Println("🔍 Discovering Prometheus endpoint...")
	
	endpoint, err := utils.GetPrometheusEndpoint()
	if err != nil {
		log.Printf("⚠️  Warning: Could not auto-discover Prometheus endpoint: %v", err)
		log.Printf("Falling back to default/environment variable")
		cfg.PrometheusURL = getEnv("PROMETHEUS_URL", "http://localhost:9090")
	} else {
		cfg.PrometheusURL = endpoint.URL
		cfg.NodeIP = endpoint.NodeIP
		cfg.NodePort = endpoint.NodePort
	}

	// Validate Prometheus connection
	log.Println("🔗 Validating Prometheus connection...")
	if err := utils.CheckPrometheusHealth(cfg.PrometheusURL); err != nil {
		log.Printf("⚠️  Warning: Prometheus health check failed: %v", err)
		log.Printf("API will still start, but metrics may not be available")
	}

	log.Printf("✅ Configuration loaded successfully:")
	log.Printf("  - Prometheus URL: %s", cfg.PrometheusURL)
	if cfg.NodeIP != "" {
		log.Printf("  - Node IP: %s", cfg.NodeIP)
		log.Printf("  - Node Port: %s", cfg.NodePort)
	}
	log.Printf("  - API Port: %s", cfg.Port)
	log.Printf("  - Environment: %s", cfg.Environment)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

