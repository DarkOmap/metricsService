package memstats

import (
	"runtime"
	"testing"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMemStatsForServer(t *testing.T) {
	var ms runtime.MemStats
	GetMemStatsForServer(&ms)
	tests := []struct {
		name string
		want map[string]float64
	}{
		{
			name: "positive test",
			want: getWantStringsMS(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStringsMS := GetMemStatsForServer(&ms)
			assert.Equal(t, tt.want, gotStringsMS)
		})
	}
}

func getWantStringsMS() map[string]float64 {
	return map[string]float64{
		"Alloc":         0,
		"BuckHashSys":   0,
		"Frees":         0,
		"GCCPUFraction": 0,
		"GCSys":         0,
		"HeapAlloc":     0,
		"HeapIdle":      0,
		"HeapInuse":     0,
		"HeapObjects":   0,
		"HeapReleased":  0,
		"HeapSys":       0,
		"LastGC":        0,
		"Lookups":       0,
		"MCacheInuse":   0,
		"MCacheSys":     0,
		"MSpanInuse":    0,
		"MSpanSys":      0,
		"Mallocs":       0,
		"NextGC":        0,
		"NumForcedGC":   0,
		"NumGC":         0,
		"OtherSys":      0,
		"PauseTotalNs":  0,
		"StackInuse":    0,
		"StackSys":      0,
		"Sys":           0,
		"TotalAlloc":    0,
	}

}

func TestGetVirtualMemoryForServer(t *testing.T) {
	vm, err := mem.VirtualMemory()
	require.NoError(t, err)
	t.Run("positive test", func(t *testing.T) {
		gotStringsVM := GetVirtualMemoryForServer(vm)
		assert.Equal(t, map[string]float64{
			"TotalMemory": float64(vm.Total),
			"FreeMemory":  float64(vm.Free),
		}, gotStringsVM)
	})
}
