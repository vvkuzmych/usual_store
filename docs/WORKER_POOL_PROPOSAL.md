# Worker Pool Implementation Proposal

## 📋 Overview
Add worker pool pattern to improve concurrent processing in `usual_store` project.

---

## 🎯 Where to Use Worker Pool

### 1. **Messaging Service (HIGH PRIORITY)** ⭐⭐⭐
**File:** `cmd/messaging-service/main.go` + `internal/messaging/consumer.go`

**Problem:**
- Currently processes Kafka messages sequentially (one at a time)
- Slow email sending blocks other messages
- No concurrent processing

**Solution:**
- Worker pool with 5-10 workers
- Each worker sends emails concurrently
- Better throughput for high email volume

**Benefits:**
- ✅ 10x faster email processing
- ✅ Handle burst traffic
- ✅ Graceful degradation under load

---

### 2. **Batch Operations (MEDIUM PRIORITY)** ⭐⭐
**Potential locations:**
- Bulk customer imports
- Invoice generation
- Report processing
- Data exports

**Use when:**
- Processing multiple items (>100)
- Each item independent
- I/O-bound operations

---

## 📦 Implementation Plan

### Step 1: Create Worker Pool Package
**Location:** `internal/workerpool/workerpool.go`

**Features:**
- Generic worker pool
- Configurable worker count
- Graceful shutdown
- Error handling
- Metrics (optional)

---

### Step 2: Update Messaging Service

**Changes needed:**

1. **Add worker pool to consumer:**
```go
// internal/messaging/consumer.go
type Consumer struct {
    reader     *kafka.Reader
    logger     *log.Logger
    handler    EmailHandler
    workerPool *workerpool.Pool  // ← NEW
}
```

2. **Process messages concurrently:**
```go
// Submit to worker pool instead of direct processing
func (c *Consumer) Start(ctx context.Context) error {
    pool := workerpool.New(10) // 10 workers
    pool.Start()
    defer pool.Stop()
    
    for {
        msg, err := c.reader.ReadMessage(ctx)
        // Submit job to pool (non-blocking)
        pool.Submit(Job{Message: msg})
    }
}
```

---

## 🗂️ Files to Create/Modify

### **New Files:**
1. `internal/workerpool/workerpool.go` - Main worker pool implementation
2. `internal/workerpool/workerpool_test.go` - Tests
3. `internal/workerpool/job.go` - Job interface/types
4. `cmd/messaging-service/worker_handler.go` - Email worker handler

### **Modified Files:**
1. `internal/messaging/consumer.go` - Add worker pool integration
2. `cmd/messaging-service/main.go` - Configure worker count
3. `cmd/messaging-service/email_handler.go` - Make thread-safe if needed

---

## 🎛️ Configuration

**Add to `.env`:**
```bash
# Worker Pool Settings
EMAIL_WORKER_COUNT=10        # Number of concurrent email workers
EMAIL_WORKER_BUFFER=100      # Job queue buffer size
EMAIL_MAX_RETRIES=3          # Retry failed emails
```

**Add to config struct:**
```go
type config struct {
    // ... existing fields
    workers struct {
        count      int
        bufferSize int
        maxRetries int
    }
}
```

---

## 📊 Expected Results

### **Before (Current):**
```
Kafka → Consumer → Email Send (sequential)
         ↓ (blocks)
       Next message waits
       
Throughput: ~10 emails/second
```

### **After (With Worker Pool):**
```
Kafka → Consumer → Worker Pool → [Worker 1, 2, 3...10] → SMTP
                                   (parallel sending)
                                   
Throughput: ~100 emails/second (10x improvement)
```

---

## ⚙️ Implementation Details

### **Worker Pool Structure:**
```go
type Pool struct {
    numWorkers int
    jobs       chan Job
    results    chan Result
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

type Job struct {
    Message kafka.Message
    Handler EmailHandler
}

type Result struct {
    JobID   string
    Success bool
    Error   error
}
```

### **Key Features:**
- ✅ Bounded concurrency (configurable workers)
- ✅ Buffered job queue (backpressure)
- ✅ Graceful shutdown (wait for current jobs)
- ✅ Error handling (collect failed jobs)
- ✅ Context cancellation support

---

## 🧪 Testing Strategy

1. **Unit tests:** `workerpool_test.go`
   - Test worker pool creation
   - Test job submission
   - Test graceful shutdown

2. **Integration tests:** Test with real Kafka + SMTP
   - Process 1000 messages
   - Measure throughput improvement
   - Test error scenarios

3. **Load tests:**
   - Simulate high message volume
   - Monitor resource usage
   - Verify no memory leaks

---

## 📈 Monitoring (Optional)

**Add metrics:**
- Jobs processed per second
- Active workers count
- Queue length
- Success/failure rate
- Average processing time

**Integration:**
- Prometheus metrics
- Grafana dashboards
- Alert on failures

---

## 🚀 Rollout Plan

### **Phase 1: Core Implementation (Week 1)**
- Create `internal/workerpool` package
- Write tests
- Document usage

### **Phase 2: Messaging Service Integration (Week 2)**
- Update consumer to use worker pool
- Add configuration
- Test with staging environment

### **Phase 3: Production Deployment (Week 3)**
- Start with 5 workers
- Monitor performance
- Gradually increase to 10 workers
- Collect metrics

---

## 🎯 Success Criteria

- ✅ Email throughput increased 5-10x
- ✅ No message loss
- ✅ Graceful shutdown works
- ✅ Error handling robust
- ✅ Resource usage acceptable (<500MB memory)
- ✅ Tests pass (>80% coverage)

---

## ❓ Questions to Resolve

1. **Worker count:** Start with 5, 10, or configurable?
   - **Recommendation:** Start with 10, make configurable

2. **Retry strategy:** Immediate, exponential backoff?
   - **Recommendation:** 3 retries with 1s, 5s, 10s delays

3. **Failed messages:** Dead letter queue or log?
   - **Recommendation:** Log errors + optional DLQ

4. **Monitoring:** Basic logs or full metrics?
   - **Recommendation:** Start with logs, add metrics later

---

## 📝 Example Usage (After Implementation)

```go
// cmd/messaging-service/main.go
func main() {
    // ... existing setup
    
    // Create worker pool
    pool := workerpool.New(cfg.workers.count)
    pool.Start()
    defer pool.Stop()
    
    // Create consumer with worker pool
    consumer := messaging.NewConsumerWithPool(
        cfg.kafka.brokers,
        cfg.kafka.topic,
        cfg.kafka.groupID,
        emailHandler,
        pool,  // ← Pass worker pool
        infoLog,
    )
    
    // Start consuming (non-blocking)
    consumer.Start(ctx)
}
```

---

## 🔧 Alternatives Considered

1. **Go routines per message:** Too many goroutines → memory issues
2. **Semaphore only:** No job queue → messages dropped under load
3. **External queue (Redis):** Adds complexity, overkill for now

**Decision:** Worker pool is the right balance of simplicity + performance.

---

## 📚 References

- `/Users/vkuzm/GolandProjects/golang_practice/week_26/concurrency_patterns/worker_pool/`
- `/Users/vkuzm/GolandProjects/golang_practice/week_26/concurrency_patterns/top-patterns-comparison.md`

---

## ✅ Approval Needed

Please review and approve/suggest changes:

- [ ] Approve overall approach
- [ ] Confirm messaging service as priority
- [ ] Agree on worker count (10?)
- [ ] Confirm file structure
- [ ] Approve configuration approach
- [ ] Any other use cases to add?

**After approval, I'll create the implementation!** 🚀
