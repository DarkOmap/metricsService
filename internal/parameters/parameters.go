package parameters

import (
	"flag"
	"os"
	"strconv"
)

type AgentParameters struct {
	ListenAddr, Key                         string
	ReportInterval, PollInterval, RateLimit uint
}

func ParseFlagsAgent() (p AgentParameters) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.StringVar(&p.ListenAddr, "a", "localhost:8080", "address and port to server")
	f.StringVar(&p.Key, "k", "", "hash key")
	f.UintVar(&p.ReportInterval, "r", 10, "report interval")
	f.UintVar(&p.PollInterval, "p", 2, "poll interval")
	f.UintVar(&p.RateLimit, "l", 10, "rate limit")
	f.Parse(os.Args[1:])

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		p.ListenAddr = envAddr
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		p.Key = envKey
	}

	if envRI := os.Getenv("REPORT_INTERVAL"); envRI != "" {
		intRI, err := strconv.ParseUint(envRI, 10, 32)

		if err == nil {
			p.ReportInterval = uint(intRI)
		}
	}

	if envPI := os.Getenv("POLL_INTERVAL"); envPI != "" {
		intPI, err := strconv.ParseUint(envPI, 10, 32)

		if err == nil {
			p.PollInterval = uint(intPI)
		}
	}

	if envRL := os.Getenv("RATE_LIMIT"); envRL != "" {
		intRL, err := strconv.ParseUint(envRL, 10, 32)

		if err == nil {
			p.RateLimit = uint(intRL)
		}
	}

	return
}

type ServerParameters struct {
	FlagRunAddr, FileStoragePath string
	DataBaseDSN, Key             string
	StoreInterval                uint
	Restore                      bool
}

func ParseFlagsServer() (p ServerParameters) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.StringVar(&p.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	f.StringVar(&p.Key, "k", "", "hash key")
	f.StringVar(&p.FileStoragePath, "f", "/tmp/metrics-db.json", "path to save storage")
	f.StringVar(
		&p.DataBaseDSN,
		"d",
		"",
		"connection string to database",
	)
	f.UintVar(&p.StoreInterval, "i", 300, "interval in seconds for save storage")
	f.BoolVar(&p.Restore, "r", true, "flag for upload storage from file")
	f.Parse(os.Args[1:])

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		p.FlagRunAddr = envAddr
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		p.Key = envKey
	}

	if envSP := os.Getenv("FILE_STORAGE_PATH"); envSP != "" {
		p.FileStoragePath = envSP
	}

	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		p.DataBaseDSN = envDB
	}

	if envSI := os.Getenv("STORE_INTERVAL"); envSI != "" {
		if unitSI, err := strconv.ParseUint(envSI, 10, 32); err == nil {
			p.StoreInterval = uint(unitSI)
		}
	}

	if envR := os.Getenv("RESTORE"); envR != "" {
		if boolR, err := strconv.ParseBool(envR); err == nil {
			p.Restore = boolR
		}
	}

	return
}
