package monitoring

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/praction-networks/common/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type ConnectionMonitor struct {
	mongoClient *mongo.Client
	isHealthy   atomic.Bool
	metrics     struct {
		checks        atomic.Int64
		failures      atomic.Int64
		lastCheckTime atomic.Value // time.Time
	}
}

func New(mongo *mongo.Client) *ConnectionMonitor {
	cm := &ConnectionMonitor{
		mongoClient: mongo,
	}
	cm.isHealthy.Store(true)
	cm.metrics.lastCheckTime.Store(time.Now())
	return cm
}

func (m *ConnectionMonitor) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial check
	m.checkHealth(ctx)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping connection monitor")
			return
		case <-ticker.C:
			m.checkHealth(ctx)
		}
	}
}

func (m *ConnectionMonitor) IsHealthy() bool {
	return m.isHealthy.Load()
}

func (m *ConnectionMonitor) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"healthy":         m.isHealthy.Load(),
		"total_checks":    m.metrics.checks.Load(),
		"total_failures":  m.metrics.failures.Load(),
		"last_check_time": m.metrics.lastCheckTime.Load().(time.Time),
	}
}

func (m *ConnectionMonitor) checkHealth(ctx context.Context) {
	m.metrics.checks.Add(1)
	m.metrics.lastCheckTime.Store(time.Now())

	mongoHealthy := m.checkMongo(ctx)

	healthy := mongoHealthy
	if !healthy {
		m.metrics.failures.Add(1)
	}

	prev := m.isHealthy.Swap(healthy)
	if prev != healthy {
		logger.Warn("Connection health status changed",
			"newStatus", healthy,
			"mongo", mongoHealthy,
		)
	}
}

func (m *ConnectionMonitor) checkMongo(ctx context.Context) bool {
	if m.mongoClient == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := m.mongoClient.Ping(ctx, nil); err != nil {
		logger.Error("MongoDB health check failed", "error", err)
		return false
	}
	return true
}
