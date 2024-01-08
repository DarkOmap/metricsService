package memstats

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMemStatsForServer(t *testing.T) {
	var ms runtime.MemStats
	GetMemStatsForServer(&ms)
	tests := []struct {
		name          string
		wantStringsMS []StringMS
	}{
		{
			name:          "positive test",
			wantStringsMS: getWantStringsMS(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStringsMS := GetMemStatsForServer(&ms)
			assert.ElementsMatch(t, tt.wantStringsMS, gotStringsMS)
		})
	}
}

func getWantStringsMS() []StringMS {
	return []StringMS{
		{"Alloc", "0"},
		{"BuckHashSys", "0"},
		{"Frees", "0"},
		{"GCCPUFraction", "0"},
		{"GCSys", "0"},
		{"HeapAlloc", "0"},
		{"HeapIdle", "0"},
		{"HeapInuse", "0"},
		{"HeapObjects", "0"},
		{"HeapReleased", "0"},
		{"HeapSys", "0"},
		{"LastGC", "0"},
		{"Lookups", "0"},
		{"MCacheInuse", "0"},
		{"MCacheSys", "0"},
		{"MSpanInuse", "0"},
		{"MSpanSys", "0"},
		{"Mallocs", "0"},
		{"NextGC", "0"},
		{"NumForcedGC", "0"},
		{"NumGC", "0"},
		{"OtherSys", "0"},
		{"PauseTotalNs", "0"},
		{"StackInuse", "0"},
		{"StackSys", "0"},
		{"Sys", "0"},
		{"TotalAlloc", "0"},
	}

}
