package memstats

import (
	"runtime"
	"strconv"
)

type StringMS struct {
	Name, Value string
}

func GetMemStatsForServer(ms *runtime.MemStats) (stringsMS []StringMS) {
	stringsMS = append(stringsMS, StringMS{"Alloc", strconv.FormatUint(ms.Alloc, 10)})
	stringsMS = append(stringsMS, StringMS{"BuckHashSys", strconv.FormatUint(ms.BuckHashSys, 10)})
	stringsMS = append(stringsMS, StringMS{"Frees", strconv.FormatUint(ms.Frees, 10)})
	stringsMS = append(stringsMS, StringMS{"GCCPUFraction", strconv.FormatFloat(ms.GCCPUFraction, 'f', -1, 64)})
	stringsMS = append(stringsMS, StringMS{"GCSys", strconv.FormatUint(ms.GCSys, 10)})
	stringsMS = append(stringsMS, StringMS{"HeapAlloc", strconv.FormatUint(ms.HeapAlloc, 10)})
	stringsMS = append(stringsMS, StringMS{"HeapIdle", strconv.FormatUint(ms.HeapIdle, 10)})
	stringsMS = append(stringsMS, StringMS{"HeapInuse", strconv.FormatUint(ms.HeapInuse, 10)})
	stringsMS = append(stringsMS, StringMS{"HeapObjects", strconv.FormatUint(ms.HeapObjects, 10)})
	stringsMS = append(stringsMS, StringMS{"HeapReleased", strconv.FormatUint(ms.HeapReleased, 10)})
	stringsMS = append(stringsMS, StringMS{"HeapSys", strconv.FormatUint(ms.HeapSys, 10)})
	stringsMS = append(stringsMS, StringMS{"LastGC", strconv.FormatUint(ms.LastGC, 10)})
	stringsMS = append(stringsMS, StringMS{"Lookups", strconv.FormatUint(ms.Lookups, 10)})
	stringsMS = append(stringsMS, StringMS{"MCacheInuse", strconv.FormatUint(ms.MCacheInuse, 10)})
	stringsMS = append(stringsMS, StringMS{"MCacheSys", strconv.FormatUint(ms.MCacheSys, 10)})
	stringsMS = append(stringsMS, StringMS{"MSpanInuse", strconv.FormatUint(ms.MSpanInuse, 10)})
	stringsMS = append(stringsMS, StringMS{"MSpanSys", strconv.FormatUint(ms.MSpanSys, 10)})
	stringsMS = append(stringsMS, StringMS{"Mallocs", strconv.FormatUint(ms.Mallocs, 10)})
	stringsMS = append(stringsMS, StringMS{"NextGC", strconv.FormatUint(ms.NextGC, 10)})
	stringsMS = append(stringsMS, StringMS{"NumForcedGC", strconv.FormatUint(uint64(ms.NumForcedGC), 10)})
	stringsMS = append(stringsMS, StringMS{"NumGC", strconv.FormatUint(uint64(ms.NumGC), 10)})
	stringsMS = append(stringsMS, StringMS{"OtherSys", strconv.FormatUint(ms.OtherSys, 10)})
	stringsMS = append(stringsMS, StringMS{"PauseTotalNs", strconv.FormatUint(ms.PauseTotalNs, 10)})
	stringsMS = append(stringsMS, StringMS{"StackInuse", strconv.FormatUint(ms.StackInuse, 10)})
	stringsMS = append(stringsMS, StringMS{"StackSys", strconv.FormatUint(ms.StackSys, 10)})
	stringsMS = append(stringsMS, StringMS{"Sys", strconv.FormatUint(ms.Sys, 10)})
	stringsMS = append(stringsMS, StringMS{"TotalAlloc", strconv.FormatUint(ms.TotalAlloc, 10)})

	return
}
