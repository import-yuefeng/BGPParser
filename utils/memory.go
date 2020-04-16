package utils

import (
	"runtime"

	log "github.com/sirupsen/logrus"
)

func TraceMemStats() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	log.Printf("Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes) \n", ms.Alloc, ms.HeapIdle, ms.HeapReleased)
}
