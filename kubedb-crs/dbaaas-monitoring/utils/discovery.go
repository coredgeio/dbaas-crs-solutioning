package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type PrometheusEndpoint struct {
	URL      string
	NodeIP   string
	NodePort string
}

type DiscoveryConfig struct {
	Namespace   string
	ServiceName string
}

// GetPrometheusEndpoint dynamically discovers Prometheus endpoint
func GetPrometheusEndpoint() (*PrometheusEndpoint, error) {
	// Check if running inside Kubernetes cluster
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		log.Println("Running inside Kubernetes cluster")
		return getInClusterEndpoint()
	}

	log.Println("Running outside Kubernetes cluster")
	
	// Try to get from environment variable first
	if url := os.Getenv("PROMETHEUS_URL"); url != "" {
		log.Printf("Using Prometheus URL from environment: %s", url)
		return &PrometheusEndpoint{URL: url}, nil
	}

	// Auto-discover from kubectl
	config := DiscoveryConfig{
		Namespace:   getEnvOrDefault("PROMETHEUS_NAMESPACE", "monitoring"),
		ServiceName: getEnvOrDefault("PROMETHEUS_SERVICE", "kube-prom-stack-kube-prome-prometheus"),
	}
	return discoverFromKubectl(config)
}

// getInClusterEndpoint returns in-cluster service endpoint
func getInClusterEndpoint() (*PrometheusEndpoint, error) {
	namespace := getEnvOrDefault("PROMETHEUS_NAMESPACE", "monitoring")
	service := getEnvOrDefault("PROMETHEUS_SERVICE", "kube-prom-stack-kube-prome-prometheus")
	
	url := fmt.Sprintf("http://%s.%s.svc.cluster.local:9090", service, namespace)
	log.Printf("Using in-cluster Prometheus endpoint: %s", url)
	
	return &PrometheusEndpoint{
		URL: url,
	}, nil
}

// discoverFromKubectl discovers Prometheus endpoint using kubectl
func discoverFromKubectl(config DiscoveryConfig) (*PrometheusEndpoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get NodePort
	nodePortCmd := exec.CommandContext(ctx, "kubectl", "get", "svc", 
		"-n", config.Namespace,
		config.ServiceName,
		"-o", "jsonpath={.spec.ports[0].nodePort}")
	
	nodePortOutput, err := nodePortCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get NodePort: %v", err)
	}
	nodePort := strings.TrimSpace(string(nodePortOutput))

	if nodePort == "" {
		return nil, fmt.Errorf("NodePort not found")
	}

	// Get Node IP
	nodeIPCmd := exec.CommandContext(ctx, "kubectl", "get", "nodes",
		"-o", "jsonpath={.items[0].status.addresses[?(@.type==\"InternalIP\")].address}")
	
	nodeIPOutput, err := nodeIPCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get Node IP: %v", err)
	}
	nodeIP := strings.TrimSpace(string(nodeIPOutput))

	if nodeIP == "" {
		return nil, fmt.Errorf("Node IP not found")
	}

	url := fmt.Sprintf("http://%s:%s", nodeIP, nodePort)
	
	log.Printf("✓ Discovered Prometheus endpoint:")
	log.Printf("  - Namespace: %s", config.Namespace)
	log.Printf("  - Service: %s", config.ServiceName)
	log.Printf("  - Node IP: %s", nodeIP)
	log.Printf("  - NodePort: %s", nodePort)
	log.Printf("  - Full URL: %s", url)

	return &PrometheusEndpoint{
		URL:      url,
		NodeIP:   nodeIP,
		NodePort: nodePort,
	}, nil
}

// CheckPrometheusHealth checks if Prometheus is accessible
func CheckPrometheusHealth(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to query Prometheus
	testURL := fmt.Sprintf("%s/api/v1/query?query=up", url)
	
	cmd := exec.CommandContext(ctx, "curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", testURL)
	output, err := cmd.Output()
	
	if err != nil {
		return fmt.Errorf("failed to connect to Prometheus: %v", err)
	}

	statusCode := strings.TrimSpace(string(output))
	if statusCode != "200" {
		return fmt.Errorf("Prometheus returned status code: %s", statusCode)
	}

	log.Printf("✓ Prometheus health check passed")
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
