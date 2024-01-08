package parameters

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlagsAgent(t *testing.T) {
	tests := []struct {
		name               string
		f                  func()
		wantListenAddr     string
		wantReportInterval uint
		wantPollInterval   uint
	}{
		{
			name:               "test env",
			f:                  setEnv,
			wantListenAddr:     "testEnv",
			wantReportInterval: 10,
			wantPollInterval:   10,
		},
		{
			name:               "test flags",
			f:                  setFlags,
			wantListenAddr:     "testFlags",
			wantReportInterval: 100,
			wantPollInterval:   100,
		},
		{
			name:               "test default",
			f:                  nil,
			wantListenAddr:     "localhost:8080",
			wantReportInterval: 10,
			wantPollInterval:   2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.f != nil {
				tt.f()
			}

			gotListenAddr, gotReportInterval, gotPollInterval := ParseFlagsAgent()

			assert.Equal(t, tt.wantListenAddr, gotListenAddr)
			assert.Equal(t, tt.wantReportInterval, gotReportInterval)
			assert.Equal(t, tt.wantPollInterval, gotPollInterval)

			delParameters()
		})
	}
}

func setEnv() {
	os.Setenv("ADDRESS", "testEnv")
	os.Setenv("REPORT_INTERVAL", "10")
	os.Setenv("POLL_INTERVAL", "10")
}

func setFlags() {
	os.Args = []string{"test", "-a=testFlags", "-r=100", "-p=100"}
}

func delParameters() {
	os.Clearenv()
	os.Args = []string{"test"}
}

func TestParseFlagsServer(t *testing.T) {
	tests := []struct {
		name            string
		f               func()
		wantFlagRunAddr string
	}{
		{
			name:            "test env",
			f:               setEnv,
			wantFlagRunAddr: "testEnv",
		},
		{
			name:            "test flags",
			f:               setFlags,
			wantFlagRunAddr: "testFlags",
		},
		{
			name:            "test default",
			f:               nil,
			wantFlagRunAddr: "localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.f != nil {
				tt.f()
			}

			gotFlagRunAddr := ParseFlagsServer()

			assert.Equal(t, tt.wantFlagRunAddr, gotFlagRunAddr)
			delParameters()
		})
	}
}
