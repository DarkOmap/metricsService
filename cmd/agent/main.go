package main

import (
	"github.com/DarkOmap/metricsService/internal/agent"
	"github.com/DarkOmap/metricsService/internal/client"
	"github.com/DarkOmap/metricsService/internal/logger"
	"github.com/DarkOmap/metricsService/internal/parameters"
	"go.uber.org/zap"
)

func main() {
	p := parameters.ParseFlagsAgent()

	if err := logger.Initialize("INFO", "stderr"); err != nil {
		panic(err)
	}

	logger.Log.Info("Create client")
	c := client.NewClient(p.ListenAddr, p.Key)
	logger.Log.Info("Create agent")
	a := agent.NewAgent(c, p.ReportInterval, p.PollInterval, p.RateLimit)
	logger.Log.Info("Agent start")
	err := a.Run()
	if err != nil {
		logger.Log.Fatal("Run agent", zap.Error(err))
	}
}
