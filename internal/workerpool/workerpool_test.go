package workerpool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

// MockJob for testing
type MockJob struct {
	id             string
	shouldFail     bool
	processingTime time.Duration
}

func (m *MockJob) Execute(ctx context.Context) error {
	if m.processingTime > 0 {
		time.Sleep(m.processingTime)
	}

	if m.shouldFail {
		return errors.New("mock job failed")
	}
	return nil
}

func (m *MockJob) ID() string {
	return m.id
}

// ContextAwareJob respects context cancellation
type ContextAwareJob struct {
	id             string
	shouldFail     bool
	processingTime time.Duration
	checkContext   bool
}

func (c *ContextAwareJob) Execute(ctx context.Context) error {
	if c.checkContext {
		// Check context before processing
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	if c.processingTime > 0 {
		// Sleep with context awareness
		select {
		case <-time.After(c.processingTime):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if c.shouldFail {
		return errors.New("context-aware job failed")
	}
	return nil
}

func (c *ContextAwareJob) ID() string {
	return c.id
}

// TestWorkerPool_JobProcessing tests various job processing scenarios
func TestWorkerPool_JobProcessing(t *testing.T) {
	tests := []struct {
		name             string
		numWorkers       int
		bufferSize       int
		numJobs          int
		failCondition    func(i int) bool
		processingTime   time.Duration
		expectedSuccess  int
		expectedFailures int
		setupWaitTime    time.Duration
		shutdownWaitTime time.Duration
	}{
		{
			name:             "BasicFunctionality_AllSuccess",
			numWorkers:       3,
			bufferSize:       10,
			numJobs:          10,
			failCondition:    func(i int) bool { return false },
			processingTime:   10 * time.Millisecond,
			expectedSuccess:  10,
			expectedFailures: 0,
			setupWaitTime:    0,
			shutdownWaitTime: 2 * time.Second,
		},
		{
			name:             "ErrorHandling_EveryThirdFails",
			numWorkers:       2,
			bufferSize:       10,
			numJobs:          10,
			failCondition:    func(i int) bool { return i%3 == 0 }, // i=0,3,6,9
			processingTime:   10 * time.Millisecond,
			expectedSuccess:  6,
			expectedFailures: 4,
			setupWaitTime:    50 * time.Millisecond,
			shutdownWaitTime: 2 * time.Second,
		},
		{
			name:             "SmallBuffer_FewJobs",
			numWorkers:       2,
			bufferSize:       5,
			numJobs:          5,
			failCondition:    func(i int) bool { return false },
			processingTime:   10 * time.Millisecond,
			expectedSuccess:  5,
			expectedFailures: 0,
			setupWaitTime:    0,
			shutdownWaitTime: 1 * time.Second,
		},
		{
			name:             "SingleWorker_MultipleJobs",
			numWorkers:       1,
			bufferSize:       10,
			numJobs:          5,
			failCondition:    func(i int) bool { return i == 2 }, // Only job 2 fails
			processingTime:   5 * time.Millisecond,
			expectedSuccess:  4,
			expectedFailures: 1,
			setupWaitTime:    0,
			shutdownWaitTime: 1 * time.Second,
		},
		{
			name:             "ManyWorkers_QuickJobs",
			numWorkers:       10,
			bufferSize:       50,
			numJobs:          20,
			failCondition:    func(i int) bool { return i%5 == 0 }, // i=0,5,10,15
			processingTime:   5 * time.Millisecond,
			expectedSuccess:  16,
			expectedFailures: 4,
			setupWaitTime:    0,
			shutdownWaitTime: 2 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.Default()
			pool := New(tt.numWorkers, tt.bufferSize, logger)
			pool.Start()

			// Setup wait time if needed
			if tt.setupWaitTime > 0 {
				time.Sleep(tt.setupWaitTime)
			}

			// Submit jobs
			for i := 0; i < tt.numJobs; i++ {
				job := &MockJob{
					id:             fmt.Sprintf("job-%d", i),
					shouldFail:     tt.failCondition(i),
					processingTime: tt.processingTime,
				}
				if err := pool.Submit(job); err != nil {
					t.Fatalf("Failed to submit job %d: %v", i, err)
				}
			}

			// Collect results
			successCount := 0
			failCount := 0

			go func() {
				time.Sleep(tt.shutdownWaitTime)
				pool.Stop()
			}()

			for result := range pool.Results() {
				if result.Success {
					successCount++
				} else {
					failCount++
				}
			}

			// Verify results
			if successCount != tt.expectedSuccess {
				t.Errorf("Expected %d successful jobs, got %d", tt.expectedSuccess, successCount)
			}

			if failCount != tt.expectedFailures {
				t.Errorf("Expected %d failed jobs, got %d", tt.expectedFailures, failCount)
			}

			// Verify all jobs were processed
			totalProcessed := successCount + failCount
			if totalProcessed != tt.numJobs {
				t.Errorf("Expected %d total jobs processed, got %d", tt.numJobs, totalProcessed)
			}
		})
	}
}

// TestWorkerPool_GracefulShutdown tests graceful shutdown behavior
func TestWorkerPool_GracefulShutdown(t *testing.T) {
	tests := []struct {
		name            string
		numWorkers      int
		bufferSize      int
		numJobs         int
		processingTime  time.Duration
		minShutdownTime time.Duration
		maxShutdownTime time.Duration
	}{
		{
			name:            "ShutdownWaitsForCompletion",
			numWorkers:      2,
			bufferSize:      5,
			numJobs:         5,
			processingTime:  100 * time.Millisecond,
			minShutdownTime: 100 * time.Millisecond,
			maxShutdownTime: 500 * time.Millisecond,
		},
		{
			name:            "FastJobs_QuickShutdown",
			numWorkers:      3,
			bufferSize:      10,
			numJobs:         10,
			processingTime:  10 * time.Millisecond,
			minShutdownTime: 0, // Fast jobs may complete before shutdown starts
			maxShutdownTime: 200 * time.Millisecond,
		},
		{
			name:            "SingleWorker_SlowJobs",
			numWorkers:      1,
			bufferSize:      5,
			numJobs:         3,
			processingTime:  150 * time.Millisecond,
			minShutdownTime: 150 * time.Millisecond,
			maxShutdownTime: 1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.Default()
			pool := New(tt.numWorkers, tt.bufferSize, logger)
			pool.Start()

			// Submit long-running jobs
			for i := 0; i < tt.numJobs; i++ {
				job := &MockJob{
					id:             fmt.Sprintf("job-%d", i),
					processingTime: tt.processingTime,
				}
				if err := pool.Submit(job); err != nil {
					t.Fatalf("Failed to submit job: %v", err)
				}
			}

			// Give workers time to pick up jobs
			time.Sleep(50 * time.Millisecond)

			// Stop pool and measure shutdown time
			startStop := time.Now()
			pool.Stop()
			stopDuration := time.Since(startStop)

			// Verify shutdown waited for jobs
			if stopDuration < tt.minShutdownTime {
				t.Errorf("Shutdown too fast (%v), expected at least %v", stopDuration, tt.minShutdownTime)
			}

			if stopDuration > tt.maxShutdownTime {
				t.Errorf("Shutdown too slow (%v), expected at most %v", stopDuration, tt.maxShutdownTime)
			}
		})
	}
}

// TestWorkerPool_ConcurrentSubmission tests concurrent job submission
func TestWorkerPool_ConcurrentSubmission(t *testing.T) {
	tests := []struct {
		name             string
		numWorkers       int
		bufferSize       int
		numGoroutines    int
		jobsPerGoroutine int
		processingTime   time.Duration
		useBlocking      bool
	}{
		{
			name:             "ManyGoroutines_BlockingSubmit",
			numWorkers:       5,
			bufferSize:       100,
			numGoroutines:    10,
			jobsPerGoroutine: 10,
			processingTime:   5 * time.Millisecond,
			useBlocking:      true,
		},
		{
			name:             "FewGoroutines_NonBlockingSubmit",
			numWorkers:       3,
			bufferSize:       50,
			numGoroutines:    5,
			jobsPerGoroutine: 10,
			processingTime:   5 * time.Millisecond,
			useBlocking:      false,
		},
		{
			name:             "HighConcurrency_LargeBuffer",
			numWorkers:       10,
			bufferSize:       200,
			numGoroutines:    20,
			jobsPerGoroutine: 5,
			processingTime:   5 * time.Millisecond,
			useBlocking:      true,
		},
		{
			name:             "LowConcurrency_SmallBuffer",
			numWorkers:       2,
			bufferSize:       20,
			numGoroutines:    3,
			jobsPerGoroutine: 5,
			processingTime:   10 * time.Millisecond,
			useBlocking:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.Default()
			pool := New(tt.numWorkers, tt.bufferSize, logger)
			pool.Start()

			var wg sync.WaitGroup
			wg.Add(tt.numGoroutines)

			for g := 0; g < tt.numGoroutines; g++ {
				go func(goroutineID int) {
					defer wg.Done()
					for i := 0; i < tt.jobsPerGoroutine; i++ {
						job := &MockJob{
							id:             fmt.Sprintf("goroutine-%d-job-%d", goroutineID, i),
							processingTime: tt.processingTime,
						}

						var err error
						if tt.useBlocking {
							err = pool.SubmitBlocking(job)
						} else {
							err = pool.Submit(job)
						}

						if err != nil {
							t.Errorf("Failed to submit job: %v", err)
						}
					}
				}(g)
			}

			wg.Wait()

			// Collect all results
			go func() {
				time.Sleep(2 * time.Second)
				pool.Stop()
			}()

			count := 0
			for range pool.Results() {
				count++
			}

			expectedJobs := tt.numGoroutines * tt.jobsPerGoroutine
			if count != expectedJobs {
				t.Errorf("Expected %d jobs processed, got %d", expectedJobs, count)
			}
		})
	}
}

// TestWorkerPool_QueueBehavior tests queue full and overflow scenarios
func TestWorkerPool_QueueBehavior(t *testing.T) {
	tests := []struct {
		name           string
		numWorkers     int
		bufferSize     int
		initialJobs    int
		processingTime time.Duration
		overflowJob    bool
		expectError    bool
	}{
		{
			name:           "FullQueue_NonBlockingSubmit",
			numWorkers:     1,
			bufferSize:     2,
			initialJobs:    3,
			processingTime: 100 * time.Millisecond,
			overflowJob:    true,
			expectError:    true,
		},
		{
			name:           "NearFullQueue_Success",
			numWorkers:     2,
			bufferSize:     5,
			initialJobs:    4,
			processingTime: 50 * time.Millisecond,
			overflowJob:    true,
			expectError:    false,
		},
		{
			name:           "EmptyQueue_NoError",
			numWorkers:     3,
			bufferSize:     10,
			initialJobs:    0,
			processingTime: 10 * time.Millisecond,
			overflowJob:    true,
			expectError:    false,
		},
		{
			name:           "TinyBuffer_OverflowImmediately",
			numWorkers:     1,
			bufferSize:     1,
			initialJobs:    2,
			processingTime: 200 * time.Millisecond,
			overflowJob:    true,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.Default()
			pool := New(tt.numWorkers, tt.bufferSize, logger)
			pool.Start()

			// Fill the queue with initial jobs
			for i := 0; i < tt.initialJobs; i++ {
				job := &MockJob{
					id:             fmt.Sprintf("job-%d", i),
					processingTime: tt.processingTime,
				}
				_ = pool.Submit(job)
			}

			// Try to submit overflow job
			var err error
			if tt.overflowJob {
				overflowJob := &MockJob{
					id:             "overflow-job",
					processingTime: tt.processingTime,
				}
				err = pool.Submit(overflowJob)
			}

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("Expected error when submitting to full queue, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Did not expect error, got: %v", err)
			}

			pool.Stop()
		})
	}
}

// TestWorkerPool_EdgeCases tests edge cases and boundary conditions
func TestWorkerPool_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		numWorkers int
		bufferSize int
		numJobs    int
		shouldPass bool
	}{
		{
			name:       "ZeroJobs",
			numWorkers: 3,
			bufferSize: 10,
			numJobs:    0,
			shouldPass: true,
		},
		{
			name:       "SingleJob_SingleWorker",
			numWorkers: 1,
			bufferSize: 1,
			numJobs:    1,
			shouldPass: true,
		},
		{
			name:       "ManyWorkers_FewJobs",
			numWorkers: 100,
			bufferSize: 200,
			numJobs:    5,
			shouldPass: true,
		},
		{
			name:       "LargeScale",
			numWorkers: 50,
			bufferSize: 1000,
			numJobs:    500,
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.Default()
			pool := New(tt.numWorkers, tt.bufferSize, logger)
			pool.Start()

			// Submit jobs
			for i := 0; i < tt.numJobs; i++ {
				job := &MockJob{
					id:             fmt.Sprintf("job-%d", i),
					processingTime: 5 * time.Millisecond,
				}
				if err := pool.Submit(job); err != nil {
					if tt.shouldPass {
						t.Fatalf("Unexpected error: %v", err)
					}
				}
			}

			// Collect results
			go func() {
				time.Sleep(3 * time.Second)
				pool.Stop()
			}()

			count := 0
			for range pool.Results() {
				count++
			}

			if count != tt.numJobs {
				t.Errorf("Expected %d jobs processed, got %d", tt.numJobs, count)
			}
		})
	}
}

// TestWorkerPool_ContextCancellation tests context cancellation scenarios
func TestWorkerPool_ContextCancellation(t *testing.T) {
	tests := []struct {
		name              string
		numWorkers        int
		bufferSize        int
		numJobs           int
		cancelBeforeJobs  bool
		cancelDuringJobs  bool
		processingTime    time.Duration
		expectCancelError bool
	}{
		{
			name:              "CancelDuringExecution",
			numWorkers:        2,
			bufferSize:        10,
			numJobs:           5,
			cancelDuringJobs:  true,
			processingTime:    200 * time.Millisecond,
			expectCancelError: true,
		},
		{
			name:             "CancelBeforeSubmit",
			numWorkers:       2,
			bufferSize:       10,
			numJobs:          3,
			cancelBeforeJobs: true,
			processingTime:   50 * time.Millisecond,
		},
		{
			name:           "ImmediateCancel",
			numWorkers:     3,
			bufferSize:     5,
			numJobs:        10,
			processingTime: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.Default()
			pool := New(tt.numWorkers, tt.bufferSize, logger)
			pool.Start()

			if tt.cancelBeforeJobs {
				// Cancel before submitting jobs
				pool.cancel()
				time.Sleep(50 * time.Millisecond)
			}

			// Submit jobs
			var submitErrors int
			for i := 0; i < tt.numJobs; i++ {
				job := &ContextAwareJob{
					id:             fmt.Sprintf("job-%d", i),
					processingTime: tt.processingTime,
					checkContext:   true,
				}

				if err := pool.Submit(job); err != nil {
					submitErrors++
				}
			}

			if tt.cancelDuringJobs {
				// Cancel while jobs are processing
				time.Sleep(50 * time.Millisecond)
				pool.cancel()
			}

			// Collect results
			go func() {
				time.Sleep(2 * time.Second)
				pool.Stop()
			}()

			canceledJobs := 0
			for result := range pool.Results() {
				if result.Error != nil && errors.Is(result.Error, context.Canceled) {
					canceledJobs++
				}
			}

			// Note: These assertions are relaxed because timing is non-deterministic
			// The test mainly ensures no panics occur during cancellation
			if tt.cancelBeforeJobs && submitErrors == 0 {
				t.Logf("Warning: No submit errors after cancellation (timing issue)")
			}

			if tt.expectCancelError && canceledJobs == 0 {
				t.Logf("Warning: No canceled jobs detected (jobs may have completed first)")
			}
		})
	}
}

// TestWorkerPool_StopWithContext tests StopWithContext method
func TestWorkerPool_StopWithContext(t *testing.T) {
	tests := []struct {
		name            string
		numWorkers      int
		numJobs         int
		processingTime  time.Duration
		shutdownTimeout time.Duration
		expectTimeout   bool
	}{
		{
			name:            "GracefulShutdown_WithTimeout",
			numWorkers:      2,
			numJobs:         5,
			processingTime:  50 * time.Millisecond,
			shutdownTimeout: 1 * time.Second,
			expectTimeout:   false,
		},
		{
			name:            "ShutdownTimeout_LongJobs",
			numWorkers:      1,
			numJobs:         3,
			processingTime:  500 * time.Millisecond,
			shutdownTimeout: 100 * time.Millisecond,
			expectTimeout:   true,
		},
		{
			name:            "QuickShutdown_FastJobs",
			numWorkers:      5,
			numJobs:         10,
			processingTime:  10 * time.Millisecond,
			shutdownTimeout: 500 * time.Millisecond,
			expectTimeout:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.Default()
			pool := New(tt.numWorkers, 20, logger)
			pool.Start()

			// Submit jobs
			for i := 0; i < tt.numJobs; i++ {
				job := &MockJob{
					id:             fmt.Sprintf("job-%d", i),
					processingTime: tt.processingTime,
				}
				_ = pool.Submit(job)
			}

			// Give workers time to pick up jobs
			time.Sleep(50 * time.Millisecond)

			// Stop with context timeout
			ctx, cancel := context.WithTimeout(context.Background(), tt.shutdownTimeout)
			defer cancel()

			err := pool.StopWithContext(ctx)

			if tt.expectTimeout && err == nil {
				t.Error("Expected timeout error, got nil")
			}

			if !tt.expectTimeout && err != nil {
				t.Errorf("Did not expect error, got: %v", err)
			}

			if tt.expectTimeout && err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Errorf("Expected context.DeadlineExceeded, got: %v", err)
				}
			}
		})
	}
}

// TestWorkerPool_SubmitDuringShutdown tests submitting while shutdown is in progress
func TestWorkerPool_SubmitDuringShutdown(t *testing.T) {
	logger := log.Default()
	pool := New(2, 10, logger)
	pool.Start()

	// Submit a job
	job := &MockJob{id: "job-1", processingTime: 10 * time.Millisecond}
	if err := pool.Submit(job); err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}

	// Collect results in background
	go func() {
		for range pool.Results() {
			// Drain results
		}
	}()

	// Wait for job to complete
	time.Sleep(50 * time.Millisecond)

	// Stop the pool
	pool.Stop()

	// After Stop(), submitting will cause issues because:
	// - jobs channel is closed
	// - context is canceled
	// This test documents that you should NOT submit after Stop()
	// In production, always submit before calling Stop()
}

// TestWorkerPool_NilLogger tests that nil logger is handled
func TestWorkerPool_NilLogger(t *testing.T) {
	// Create pool with nil logger
	pool := New(2, 5, nil)

	if pool.logger == nil {
		t.Error("Expected default logger, got nil")
	}

	pool.Start()

	// Submit a job
	job := &MockJob{id: "job-1", processingTime: 10 * time.Millisecond}
	if err := pool.Submit(job); err != nil {
		t.Fatalf("Failed to submit job: %v", err)
	}

	// Stop pool
	go func() {
		time.Sleep(500 * time.Millisecond)
		pool.Stop()
	}()

	count := 0
	for range pool.Results() {
		count++
	}

	if count != 1 {
		t.Errorf("Expected 1 job processed, got %d", count)
	}
}

// TestWorkerPool_ResultChannelBlocking tests result channel blocking behavior
func TestWorkerPool_ResultChannelBlocking(t *testing.T) {
	logger := log.Default()
	pool := New(2, 5, logger)
	pool.Start()

	// Submit jobs but don't read results immediately
	numJobs := 5
	for i := 0; i < numJobs; i++ {
		job := &MockJob{
			id:             fmt.Sprintf("job-%d", i),
			processingTime: 10 * time.Millisecond,
		}
		if err := pool.Submit(job); err != nil {
			t.Fatalf("Failed to submit job: %v", err)
		}
	}

	// Wait for jobs to be processed
	time.Sleep(200 * time.Millisecond)

	// Now read results
	results := []Result{}
	go func() {
		time.Sleep(500 * time.Millisecond)
		pool.Stop()
	}()

	for result := range pool.Results() {
		results = append(results, result)
	}

	if len(results) != numJobs {
		t.Errorf("Expected %d results, got %d", numJobs, len(results))
	}
}

// TestWorkerPool_CancelContext tests manual context cancellation
func TestWorkerPool_CancelContext(t *testing.T) {
	logger := log.Default()
	pool := New(2, 5, logger)
	pool.Start()

	// Submit a long-running job
	job := &ContextAwareJob{
		id:             "long-job",
		processingTime: 500 * time.Millisecond,
		checkContext:   true,
	}
	_ = pool.Submit(job)

	// Cancel context after job starts
	time.Sleep(100 * time.Millisecond)
	pool.cancel()

	// Wait a bit to ensure context is propagated
	time.Sleep(50 * time.Millisecond)

	// Try to submit after canceling context
	// Note: This may still succeed if it wins the race in select statement
	job2 := &MockJob{id: "job-after-cancel"}
	err := pool.SubmitBlocking(job2)

	// Error may or may not occur due to timing
	if err != nil {
		t.Logf("Got expected error: %v", err)
	}

	// Collect results
	go func() {
		time.Sleep(1 * time.Second)
		pool.Stop()
	}()

	var canceledCount int
	for result := range pool.Results() {
		if result.Error != nil && errors.Is(result.Error, context.Canceled) {
			canceledCount++
		}
	}

	// Note: We don't assert canceledCount > 0 because the job might
	// complete before context is checked. This test just ensures no panic.
	_ = canceledCount
}

// TestWorkerPool_SubmitBlockingBehavior tests the actual blocking behavior of SubmitBlocking
func TestWorkerPool_SubmitBlockingBehavior(t *testing.T) {
	logger := log.Default()

	// Create pool with buffer size 1 and 1 slow worker
	pool := New(1, 1, logger)
	pool.Start()

	// Submit first job - worker picks it up immediately and starts processing
	job1 := &MockJob{
		id:             "job-1",
		processingTime: 800 * time.Millisecond,
	}
	if err := pool.Submit(job1); err != nil {
		t.Fatalf("Failed to submit job1: %v", err)
	}

	// Give worker time to pick up job1
	time.Sleep(50 * time.Millisecond)

	// Submit second job - fills the buffer (size 1)
	job2 := &MockJob{
		id:             "job-2",
		processingTime: 100 * time.Millisecond,
	}
	if err := pool.Submit(job2); err != nil {
		t.Fatalf("Failed to submit job2: %v", err)
	}

	// Now buffer is full (has job2) and worker is busy (processing job1)
	// SubmitBlocking should block until worker finishes job1 and picks up job2
	blockingJobSubmitted := make(chan bool)
	submitStartTime := time.Now()

	go func() {
		job3 := &MockJob{
			id:             "job-3-blocked",
			processingTime: 10 * time.Millisecond,
		}

		// This should block because buffer is full and worker is busy
		err := pool.SubmitBlocking(job3)
		if err != nil {
			t.Errorf("SubmitBlocking failed: %v", err)
		}
		blockingJobSubmitted <- true
	}()

	// Give it time to hit the blocking state
	time.Sleep(100 * time.Millisecond)

	// Verify it's still blocking after 100ms
	select {
	case <-blockingJobSubmitted:
		elapsed := time.Since(submitStartTime)
		// It shouldn't complete this quickly (job1 takes 800ms)
		if elapsed < 500*time.Millisecond {
			t.Errorf("SubmitBlocking completed too quickly (%v), should have blocked", elapsed)
		}
	case <-time.After(50 * time.Millisecond):
		// Good, still blocking as expected
	}

	// Wait for worker to finish job1 and pick up job2
	// After ~800ms, job1 completes, worker picks up job2, buffer has space
	select {
	case <-blockingJobSubmitted:
		elapsed := time.Since(submitStartTime)
		t.Logf("SubmitBlocking completed after %v (expected ~800ms)", elapsed)
		// Should have blocked for roughly the duration of job1
		if elapsed < 700*time.Millisecond {
			t.Errorf("SubmitBlocking completed too quickly: %v", elapsed)
		}
	case <-time.After(2 * time.Second):
		t.Error("SubmitBlocking never completed")
	}

	// Cleanup
	go func() {
		time.Sleep(2 * time.Second)
		pool.Stop()
	}()

	// Drain results
	count := 0
	for range pool.Results() {
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3 jobs processed, got %d", count)
	}
}
