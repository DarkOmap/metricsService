package client

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	memstats "github.com/DarkOmap/metricsService/internal/memstats"
	"github.com/go-resty/resty/v2"
)

var serviceAddr string

func Run(listenAddr string, reportInterval, pollInterval uint) <-chan struct{} {
	done := make(chan struct{})
	var ms runtime.MemStats

	serviceAddr = listenAddr + "/update"

	pollCount := 0

	go func() {
		for {
			time.Sleep(time.Duration(reportInterval) * time.Second)
			msForServer := memstats.GetMemStatsForServer(&ms)
			err := pushStats(msForServer)

			if err != nil {
				log.Print(err.Error())
			}

			pollCountString := strconv.Itoa(pollCount)
			err = sendCounter("PollCount", pollCountString)

			if err != nil {
				log.Print(err.Error())
			}

			err = sendGauge("RandomValue", strconv.FormatFloat(rand.Float64(), 'f', -1, 64))

			if err != nil {
				log.Print(err.Error())
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(pollInterval) * time.Second)
			runtime.ReadMemStats(&ms)
			pollCount++
		}
	}()

	return done
}

func sendGauge(name, value string) error {
	param := map[string]string{"name": name, "value": value}

	client := resty.New()
	resp, err := client.R().SetPathParams(param).
		SetHeader("Content-Type", "text/plain").
		Post("http://" + serviceAddr + "/gauge/{name}/{value}")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("status not 200, current status %d", resp.StatusCode())
	}

	return nil
}

func sendCounter(name, value string) error {
	param := map[string]string{"name": name, "value": value}

	client := resty.New()

	resp, err := client.R().SetPathParams(param).
		SetHeader("Content-Type", "text/plain").
		Post("http://" + serviceAddr + "/counter/{name}/{value}")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("status not 200, current status %d", resp.StatusCode())
	}

	return nil
}

func pushStats(ms []memstats.StringMS) error {
	for _, val := range ms {
		err := sendGauge(val.Name, val.Value)

		if err != nil {
			return err
		}
	}

	return nil
}
