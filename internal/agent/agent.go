package agent

import (
	"context"
	"fmt"
	"math/rand"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/DarkOmap/metricsService/internal/client"
	"github.com/DarkOmap/metricsService/internal/logger"
	"github.com/DarkOmap/metricsService/internal/memstats"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Agent struct {
	reportInterval, pollInterval uint
	client                       *client.Client
	pollCount                    atomic.Int64
	ms                           runtime.MemStats
}

func (a *Agent) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		logger.Log.Info("Send report start")
		for {
			select {
			case <-time.After(time.Duration(a.reportInterval) * time.Second):
				a.sendReport(egCtx)
			case <-egCtx.Done():
				logger.Log.Info("Send report done")
				return nil
			}
		}
	})

	eg.Go(func() error {
		logger.Log.Info("Read mem stats start")
		for {
			select {
			case <-time.After(time.Duration(a.pollInterval) * time.Second):
				runtime.ReadMemStats(&a.ms)
				a.pollCount.Add(1)
			case <-ctx.Done():
				logger.Log.Info("Read mem stats done")
				return nil
			}
		}
	})

	if err := eg.Wait(); err != nil {
		logger.Log.Error("Problem with working agent", zap.Error(err))
		return fmt.Errorf("problem with working agent: %w", err)
	}

	return nil
}

func (a *Agent) sendReport(ctx context.Context) {
	msForServer := memstats.GetMemStatsForServer(&a.ms)
	msForServer["RandomValue"] = rand.Float64()

	err := a.client.SendBatch(ctx, msForServer)

	if err != nil {
		logger.Log.Warn(
			"Push memstats",
			zap.Error(err),
		)
	}

	err = a.client.SendCounter(ctx, "PollCount", a.pollCount.Load())

	if err != nil {
		logger.Log.Warn(
			"Error on sending poll count",
			zap.Int64("value", a.pollCount.Load()),
			zap.Error(err),
		)
	}
}

func NewAgent(client *client.Client, reportInterval, pollInterval uint) (a *Agent) {
	a = &Agent{reportInterval: reportInterval, pollInterval: pollInterval, client: client}
	runtime.ReadMemStats(&a.ms)
	return
}
