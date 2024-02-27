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
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Agent struct {
	reportInterval, pollInterval, rateLimit uint
	client                                  *client.Client
	pollCount                               atomic.Int64
	ms                                      runtime.MemStats
	vm                                      *mem.VirtualMemoryStat
	CPUutilization                          float64
}

func (a *Agent) Run() error {
	jobs := make(chan func() error, a.rateLimit)
	defer close(jobs)

	for w := uint(1); w <= a.rateLimit; w++ {
		go worker(jobs)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		logger.Log.Info("Send report start")
		for {
			select {
			case <-time.After(time.Duration(a.reportInterval) * time.Second):
				a.sendReport(egCtx, jobs)
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

	eg.Go(func() error {
		logger.Log.Info("Read virtual memory start")
		for {
			select {
			case <-time.After(time.Duration(a.pollInterval) * time.Second):
				var err error
				a.vm, err = mem.VirtualMemory()

				if err != nil {
					return err
				}

				CPUutilization, err := cpu.Percent(0, false)

				if err != nil {
					return err
				}

				a.CPUutilization = CPUutilization[0]

				a.pollCount.Add(1)
			case <-ctx.Done():
				logger.Log.Info("Read virtual memory done")
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

func (a *Agent) sendReport(ctx context.Context, jobs chan<- func() error) {
	jobs <- func() error {
		msForServer := memstats.GetMemStatsForServer(&a.ms)
		msForServer["RandomValue"] = rand.Float64()

		return a.client.SendBatch(ctx, msForServer)
	}

	jobs <- func() error {
		return a.client.SendCounter(ctx, "PollCount", a.pollCount.Load())
	}

	jobs <- func() error {
		vmForServer := memstats.GetVirtualMemoryForServer(a.vm)

		return a.client.SendBatch(ctx, vmForServer)
	}

	jobs <- func() error {
		return a.client.SendBatch(ctx, map[string]float64{"CPUutilization1": a.CPUutilization})
	}
}

func worker(jobs <-chan func() error) {
	for j := range jobs {
		err := j()

		if err != nil {
			logger.Log.Warn(
				"Error on sending to server",
				zap.Error(err),
			)
		}
	}
}

func NewAgent(client *client.Client, reportInterval, pollInterval, rateLimit uint) (a *Agent) {
	a = &Agent{reportInterval: reportInterval, pollInterval: pollInterval, client: client,
		rateLimit: rateLimit,
	}
	runtime.ReadMemStats(&a.ms)
	a.vm, _ = mem.VirtualMemory()
	CPUutilization, _ := cpu.Percent(0, false)
	a.CPUutilization = CPUutilization[0]
	return
}
