package parameters

import (
	"flag"
	"os"
	"strconv"
)

func ParseFlagsAgent() (listenAddr string, reportInterval, pollInterval uint) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.StringVar(&listenAddr, "a", "localhost:8080", "address and port to server")
	f.UintVar(&reportInterval, "r", 10, "report interval")
	f.UintVar(&pollInterval, "p", 2, "poll interval")
	f.Parse(os.Args[1:])

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		listenAddr = envAddr
	}

	if envRI := os.Getenv("REPORT_INTERVAL"); envRI != "" {
		intRI, err := strconv.ParseUint(envRI, 10, 32)

		if err == nil {
			reportInterval = uint(intRI)
		}
	}

	if envPI := os.Getenv("POLL_INTERVAL"); envPI != "" {
		intPI, err := strconv.ParseUint(envPI, 10, 32)

		if err == nil {
			pollInterval = uint(intPI)
		}
	}

	return
}

func ParseFlagsServer() (flagRunAddr string) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	f.Parse(os.Args[1:])

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		flagRunAddr = envAddr
	}

	return
}
