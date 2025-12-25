# OpenTelemetry Tracing Guide for Usual Store

## 🔍 Overview

OpenTelemetry provides **distributed tracing** to help you understand and debug your application's request flows. This guide shows you what tracing will look like once it's enabled.

## ⏸️ Current Status

**Status:** Temporarily disabled due to TLS certificate verification issues on your system.

**Error:** `x509: OSStatus -26276`

**Files Created:**
- ✅ `internal/telemetry/telemetry.go` - Core tracing setup
- ✅ `internal/telemetry/database.go` - Database query tracing
- ✅ `internal/telemetry/stripe.go` - Stripe API tracing
- ✅ `cmd/api/api.go` - Integration (commented out)
- ✅ `cmd/api/routes-api.go` - HTTP middleware (commented out)
- ✅ `docker-compose.yml` - Jaeger service added

---

## 🚀 How to Enable (Once Certificate Issue is Fixed)

### Step 1: Fix Certificate Issue

Option A: Check System Certificates
```bash
# macOS - check certificate validity
security find-certificate -a -p /System/Library/Keychains/SystemRootCertificates.keychain

# Update certificates if needed
```

Option B: Use GOSUMDB=off (Development Only)
```bash
export GOSUMDB=off
go mod tidy
```

Option C: Use go mod vendor
```bash
go mod vendor
go build -mod=vendor
```

### Step 2: Uncomment OpenTelemetry Code

**File: `cmd/api/api.go`**
```go
// Uncomment this import:
// import "usual_store/internal/telemetry"

// Uncomment this section (around line 162):
// Initialize OpenTelemetry if enabled
var telemetryShutdown func(context.Context) error
if os.Getenv("OTEL_ENABLED") == "true" {
    otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
    if otelEndpoint == "" {
        otelEndpoint = "localhost:4318"
    }

    otelCfg := telemetry.Config{
        ServiceName:    getEnvOrDefault("OTEL_SERVICE_NAME", "usual-store-api"),
        ServiceVersion: getEnvOrDefault("OTEL_SERVICE_VERSION", version),
        Environment:    getEnvOrDefault("OTEL_ENVIRONMENT", cfg.env),
        OTLPEndpoint:   otelEndpoint,
    }

    shutdown, err := telemetry.InitTracer(otelCfg)
    if err != nil {
        errorLog.Printf("Failed to initialize OpenTelemetry: %v", err)
    } else {
        infoLog.Printf("OpenTelemetry initialized successfully (endpoint: %s)", otelEndpoint)
        telemetryShutdown = shutdown
    }
}
```

**File: `cmd/api/routes-api.go`**
```go
// Uncomment these imports:
// import "os"
// import "go.opentelemetry.io/contrib/instrumentation/github.com/go-chi/chi/v5/otelchi"

// Uncomment this middleware:
if os.Getenv("OTEL_ENABLED") == "true" {
    serviceName := os.Getenv("OTEL_SERVICE_NAME")
    if serviceName == "" {
        serviceName = "usual-store-api"
    }
    mux.Use(otelchi.Middleware(serviceName, otelchi.WithChiRoutes(mux)))
}
```

### Step 3: Update go.mod

```bash
# Add the OpenTelemetry packages with correct versions
go get go.opentelemetry.io/otel@latest
go get go.opentelemetry.io/otel/sdk@latest
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp@latest
go get go.opentelemetry.io/contrib/instrumentation/github.com/go-chi/chi/v5/otelchi@latest

go mod tidy
```

### Step 4: Start Jaeger

```bash
# Start Jaeger tracing backend
docker-compose up -d jaeger

# Verify Jaeger is running
docker-compose ps jaeger
```

### Step 5: Enable OpenTelemetry in docker-compose.yml

The environment variables are already configured in `docker-compose.yml`:
```yaml
back-end:
  environment:
    - OTEL_ENABLED=true
    - OTEL_SERVICE_NAME=usual-store-api
    - OTEL_SERVICE_VERSION=1.0.0
    - OTEL_ENVIRONMENT=development
    - OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4318
```

### Step 6: Rebuild and Restart Backend

```bash
docker-compose build back-end
docker-compose up -d back-end
```

### Step 7: Access Jaeger UI

```
IPv4: http://localhost:16686
IPv6: http://[::1]:16686
```

---

## 📊 What You'll See: Example Traces

### Example 1: Widget API Request

When you make a request like `GET /api/widgets/1`, you'll see:

```
┌─────────────────────────────────────────────────────────┐
│ GET /api/widgets/1                      Duration: 15ms  │
├─────────────────────────────────────────────────────────┤
│ ├─ database.query_row                   Duration: 12ms  │
│ │  └─ SELECT * FROM widgets WHERE id=$1               │
│ │     • db.system: postgresql                          │
│ │     • db.statement: SELECT * FROM widgets...         │
│ │     • db.duration_ms: 12                            │
│ └─ http.response                        Duration: 3ms   │
│    └─ status: 200                                      │
└─────────────────────────────────────────────────────────┘
```

**Insights:**
- Total request time: 15ms
- Database query: 12ms (80% of time)
- HTTP processing: 3ms
- ✅ Performance is good!

### Example 2: Payment Intent Creation

When you create a Stripe payment intent:

```
┌─────────────────────────────────────────────────────────┐
│ POST /api/payment-intent                Duration: 450ms │
├─────────────────────────────────────────────────────────┤
│ ├─ database.query_row                   Duration: 12ms  │
│ │  └─ SELECT * FROM users WHERE email=$1              │
│ │     • db.system: postgresql                          │
│ │     • db.duration_ms: 12                            │
│ │                                                       │
│ ├─ stripe.CreatePaymentIntent           Duration: 380ms│
│ │  └─ External API call to Stripe                     │
│ │     • stripe.operation: CreatePaymentIntent          │
│ │     • external.service: stripe                       │
│ │     • stripe.duration_ms: 380                       │
│ │                                                       │
│ ├─ database.exec                        Duration: 8ms  │
│ │  └─ INSERT INTO transactions...                     │
│ │     • db.system: postgresql                          │
│ │     • db.duration_ms: 8                             │
│ │                                                       │
│ └─ http.response                        Duration: 50ms  │
│    └─ status: 200                                      │
└─────────────────────────────────────────────────────────┘
```

**Insights:**
- Total request time: 450ms
- Stripe API: 380ms (84% of time) ⚠️ **Bottleneck!**
- Database queries: 20ms total
- ✅ Database is fast
- ⚠️ Consider caching Stripe responses if possible

### Example 3: Subscription Creation

```
┌─────────────────────────────────────────────────────────┐
│ POST /api/create-customer-and-subscribe  Duration: 950ms│
├─────────────────────────────────────────────────────────┤
│ ├─ database.query_row                   Duration: 10ms  │
│ │  └─ SELECT * FROM users WHERE id=$1                 │
│ │                                                       │
│ ├─ stripe.CreateCustomer                Duration: 420ms│
│ │  └─ External API call to Stripe                     │
│ │     • stripe.operation: CreateCustomer               │
│ │     • stripe.duration_ms: 420                       │
│ │                                                       │
│ ├─ stripe.CreateSubscription            Duration: 450ms│
│ │  └─ External API call to Stripe                     │
│ │     • stripe.operation: CreateSubscription           │
│ │     • stripe.duration_ms: 450                       │
│ │                                                       │
│ ├─ database.exec                        Duration: 15ms  │
│ │  └─ INSERT INTO orders...                           │
│ │                                                       │
│ └─ http.response                        Duration: 55ms  │
│    └─ status: 200                                      │
└─────────────────────────────────────────────────────────┘
```

**Insights:**
- Total: 950ms (almost 1 second)
- Two Stripe API calls: 870ms (91%)
- Database: 25ms total
- ✅ Two external calls are expected
- ℹ️ Consider parallel execution if possible

### Example 4: Error Scenario

When an error occurs:

```
┌─────────────────────────────────────────────────────────┐
│ POST /api/payment-intent                Duration: 25ms  │
│ ❌ ERROR: insufficient funds                           │
├─────────────────────────────────────────────────────────┤
│ ├─ database.query_row                   Duration: 10ms  │
│ │  └─ SELECT * FROM users WHERE id=$1                 │
│ │                                                       │
│ ├─ stripe.CreatePaymentIntent           Duration: 380ms│
│ │  ❌ ERROR                                            │
│ │  └─ error: card_declined                             │
│ │     • error.type: StripeError                       │
│ │     • error.message: Your card was declined          │
│ │                                                       │
│ └─ http.response                        Duration: 5ms   │
│    └─ status: 400                                      │
└─────────────────────────────────────────────────────────┘
```

**Insights:**
- Error clearly visible in trace
- Can see exact point of failure (Stripe)
- Error details captured
- 🔍 Easy to debug!

---

## 🎯 Key Benefits of Tracing

### 1. Performance Monitoring
- **See exact timings** for each operation
- **Identify bottlenecks** immediately
- **Compare** different request types

### 2. Debugging
- **Track errors** through entire request flow
- **See context** of what happened before/after error
- **Reproduce issues** with exact request details

### 3. Dependency Tracking
- **Visualize** how services interact
- **Monitor** external API performance (Stripe, etc.)
- **Detect** slow database queries

### 4. Production Insights
- **Real user requests** traced
- **P99 latency** tracking
- **Failure rate** monitoring

---

## 📈 Jaeger UI Features

### Service Map
```
┌─────────┐     ┌──────────┐     ┌──────────┐
│ Frontend│────▶│  API     │────▶│ Database │
└─────────┘     └──────────┘     └──────────┘
                      │
                      ▼
                ┌──────────┐
                │  Stripe  │
                └──────────┘
```

### Trace Search
- Filter by:
  - Service name
  - Operation name
  - Tags (user_id, status_code, etc.)
  - Duration (find slow requests)
  - Time range

### Trace Details
- Complete request timeline
- All database queries
- All external API calls
- Error details
- Custom tags

---

## 🛠️ Custom Instrumentation

You can add custom tracing to your code:

```go
import "usual_store/internal/telemetry"

func MyFunction(ctx context.Context) error {
    tracer := telemetry.Tracer("my-service")
    ctx, span := tracer.Start(ctx, "MyFunction")
    defer span.End()

    // Your code here
    // Pass ctx to child functions for distributed tracing

    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return err
    }

    span.SetStatus(codes.Ok, "success")
    return nil
}
```

### Add Custom Attributes

```go
span.SetAttributes(
    attribute.String("user.email", email),
    attribute.Int("order.amount", amount),
    attribute.String("product.name", productName),
)
```

---

## 🔧 Configuration

### Environment Variables

```bash
# Enable tracing
OTEL_ENABLED=true

# Service identification
OTEL_SERVICE_NAME=usual-store-api
OTEL_SERVICE_VERSION=1.0.0
OTEL_ENVIRONMENT=production

# Jaeger endpoint
OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4318
```

### Sampling

The current configuration uses `AlwaysSample()` which traces **every request**. For production, you might want to sample:

```go
// In internal/telemetry/telemetry.go
traceProvider := sdktrace.NewTracerProvider(
    // Sample 10% of requests
    sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)),
    // ...
)
```

---

## 📊 Example: Debugging a Slow Request

### Problem
Users report slow checkout (2-3 seconds)

### Steps
1. Open Jaeger UI: http://localhost:16686
2. Select service: `usual-store-api`
3. Filter by operation: `POST /api/payment-intent`
4. Sort by duration: longest first
5. Click on slow trace

### What You Find
```
Total: 2,450ms
  ├─ database.query_row: 15ms ✓ Fast
  ├─ stripe.CreatePaymentIntent: 2,380ms ⚠️ SLOW!
  └─ database.exec: 55ms ⚠️ Also slow
```

### Solution
- **Stripe API is slow**: Check Stripe status page, or add timeout
- **Database insert is slow**: Missing index? Add:
  ```sql
  CREATE INDEX idx_transactions_user_id ON transactions(user_id);
  ```

### Result
After adding index:
```
Total: 420ms (6x faster!)
  ├─ database.query_row: 12ms ✓
  ├─ stripe.CreatePaymentIntent: 395ms ✓ (Stripe's normal)
  └─ database.exec: 13ms ✓ Fixed!
```

---

## 🎓 Learning Resources

- **OpenTelemetry Docs**: https://opentelemetry.io/docs/languages/go/
- **Jaeger UI Guide**: https://www.jaegertracing.io/docs/latest/frontend-ui/
- **Chi Integration**: https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/go-chi/chi/v5/otelchi

---

## ⚠️ Troubleshooting

### Certificate Error (Current Issue)
```
Error: x509: OSStatus -26276
```

**Solutions:**
1. Check system time is correct
2. Update macOS: `softwareupdate -i -a`
3. Use `GOSUMDB=off` (dev only)
4. Use `go mod vendor`

### Jaeger Not Starting
```bash
# Check logs
docker-compose logs jaeger

# Restart
docker-compose restart jaeger
```

### No Traces Appearing
1. Check `OTEL_ENABLED=true` is set
2. Verify Jaeger is running: `docker-compose ps jaeger`
3. Check backend logs for "OpenTelemetry initialized"
4. Verify endpoint: `OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4318`

---

## 🎉 Summary

Once enabled, OpenTelemetry will provide:
- ✅ **Complete visibility** into request flows
- ✅ **Performance insights** for every endpoint
- ✅ **Easy debugging** of production issues
- ✅ **Bottleneck identification** in seconds
- ✅ **External API monitoring** (Stripe, etc.)

**Status:** Code is ready, just needs certificate issue resolved!

---

**Last Updated**: December 25, 2025
**Version**: 1.0.0
**Status**: ⏸️ Temporarily Disabled (Ready to Enable)

