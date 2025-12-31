# UsualStore - Design Patterns Used

This document catalogs all design patterns implemented in the UsualStore project with references to actual code.

## Table of Contents

1. [Creational Patterns](#creational-patterns)
2. [Structural Patterns](#structural-patterns)
3. [Behavioral Patterns](#behavioral-patterns)
4. [Architectural Patterns](#architectural-patterns)
5. [Concurrency Patterns](#concurrency-patterns)
6. [Testing Patterns](#testing-patterns)

---

## Creational Patterns

### 1. Factory Pattern

**Purpose**: Create objects without specifying the exact class

#### Implementation: New* Constructor Functions

**Location**: Throughout codebase

**Examples**:

```go
// AI Service Factory
// File: internal/ai/service.go:22-28
func NewService(db *sql.DB, aiClient AIClient, logger *log.Logger) *Service {
    return &Service{
        db:       db,
        aiClient: aiClient,
        logger:   logger,
    }
}
```

**More Examples**:
- `internal/ai/openai_client.go:30-47` - `NewOpenAIClient()`
- `internal/cards/cards.go` - `New()` for card processing
- `internal/driver/driver.go` - Database connection factory
- `internal/encryption/encryption.go` - `New()` for encryption service
- `internal/urlsigner/signer.go` - `New()` for URL signing
- `internal/validator/validator.go` - `New()` for validation

**Benefits**:
- Encapsulates complex initialization
- Dependency injection ready
- Consistent object creation
- Easy to test with mocks

---

### 2. Builder Pattern (Implicit)

**Purpose**: Construct complex objects step by step

#### Implementation: Request/Response Builders

**Location**: `internal/ai/models.go:81-107`

```go
// ChatRequest - builds complex request
type ChatRequest struct {
    SessionID  string  `json:"session_id"`
    Message    string  `json:"message"`
    UserID     *int    `json:"user_id,omitempty"`
    UserAgent  *string `json:"user_agent,omitempty"`
    IPAddress  *string `json:"ip_address,omitempty"`
}

// ChatResponse - builds complex response
type ChatResponse struct {
    Message        string                `json:"message"`
    Products       []RecommendedProduct  `json:"products,omitempty"`
    Suggestions    []string              `json:"suggestions,omitempty"`
    TokensUsed     int                   `json:"tokens_used"`
    ResponseTimeMs int                   `json:"response_time_ms"`
}
```

**More Examples**:
- `internal/models/models.go` - Various model builders
- `internal/messaging/types.go` - Email message builder

---

### 3. Singleton Pattern

**Purpose**: Ensure a class has only one instance

#### Implementation: Database Connection Pool

**Location**: `internal/driver/driver.go`

```go
// Database driver maintains single connection pool
type DB struct {
    SQL *sql.DB
}

// OpenDB creates a single database connection pool
func OpenDB(dsn string) (*DB, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }
    // Connection pool is singleton per database
    return &DB{SQL: db}, nil
}
```

**More Examples**:
- WebSocket Hub: `internal/support/hub.go` - Single hub per server
- Kafka Producers: `internal/messaging/producer.go` - Singleton producer

---

## Structural Patterns

### 4. Adapter Pattern

**Purpose**: Allow incompatible interfaces to work together

#### Implementation: AIClient Interface

**Location**: `internal/ai/models.go:134-139`

```go
// AIClient interface adapts different AI providers
type AIClient interface {
    GenerateResponse(messages []Message, context string) (*ChatResponse, error)
    GetEmbedding(text string) ([]float64, error)
    SpeechToText(audioData []byte, language string) (string, error)
    TextToSpeech(text, voice string) ([]byte, error)
}
```

**Implementation**: `internal/ai/openai_client.go:20-47`

```go
// OpenAIClient adapts OpenAI API to our interface
type OpenAIClient struct {
    APIKey      string
    Model       string
    Temperature float64
    MaxTokens   int
    HTTPClient  *http.Client
}

// Implements AIClient interface methods
func (c *OpenAIClient) GenerateResponse(...) {...}
func (c *OpenAIClient) GetEmbedding(...) {...}
func (c *OpenAIClient) SpeechToText(...) {...}
func (c *OpenAIClient) TextToSpeech(...) {...}
```

**Benefits**:
- Easy to swap AI providers (OpenAI → Anthropic → Custom)
- Testing with mock implementations
- Decoupling from external APIs

---

### 5. Facade Pattern

**Purpose**: Provide a simplified interface to a complex subsystem

#### Implementation: Service Layer

**Location**: `internal/ai/service.go:14-28`

```go
// Service facade hides complexity of AI operations
type Service struct {
    db       *sql.DB
    aiClient AIClient
    logger   *log.Logger
}

// Simple interface hiding complex operations:
// - Database queries
// - AI API calls
// - Context building
// - Message history management
func (s *Service) HandleChat(req ChatRequest) (*ChatResponse, error) {
    // Internally orchestrates:
    // 1. Get/create conversation
    // 2. Load history
    // 3. Get product context
    // 4. Generate AI response
    // 5. Save messages
    // 6. Update stats
}
```

**More Examples**:
- `internal/cards/cards.go` - Payment processing facade
- `internal/support/client.go` - Support service facade
- `cmd/api/handlers-api.go` - HTTP handler facade

---

### 6. Decorator Pattern

**Purpose**: Add behavior to objects dynamically

#### Implementation: Middleware Chain

**Location**: `cmd/api/middleware.go`

```go
// Auth middleware decorates handlers
func (app *application) Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add authentication behavior
        // Then call next handler
        next.ServeHTTP(w, r)
    })
}

// Logging middleware decorates handlers
func (app *application) Log(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add logging behavior
        // Then call next handler
        next.ServeHTTP(w, r)
    })
}

// Chain decorators:
// router.Use(app.Log, app.Auth, app.CORS)
```

**More Examples**:
- Telemetry decorators: `internal/telemetry/`
- Retry logic in messaging: `internal/messaging/`

---

### 7. Proxy Pattern

**Purpose**: Control access to an object

#### Implementation: API Gateway (Kong)

**Location**: `terraform/modules/api_gateway/main.tf`

```hcl
# Kong acts as proxy for all services
resource "docker_container" "kong" {
  # Routes requests to appropriate backend services
  # - /api/* → backend API
  # - /admin/* → support frontend
  # - /support/* → support service
}
```

**Benefits**:
- Single entry point
- Rate limiting
- Authentication
- Load balancing
- SSL termination

---

## Behavioral Patterns

### 8. Strategy Pattern

**Purpose**: Define a family of algorithms and make them interchangeable

#### Implementation: Payment Strategies

**Location**: `internal/cards/cards.go`

```go
// Card interface allows different payment strategies
type Card struct {
    Secret   string
    Key      string
    Currency string
}

// Can use different payment providers by changing implementation
// - Stripe
// - PayPal
// - Custom processor
```

**More Examples**:
- AI model selection strategy: `internal/ai/openai_client.go`
- Different database strategies: `internal/driver/driver.go`

---

### 9. Observer Pattern

**Purpose**: Notify multiple objects about state changes

#### Implementation: WebSocket Hub

**Location**: `internal/support/hub.go`

```go
// Hub broadcasts messages to all connected clients
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

// Run observes and notifies
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
        case message := <-h.broadcast:
            // Notify all observers
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}
```

**More Examples**:
- Kafka message broadcasting: `internal/messaging/producer.go`
- Event logging: `internal/telemetry/telemetry.go`

---

### 10. Command Pattern

**Purpose**: Encapsulate requests as objects

#### Implementation: CLI Commands

**Location**: `cmd/cli/main.go`

```go
// Each command encapsulates an action
// - migrate: Run database migrations
// - seed: Seed database
// - user:create: Create user
// - cache:clear: Clear cache

// Commands can be queued, logged, undone
```

---

### 11. Template Method Pattern

**Purpose**: Define algorithm skeleton, let subclasses override steps

#### Implementation: HTTP Handlers

**Location**: `cmd/api/handlers-api.go`

```go
// Base pattern for all handlers:
func (app *application) handlerTemplate(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request (can be customized)
    // 2. Validate input (can be customized)
    // 3. Business logic (can be customized)
    // 4. Response (can be customized)
}
```

---

### 12. Chain of Responsibility Pattern

**Purpose**: Pass requests along a chain of handlers

#### Implementation: Middleware Chain

**Location**: `cmd/api/routes-api.go`

```go
// Request passes through middleware chain:
// Request → Logger → Auth → CORS → RateLimit → Handler
mux := chi.NewRouter()
mux.Use(middleware.Logger)    // First handler
mux.Use(app.Auth)              // Second handler
mux.Use(app.CORS)              // Third handler
mux.Use(middleware.RateLimit(100)) // Fourth handler
// Finally reaches actual handler
```

---

## Architectural Patterns

### 13. Repository Pattern

**Purpose**: Abstract data access logic

#### Implementation: Database Interface

**Location**: `internal/ai/models.go:115-132`

```go
// DB interface abstracts database operations
type DB interface {
    GetConversation(sessionID string) (*Conversation, error)
    CreateConversation(conv *Conversation) error
    GetMessages(conversationID int, limit int) ([]Message, error)
    CreateMessage(msg *Message) error
    UpdateConversation(conv *Conversation) error
    GetUserPreferences(userID *int, sessionID *string) (*UserPreferences, error)
    SaveFeedback(feedback *Feedback) error
    GetProductContext() (string, error)
    UpdateProductPopularity(productID int) error
}
```

**Implementation**: `internal/ai/service.go` - Service implements repository methods

**Benefits**:
- Testable with mock databases
- Easy to switch database engines
- Centralized data access
- Clean separation of concerns

---

### 14. Service Layer Pattern

**Purpose**: Define application boundary and business logic

#### Implementation: AI Service

**Location**: `internal/ai/service.go:14-28`

```go
// Service layer coordinates between presentation and data layers
type Service struct {
    db       *sql.DB        // Data layer
    aiClient AIClient       // External service
    logger   *log.Logger    // Infrastructure
}

// Business logic encapsulated in service methods
func (s *Service) HandleChat(req ChatRequest) (*ChatResponse, error)
func (s *Service) HandleVoiceInput(...) (string, *ChatResponse, error)
func (s *Service) SubmitFeedback(...) error
```

**More Examples**:
- Messaging service: `internal/messaging/`
- Support service: `internal/support/`
- Card processing service: `internal/cards/`

---

### 15. Domain Model Pattern

**Purpose**: Create object model of the domain

#### Implementation: Domain Models

**Location**: `internal/ai/models.go:8-79`

```go
// Domain entities represent business concepts
type Conversation struct {
    ID                 int
    SessionID          string
    UserID             *int
    StartedAt          time.Time
    EndedAt            *time.Time
    TotalMessages      int
    ResultedInPurchase bool
    TotalTokensUsed    int
    TotalCost          float64
}

type Message struct {
    ID             int
    ConversationID int
    Role           string
    Content        string
    TokensUsed     int
    ResponseTimeMs *int
}

type UserPreferences struct {
    ID                    int
    PreferredCategories   []string
    BudgetMin             *float64
    BudgetMax             *float64
    InteractionCount      int
}
```

**More Examples**:
- `internal/models/models.go` - Core domain models
- `internal/models/customer_model.go` - Customer entity
- `internal/models/widget_model.go` - Product entity

---

### 16. Dependency Injection Pattern

**Purpose**: Inject dependencies rather than creating them

#### Implementation: Throughout Application

**Location**: `cmd/api/main.go` (example)

```go
// Dependencies injected via constructors
func main() {
    // Create dependencies
    db := driver.OpenDB(dsn)
    logger := log.New(os.Stdout, "", log.LstdFlags)
    aiClient := ai.NewOpenAIClient(apiKey, model, temperature)
    
    // Inject into service
    aiService := ai.NewService(db.SQL, aiClient, logger)
    
    // Inject into handler
    app := &application{
        config:    cfg,
        infoLog:   logger,
        errorLog:  errorLogger,
        aiService: aiService,
        db:        db,
    }
}
```

**Benefits**:
- Loose coupling
- Easy testing
- Flexible configuration
- Runtime dependency swapping

---

### 17. MVC Pattern (Modified)

**Purpose**: Separate concerns in web application

#### Implementation: HTTP Layer Structure

**Structure**:
```
cmd/api/
├── main.go           # Application entry (Controller setup)
├── handlers-api.go   # Controllers
├── routes-api.go     # Routing
└── middleware.go     # Middleware

internal/models/      # Models (Domain)
internal/ai/         # Business Logic

react-frontend/      # Views (React components)
```

**Example**:
- Model: `internal/models/models.go`
- Controller: `cmd/api/handlers-api.go`
- View: `react-frontend/src/pages/`

---

### 18. Microservices Pattern

**Purpose**: Build application as suite of small services

#### Implementation: Docker Services Architecture

**Services**:
```
┌─────────────────────────────────────────┐
│  API Gateway (Kong) - Port 8000         │
└─────────────────────────────────────────┘
              │
    ┌─────────┼─────────┬──────────┐
    │         │         │          │
    ▼         ▼         ▼          ▼
┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
│Backend │ │   AI   │ │Support │ │Messaging│
│  API   │ │Assistant│ │Service │ │ Service│
│  4001  │ │  8080  │ │  5001  │ │  6001  │
└────────┘ └────────┘ └────────┘ └────────┘
```

**Location**: `terraform/main.tf`, `docker-compose.yml`

**Services**:
- `usualstore-back-end` - Main API
- `usualstore-ai-assistant` - AI service
- `usualstore-support-service` - Support WebSocket
- `usualstore-messaging-service` - Email/notifications
- `usualstore-database` - PostgreSQL
- `api-gateway` - Kong gateway

---

## Concurrency Patterns

### 19. Worker Pool Pattern

**Purpose**: Limit concurrent operations

#### Implementation: Database Connection Pool

**Location**: `internal/driver/driver.go`

```go
// Connection pool manages concurrent database access
db.SetMaxOpenConns(25)        // Maximum workers
db.SetMaxIdleConns(25)        // Idle workers
db.SetConnMaxLifetime(5 * time.Minute)
```

---

### 20. Fan-Out/Fan-In Pattern

**Purpose**: Distribute work and collect results

#### Implementation: Kafka Message Processing

**Location**: `internal/messaging/consumer.go`

```go
// Consumer fans out messages to handlers
type EmailHandler interface {
    HandleEmail(msg EmailMessage) error
}

// Multiple consumers process messages concurrently
// Results fan in for aggregation
```

---

### 21. Pipeline Pattern

**Purpose**: Process data through series of stages

#### Implementation: Request Processing Pipeline

**Flow**:
```
Request → Middleware → Handler → Service → Repository → Database
         (Logging)    (Validate)  (Logic)   (Data)     (Store)
              ↓           ↓          ↓         ↓          ↓
         Response ← JSON  ← Result  ← Data   ← Query   ← DB
```

---

## Testing Patterns

### 22. Mock Object Pattern

**Purpose**: Replace real objects with test doubles

#### Implementation: Mock Repository

**Location**: `internal/mocks/mock_token_repository.go`

```go
// Generated mock using mockgen
type MockTokenRepository struct {
    ctrl     *gomock.Controller
    recorder *MockTokenRepositoryMockRecorder
}

// Mock methods for testing
func (m *MockTokenRepository) EXPECT() *MockTokenRepositoryMockRecorder {
    return m.recorder
}
```

**More Examples**:
- `internal/ai/service_db_test.go:890-927` - MockAIClient

```go
// MockAIClient for testing
type MockAIClient struct {
    GenerateResponseFunc func(messages []Message, context string) (*ChatResponse, error)
    GetEmbeddingFunc     func(text string) ([]float64, error)
    SpeechToTextFunc     func(audioData []byte, language string) (string, error)
    TextToSpeechFunc     func(text, voice string) ([]byte, error)
}

func (m *MockAIClient) GenerateResponse(messages []Message, context string) (*ChatResponse, error) {
    if m.GenerateResponseFunc != nil {
        return m.GenerateResponseFunc(messages, context)
    }
    return &ChatResponse{Message: "Test response"}, nil
}
```

---

### 23. Table-Driven Tests Pattern

**Purpose**: Test multiple scenarios with same logic

#### Implementation: AI Service Tests

**Location**: `internal/ai/service_handlechat_test.go:13-375`

```go
func TestHandleChat(t *testing.T) {
    tests := []struct {
        name        string
        request     ChatRequest
        setupMock   func(sqlmock.Sqlmock)
        aiClient    *MockAIClient
        expectError bool
        checkResult func(*testing.T, *ChatResponse)
    }{
        {
            name: "successful chat with existing conversation",
            request: ChatRequest{...},
            setupMock: func(mock sqlmock.Sqlmock) {...},
            aiClient: &MockAIClient{...},
            expectError: false,
        },
        {
            name: "new conversation created",
            request: ChatRequest{...},
            // ... more test cases
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

**Benefits**:
- Comprehensive test coverage
- Easy to add new test cases
- Clear test documentation
- Reduced code duplication

---

## Infrastructure Patterns

### 24. Infrastructure as Code (IaC)

**Purpose**: Manage infrastructure through code

#### Implementation: Terraform

**Location**: `terraform/`

```hcl
# Declarative infrastructure definition
resource "docker_container" "backend" {
  name  = "usualstore-back-end"
  image = docker_image.backend.image_id
  
  ports {
    internal = 4001
    external = 4001
  }
  
  networks_advanced {
    name = var.network_id
  }
}
```

---

### 25. Circuit Breaker Pattern

**Purpose**: Prevent cascade failures

#### Implementation: HTTP Client Timeouts

**Location**: `internal/ai/openai_client.go:43-45`

```go
HTTPClient: &http.Client{
    Timeout: 30 * time.Second,
},
```

---

### 26. Retry Pattern

**Purpose**: Retry failed operations

#### Implementation: Database Retries

**Location**: `internal/driver/driver.go`

```go
// Retry connection with backoff
for i := 0; i < maxRetries; i++ {
    if err := db.Ping(); err == nil {
        return nil
    }
    time.Sleep(time.Second * time.Duration(i+1))
}
```

---

## Summary

### Pattern Distribution

| Category | Count | Patterns |
|----------|-------|----------|
| **Creational** | 3 | Factory, Builder, Singleton |
| **Structural** | 4 | Adapter, Facade, Decorator, Proxy |
| **Behavioral** | 5 | Strategy, Observer, Command, Template, Chain of Responsibility |
| **Architectural** | 6 | Repository, Service Layer, Domain Model, DI, MVC, Microservices |
| **Concurrency** | 3 | Worker Pool, Fan-Out/Fan-In, Pipeline |
| **Testing** | 2 | Mock Objects, Table-Driven Tests |
| **Infrastructure** | 3 | IaC, Circuit Breaker, Retry |

**Total**: 26 Design Patterns

---

## Design Principles Applied

### SOLID Principles

1. **Single Responsibility**: Each service has one responsibility
   - `internal/ai/service.go` - AI operations only
   - `internal/cards/cards.go` - Payment processing only

2. **Open/Closed**: Extend through interfaces, don't modify
   - `AIClient` interface allows new AI providers without changing Service

3. **Liskov Substitution**: Interfaces can be replaced with implementations
   - `MockAIClient` can replace `OpenAIClient` in tests

4. **Interface Segregation**: Small, focused interfaces
   - `AIClient` interface has only AI-specific methods
   - `DB` interface has only data methods

5. **Dependency Inversion**: Depend on abstractions
   - Service depends on `AIClient` interface, not concrete `OpenAIClient`

---

## Best Practices Observed

1. **Accept interfaces, return concrete types**
   ```go
   func NewService(db *sql.DB, aiClient AIClient, logger *log.Logger) *Service
   ```

2. **Constructor injection for dependencies**
   ```go
   NewService(db, aiClient, logger)
   ```

3. **Small, focused interfaces (1-4 methods)**
   ```go
   type AIClient interface {
       GenerateResponse(...) (...)
       GetEmbedding(...) (...)
       SpeechToText(...) (...)
       TextToSpeech(...) (...)
   }
   ```

4. **Explicit error handling**
   ```go
   if err != nil {
       return fmt.Errorf("context: %w", err)
   }
   ```

5. **Table-driven tests**
   ```go
   tests := []struct{ name string; ... }{ ... }
   ```

---

## Anti-Patterns Avoided

✅ **God Object** - Avoided by separating concerns into services
✅ **Spaghetti Code** - Clear layered architecture
✅ **Tight Coupling** - Interfaces enable loose coupling
✅ **Premature Optimization** - Clear, readable code first
✅ **Not Invented Here** - Uses standard libraries and patterns

---

## References

- **Go Design Patterns**: https://github.com/tmrts/go-patterns
- **SOLID Principles**: https://dave.cheney.net/2016/08/20/solid-go-design
- **Effective Go**: https://go.dev/doc/effective_go
- **Gang of Four**: Design Patterns book
- **Clean Architecture**: Robert C. Martin

---

**Document Version**: 1.0
**Last Updated**: December 31, 2025
**Project**: UsualStore

