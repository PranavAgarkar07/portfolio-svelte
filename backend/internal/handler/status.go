package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"portfolio-backend/internal/model"
)

func (h *Handlers) HandleHealthz(c *fiber.Ctx) error {
	status := fiber.Map{
		"status": "ok",
		"time":   time.Now().Unix(),
	}

	dbOK := h.DB != nil
	if dbOK {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := h.DB.PingContext(ctx); err != nil {
			dbOK = false
		}
	}
	status["database"] = dbOK

	s3OK := h.S3Client != nil
	status["s3"] = s3OK

	status["uptime_seconds"] = int64(time.Since(h.Metrics.StartTime).Seconds())

	if !dbOK && h.DB != nil {
		return c.Status(503).JSON(status)
	}
	return c.JSON(status)
}

func (h *Handlers) HandleStatus(c *fiber.Ctx) error {
	h.CacheMu.RLock()
	valid := time.Now().Before(h.Cache.ExpiresAt)
	hasData := h.Cache.Response.Summary != ""
	expired := !valid && hasData
	resp := h.Cache.Response
	h.CacheMu.RUnlock()

	h.Metrics.Mu.Lock()
	h.Metrics.TotalCalls++
	if valid {
		h.Metrics.CacheHits++
	} else if expired {
		h.Metrics.CacheHits++
	} else {
		h.Metrics.CacheMisses++
	}
	h.Metrics.Mu.Unlock()

	if valid {
		slog.Info("Serving from cache", reqID(c), "source", "cache")
		resp.Source = "cache"
		return c.JSON(resp)
	}

	if expired {
		slog.Info("Serving stale cache, refreshing in background", reqID(c), "source", "stale-cache")
		resp.Source = "stale-cache"
		go func() {
			start := time.Now()
			summary, err := h.generateDevLog()
			if err != nil {
				slog.Error("Background refresh failed", reqID(c), "error", err)
				h.Metrics.Mu.Lock()
				h.Metrics.LLMErrors++
				h.Metrics.Mu.Unlock()
				return
			}
			h.CacheMu.Lock()
			h.Cache.Response = model.DevLogResponse{
				Summary:    summary,
				LastUpdate: time.Now().Format("2006-01-02 15:04:05"),
				Source:     "live",
			}
			h.Cache.ExpiresAt = time.Now().Add(model.CacheDuration)
			h.CacheMu.Unlock()
			slog.Info("Cache refreshed in background", "latency_ms", time.Since(start).Milliseconds())
		}()
		return c.JSON(resp)
	}

	slog.Info("Cache empty, fetching live data", reqID(c))
	startTime := time.Now()
	summary, err := h.generateDevLog()
	latency := time.Since(startTime)

	h.Metrics.Mu.Lock()
	h.Metrics.TotalLatency += latency
	if err != nil {
		h.Metrics.LLMErrors++
	}
	h.Metrics.Mu.Unlock()

	if err != nil {
		slog.Error("Live fetch failed", reqID(c), "error", err, "latency_ms", latency.Milliseconds())
		return c.Status(500).JSON(fiber.Map{"error": "internal server error", "summary": "System Update: Offline (Retrying...)"})
	}

	newResp := model.DevLogResponse{
		Summary:    summary,
		LastUpdate: time.Now().Format("2006-01-02 15:04:05"),
		Source:     "live",
	}

	h.CacheMu.Lock()
	h.Cache.Response = newResp
	h.Cache.ExpiresAt = time.Now().Add(model.CacheDuration)
	h.CacheMu.Unlock()

	return c.JSON(newResp)
}

func (h *Handlers) HandleMetrics(c *fiber.Ctx) error {
	h.Metrics.Mu.Lock()
	defer h.Metrics.Mu.Unlock()
	avgLatency := time.Duration(0)
	if h.Metrics.TotalCalls > 0 {
		avgLatency = h.Metrics.TotalLatency / time.Duration(h.Metrics.TotalCalls)
	}
	return c.JSON(fiber.Map{
		"cache_hits":     h.Metrics.CacheHits,
		"cache_misses":   h.Metrics.CacheMisses,
		"llm_errors":     h.Metrics.LLMErrors,
		"total_requests": h.Metrics.TotalCalls,
		"avg_latency_ms": avgLatency.Milliseconds(),
		"uptime_seconds": int64(time.Since(h.Metrics.StartTime).Seconds()),
	})
}
