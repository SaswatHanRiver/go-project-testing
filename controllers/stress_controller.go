package controllers

import (
	"math"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go-project-testing/utils"

	"github.com/gin-gonic/gin"
)

type StressController struct{}

func NewStressController() *StressController {
	return &StressController{}
}

// StressResult - what we return to show Go's concurrency power
type StressResult struct {
	TotalGoroutines       int     `json:"total_goroutines"`
	CompletedJobs         int64   `json:"completed_jobs"`
	DurationMs            int64   `json:"duration_ms"`
	JobsPerSecond         int64   `json:"jobs_per_second"`
	MemoryBeforeMB        float64 `json:"memory_before_mb"`
	MemoryAfterMB         float64 `json:"memory_after_mb"`
	GoRoutineMemoryUsedMB float64 `json:"goroutine_memory_used_mb"`
	AvgMemoryPerGoroutineKB float64 `json:"avg_memory_per_goroutine_kb"`
	CPUCores              int     `json:"cpu_cores"`
	JavaEquivalentMemoryGB float64 `json:"java_equivalent_memory_gb"`
	Summary               string  `json:"summary"`
}

// RunStressTest godoc
// @Summary      Goroutine concurrency stress test
// @Description  Spawns N goroutines simultaneously and shows memory + throughput vs Java
// @Tags         stress
// @Produce      json
// @Security     BearerAuth
// @Param        jobs  query     int  false  "Number of goroutines to spawn (default 20000)"
// @Success      200   {object}  utils.ApiResponse
// @Router       /api/stress [get]
func (c *StressController) RunStressTest(ctx *gin.Context) {
	// Read jobs count from query param - default 20000 (like LG U+ X-Ray TPS requirement)
	jobCount := 20000
	if q := ctx.Query("jobs"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 && n <= 100000 {
			jobCount = n
		}
	}

	// --- Memory BEFORE ---
	var memBefore runtime.MemStats
	runtime.GC() // clean GC before measuring
	runtime.ReadMemStats(&memBefore)

	var completed int64 // atomic counter - safe to increment from multiple goroutines
	var wg sync.WaitGroup // like CountDownLatch in Java

	start := time.Now()

	// Spawn N goroutines simultaneously
	// This is the key: spawning 20,000 goroutines costs almost nothing
	for i := 0; i < jobCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Simulate real work: small CPU computation (not just sleep)
			// Like processing a request, doing a calculation, or parsing data
			result := 0.0
			for j := 0; j < 100; j++ {
				result += math.Sqrt(float64(id * j))
			}
			_ = result // prevent compiler from optimizing away

			atomic.AddInt64(&completed, 1) // thread-safe increment
		}(i)
	}

	wg.Wait() // block until ALL goroutines finish - like CountDownLatch.await()
	duration := time.Since(start)

	// --- Memory AFTER ---
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// Calculate stats
	memBeforeMB := float64(memBefore.Alloc) / 1024 / 1024
	memAfterMB := float64(memAfter.TotalAlloc) / 1024 / 1024
	usedMB := memAfterMB - memBeforeMB
	if usedMB < 0 {
		usedMB = float64(memAfter.Alloc) / 1024 / 1024
	}
	avgKBPerGoroutine := (usedMB * 1024) / float64(jobCount)
	durationMs := duration.Milliseconds()
	if durationMs == 0 {
		durationMs = 1
	}
	jobsPerSec := (completed * 1000) / durationMs

	// Java comparison: 1 OS thread = ~1MB stack minimum
	javaEquivGB := float64(jobCount) / 1024.0

	result := StressResult{
		TotalGoroutines:         jobCount,
		CompletedJobs:           completed,
		DurationMs:              durationMs,
		JobsPerSecond:           jobsPerSec,
		MemoryBeforeMB:          math.Round(memBeforeMB*100) / 100,
		MemoryAfterMB:           math.Round(memAfterMB*100) / 100,
		GoRoutineMemoryUsedMB:   math.Round(usedMB*100) / 100,
		AvgMemoryPerGoroutineKB: math.Round(avgKBPerGoroutine*100) / 100,
		CPUCores:                runtime.NumCPU(),
		JavaEquivalentMemoryGB:  math.Round(javaEquivGB*100) / 100,
		Summary: "Go handled " + strconv.Itoa(jobCount) + " concurrent goroutines using ~" +
			strconv.FormatFloat(math.Round(usedMB*100)/100, 'f', 2, 64) + "MB. " +
			"Java would need ~" + strconv.FormatFloat(javaEquivGB, 'f', 1, 64) + "GB for the same concurrency using OS threads.",
	}

	utils.Success(ctx, http.StatusOK, "Stress test completed", result)
}
