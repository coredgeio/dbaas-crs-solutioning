package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/amitesh/dbaaas-monitoring/config"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type PrometheusService struct {
	client v1.API
}

func NewPrometheusService(cfg *config.Config) (*PrometheusService, error) {
	client, err := api.NewClient(api.Config{
		Address: cfg.PrometheusURL,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Prometheus client: %v", err)
	}

	promAPI := v1.NewAPI(client)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, _, err = promAPI.Query(ctx, "up", time.Now())
	if err != nil {
		log.Printf("Warning: Cannot connect to Prometheus at %s: %v", cfg.PrometheusURL, err)
	} else {
		log.Printf("Successfully connected to Prometheus at %s", cfg.PrometheusURL)
	}

	return &PrometheusService{
		client: promAPI,
	}, nil
}

func (p *PrometheusService) Query(ctx context.Context, query string) (model.Value, error) {
	result, warnings, err := p.client.Query(ctx, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}

	if len(warnings) > 0 {
		log.Printf("Query warnings: %v", warnings)
	}

	return result, nil
}

func (p *PrometheusService) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (model.Value, error) {
	result, warnings, err := p.client.QueryRange(ctx, query, v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	})
	if err != nil {
		return nil, fmt.Errorf("query range error: %v", err)
	}

	if len(warnings) > 0 {
		log.Printf("Query range warnings: %v", warnings)
	}

	return result, nil
}

// Predefined metric queries
func (p *PrometheusService) GetCPUMetrics(ctx context.Context) (model.Value, error) {
	query := `100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`
	return p.Query(ctx, query)
}

func (p *PrometheusService) GetMemoryMetrics(ctx context.Context) (model.Value, error) {
	query := `(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100`
	return p.Query(ctx, query)
}

func (p *PrometheusService) GetDiskMetrics(ctx context.Context) (model.Value, error) {
	query := `(1 - (node_filesystem_avail_bytes{fstype!~"tmpfs|fuse.lxcfs|squashfs|vfat"} / node_filesystem_size_bytes{fstype!~"tmpfs|fuse.lxcfs|squashfs|vfat"})) * 100`
	return p.Query(ctx, query)
}

func (p *PrometheusService) GetPodMetrics(ctx context.Context) (model.Value, error) {
	query := `kube_pod_info`
	return p.Query(ctx, query)
}

func (p *PrometheusService) GetNodeMetrics(ctx context.Context) (model.Value, error) {
	query := `kube_node_info`
	return p.Query(ctx, query)
}

func (p *PrometheusService) GetContainerCPU(ctx context.Context) (model.Value, error) {
	query := `sum(rate(container_cpu_usage_seconds_total{container!=""}[5m])) by (pod, namespace)`
	return p.Query(ctx, query)
}

func (p *PrometheusService) GetContainerMemory(ctx context.Context) (model.Value, error) {
	query := `sum(container_memory_working_set_bytes{container!=""}) by (pod, namespace)`
	return p.Query(ctx, query)
}

func (p *PrometheusService) GetAllTargets(ctx context.Context) (model.Value, error) {
	query := `up`
	return p.Query(ctx, query)
}

