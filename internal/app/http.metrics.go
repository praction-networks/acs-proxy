package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/praction-networks/acs-proxy/internal/config"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	metricsServer   *http.Server
	metricsStopChan chan struct{}
	metricsDoneChan chan struct{}
)

func StartMetricsServer(mongoCLient *mongo.Client) *http.Server {

	port := config.AppConfig.ServerEnv.MetricsPort

	if port == "" {
		port = "9001"
	}

	metricsStopChan = make(chan struct{})
	metricsDoneChan = make(chan struct{})

	go collectSystemMetrics()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		metrics.Registry(),
		promhttp.HandlerFor(metrics.Registry(), promhttp.HandlerOpts{EnableOpenMetrics: true}),
	))
	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	// mux.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	_, _ = w.Write([]byte(`{"status":"ready"}`))
	// })

	mux.HandleFunc("/readiness", readinessHandler(mongoCLient))

	metricsServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	return metricsServer
}

func StopMetricsServer() error {
	if metricsStopChan != nil {
		close(metricsStopChan)
	}

	if metricsServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := metricsServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("metrics server shutdown failed: %w", err)
		}

		<-metricsDoneChan
	}

	return nil
}

func collectSystemMetrics() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-metricsStopChan:
			return
		case <-ticker.C:
			updateCPUMetrics()
			updateMemoryMetrics()
			updateDiskIOMetrics()
		}
	}
}

func updateCPUMetrics() {
	if percent, err := cpu.Percent(0, false); err == nil && len(percent) > 0 {
		metrics.CPUUsage.Set(percent[0])
	}
}

func updateMemoryMetrics() {
	if memStat, err := mem.VirtualMemory(); err == nil {
		metrics.MemoryUsage.Set(float64(memStat.Used))
	}
}

func updateDiskIOMetrics() {
	if usage, err := disk.Usage("/"); err == nil {
		metrics.DiskUsagePercent.Set(usage.UsedPercent)
	}
	if ioStats, err := disk.IOCounters(); err == nil {
		for _, stat := range ioStats {
			metrics.DiskReadBytes.Set(float64(stat.ReadBytes))
			metrics.DiskWriteBytes.Set(float64(stat.WriteBytes))
			break
		}
	} else {
		logger.Warn("Failed to fetch disk I/O stats", err)
	}
}

func readinessHandler(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		mongoHealthy := mongoClient != nil && mongoClient.Ping(ctx, nil) == nil

		response := ReadinessResponse{
			Status: "ready",
		}

		statusCode := http.StatusOK

		if !mongoHealthy {
			response.Status = "unhealthy"
			statusCode = http.StatusServiceUnavailable
			if !mongoHealthy {
				response.Mongo = "down"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(response)
	}
}

type ReadinessResponse struct {
	Status string `json:"status"`
	Mongo  string `json:"mongo,omitempty"`
	NATS   string `json:"nats,omitempty"`
}
