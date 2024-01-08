package main

import (
	"github.com/DarkOmap/metricsService/internal/client"
	"github.com/DarkOmap/metricsService/internal/parameters"
)

func main() {
	listenAddr, reportInterval, pollInterval := parameters.ParseFlagsAgent()

	<-client.Run(listenAddr, reportInterval, pollInterval)
}
