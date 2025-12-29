# 🍃 MongoDB Multi-Tenant Guide

**Use MongoDB instead of PostgreSQL - locally and in AWS**

---

## 🎯 Yes, You Can Use MongoDB!

Your multi-tenant architecture works perfectly with MongoDB:

- ✅ **Locally**: MongoDB in Docker
- ✅ **AWS**: DocumentDB, MongoDB Atlas, or EC2
- ✅ **Same Pattern**: Customer chooses database name
- ✅ **Same Isolation**: Each customer = separate database
- ✅ **Same Code**: Minimal changes to your Go application

---

## 🏗️ Architecture Comparison

### Current (PostgreSQL)

```
Customer 1 → building_shop (PostgreSQL DB) → Go App
Customer 2 → hardware_store (PostgreSQL DB) → Go App
Customer 3 → bookstore_2025 (PostgreSQL DB) → Go App
```

### With MongoDB

```
Customer 1 → building_shop (MongoDB Database) → Go App
Customer 2 → hardware_store (MongoDB Database) → Go App
Customer 3 → bookstore_2025 (MongoDB Database) → Go App

Same pattern, different database engine!
```

---

## 🔄 MongoDB vs PostgreSQL

### What Changes

| Aspect | PostgreSQL | MongoDB |
|--------|-----------|---------|
| **Data Model** | Tables, rows, columns | Collections, documents, fields |
| **Query Language** | SQL | MongoDB Query Language |
| **Schema** | Fixed schema | Flexible schema |
| **Connections** | `psql`, connection string | `mongosh`, connection string |
| **Local Setup** | PostgreSQL container | MongoDB container |
| **AWS Options** | RDS PostgreSQL | DocumentDB, Atlas, EC2 |

### What Stays the Same

- ✅ Customer chooses database name
- ✅ Terraform creates databases
- ✅ Complete isolation per customer
- ✅ Same application code for all customers
- ✅ Environment variables configure which database

---

## 🐳 Local MongoDB Setup

### Docker Compose with MongoDB

```yaml
# docker-compose-mongodb.yml
version: '3.8'

services:
  mongodb:
    image: mongo:7.0
    container_name: usualstore-mongodb
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password123
    volumes:
      - mongodb_data:/data/db
    networks:
      - usualstore_network

  backend:
    build: .
    environment:
      MONGODB_URI: "mongodb://admin:password123@mongodb:27017"
      DATABASE_NAME: "building_shop"  # Customer-specific!
      TENANT_NAME: "building_shop"
    ports:
      - "4001:4001"
    depends_on:
      - mongodb
    networks:
      - usualstore_network

volumes:
  mongodb_data:

networks:
  usualstore_network:
    driver: bridge
```

### Start MongoDB Locally

```bash
# Start MongoDB
docker-compose -f docker-compose-mongodb.yml up -d

# Connect to MongoDB
docker exec -it usualstore-mongodb mongosh -u admin -p password123

# Create customer database
use building_shop
db.createCollection("customers")
db.createCollection("orders")
db.createCollection("widgets")

# Verify
show dbs
show collections
```

---

## 🔧 Go Application Changes

### 1. MongoDB Driver

```bash
# Add MongoDB driver
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options
```

### 2. Database Connection

```go
// internal/driver/mongodb.go (new file)
package driver

import (
    "context"
    "fmt"
    "time"
    
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
    Client   *mongo.Client
    Database *mongo.Database
}

func ConnectMongoDB(uri, dbName string) (*MongoDB, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Connect to MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }
    
    // Ping to verify connection
    if err := client.Ping(ctx, nil); err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }
    
    // Select database (customer-specific!)
    database := client.Database(dbName)
    
    return &MongoDB{
        Client:   client,
        Database: database,
    }, nil
}

func (m *MongoDB) Close() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    return m.Client.Disconnect(ctx)
}
```

### 3. Environment Variables

```bash
# Instead of DATABASE_DSN (PostgreSQL)
export MONGODB_URI="mongodb://admin:password123@localhost:27017"
export DATABASE_NAME="building_shop"  # Customer chooses this!
export TENANT_NAME="building_shop"

# Or combined
export MONGODB_URI="mongodb://admin:password123@localhost:27017/building_shop?authSource=admin"
```

### 4. Models Update

```go
// internal/models/models.go (MongoDB version)
package models

import (
    "context"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type Customer struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    FirstName string             `bson:"first_name" json:"first_name"`
    LastName  string             `bson:"last_name" json:"last_name"`
    Email     string             `bson:"email" json:"email"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Order struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
    Status     string             `bson:"status" json:"status"`
    Total      float64            `bson:"total" json:"total"`
    CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// MongoDB operations
type DBModel struct {
    DB *mongo.Database
}

func (m *DBModel) GetAllCustomers() ([]Customer, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    collection := m.DB.Collection("customers")
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var customers []Customer
    if err := cursor.All(ctx, &customers); err != nil {
        return nil, err
    }
    
    return customers, nil
}

func (m *DBModel) InsertCustomer(customer *Customer) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    customer.CreatedAt = time.Now()
    customer.UpdatedAt = time.Now()
    
    collection := m.DB.Collection("customers")
    result, err := collection.InsertOne(ctx, customer)
    if err != nil {
        return err
    }
    
    customer.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}
```

### 5. Main Application

```go
// cmd/api/main.go (MongoDB version)
package main

import (
    "log"
    "os"
    
    "your-app/internal/driver"
    "your-app/internal/models"
)

func main() {
    // Get environment variables
    mongoURI := os.Getenv("MONGODB_URI")
    dbName := os.Getenv("DATABASE_NAME")  // Customer-specific!
    
    if mongoURI == "" || dbName == "" {
        log.Fatal("MONGODB_URI and DATABASE_NAME must be set")
    }
    
    // Connect to MongoDB
    mongodb, err := driver.ConnectMongoDB(mongoURI, dbName)
    if err != nil {
        log.Fatalf("Cannot connect to MongoDB: %v", err)
    }
    defer mongodb.Close()
    
    log.Printf("Connected to MongoDB database: %s", dbName)
    
    // Initialize models
    app := &config{
        Models: models.DBModel{DB: mongodb.Database},
    }
    
    // Start server
    app.serve()
}
```

---

## 🏢 Terraform for Local MongoDB

### MongoDB Container per Customer

```hcl
# terraform-tenants/mongodb/main.tf

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

provider "docker" {
  host = "unix:///var/run/docker.sock"
}

# Shared MongoDB server
resource "docker_container" "mongodb" {
  name  = "usualstore-mongodb"
  image = "mongo:7.0"
  
  ports {
    internal = 27017
    external = 27017
  }
  
  env = [
    "MONGO_INITDB_ROOT_USERNAME=admin",
    "MONGO_INITDB_ROOT_PASSWORD=password123"
  ]
  
  volumes {
    volume_name    = "mongodb_data"
    container_path = "/data/db"
  }
}

# Create customer databases using null_resource
resource "null_resource" "create_customer_databases" {
  for_each = var.tenants
  
  depends_on = [docker_container.mongodb]
  
  provisioner "local-exec" {
    command = <<-EOT
      docker exec usualstore-mongodb mongosh -u admin -p password123 --eval "
        use ${each.value.database_name}
        db.createCollection('customers')
        db.createCollection('orders')
        db.createCollection('widgets')
        db.createCollection('transactions')
        db.createCollection('statuses')
        
        // Create admin user for this database
        db.createUser({
          user: '${each.value.admins[0].username}',
          pwd: '${each.value.admins[0].password}',
          roles: [
            { role: 'readWrite', db: '${each.value.database_name}' },
            { role: 'dbAdmin', db: '${each.value.database_name}' }
          ]
        })
      "
    EOT
  }
}

# Deploy customer app
resource "docker_container" "customer_app" {
  for_each = var.tenants
  
  name  = "${each.key}-app"
  image = "usualstore/api:latest"
  
  env = [
    "MONGODB_URI=mongodb://${each.value.admins[0].username}:${each.value.admins[0].password}@mongodb:27017/${each.value.database_name}?authSource=${each.value.database_name}",
    "DATABASE_NAME=${each.value.database_name}",
    "TENANT_NAME=${each.key}",
    "API_PORT=4001"
  ]
  
  ports {
    internal = 4001
    external = 4001 + index(keys(var.tenants), each.key)
  }
  
  depends_on = [
    docker_container.mongodb,
    null_resource.create_customer_databases[each.key]
  ]
}
```

---

## ☁️ AWS MongoDB Options

### Option 1: Amazon DocumentDB (Recommended)

**What is DocumentDB?**
- AWS-managed MongoDB-compatible database
- Fully managed (backups, scaling, patches)
- MongoDB 4.0 compatible
- Integrated with AWS ecosystem

**Terraform Configuration:**

```hcl
# terraform-tenants/aws-mongodb/documentdb.tf

# DocumentDB cluster (shared or per-customer)
resource "aws_docdb_cluster" "building_shop" {
  cluster_identifier      = "building-shop-docdb"
  engine                  = "docdb"
  master_username         = "admin"
  master_password         = var.master_password
  backup_retention_period = 7
  preferred_backup_window = "03:00-04:00"
  skip_final_snapshot     = false
  
  vpc_security_group_ids = [aws_security_group.docdb.id]
  db_subnet_group_name   = aws_docdb_subnet_group.main.name
  
  tags = {
    Customer = "Building Shop"
  }
}

# DocumentDB instance
resource "aws_docdb_cluster_instance" "building_shop" {
  identifier         = "building-shop-docdb-instance"
  cluster_identifier = aws_docdb_cluster.building_shop.id
  instance_class     = "db.t3.medium"
  
  tags = {
    Customer = "Building Shop"
  }
}

# Security group
resource "aws_security_group" "docdb" {
  name        = "documentdb-sg"
  vpc_id      = aws_vpc.main.id
  description = "Security group for DocumentDB"
  
  ingress {
    from_port       = 27017
    to_port         = 27017
    protocol        = "tcp"
    security_groups = [aws_security_group.lambda.id]
    description     = "MongoDB from Lambda"
  }
}

# Lambda with DocumentDB
resource "aws_lambda_function" "building_shop_api" {
  function_name = "building-shop-api"
  # ... other config ...
  
  environment {
    variables = {
      MONGODB_URI   = "mongodb://${aws_docdb_cluster.building_shop.master_username}:${var.master_password}@${aws_docdb_cluster.building_shop.endpoint}:27017"
      DATABASE_NAME = "building_shop"  # Customer database!
      TENANT_NAME   = "building_shop"
    }
  }
  
  vpc_config {
    subnet_ids         = aws_subnet.private[*].id
    security_group_ids = [aws_security_group.lambda.id]
  }
}
```

**Cost (DocumentDB):**
- db.t3.medium instance: ~$70-90/month
- Storage: ~$0.10/GB/month
- Backups: ~$0.02/GB/month
- **Total per customer: ~$80-100/month** (separate cluster)
- **Shared cluster: ~$80-100/month for multiple customers**

### Option 2: MongoDB Atlas (Cloud SaaS)

**What is Atlas?**
- Official MongoDB cloud service
- Fully managed
- Multi-cloud (AWS, GCP, Azure)
- Free tier available!

**Setup:**

1. Create Atlas account: https://www.mongodb.com/cloud/atlas
2. Create cluster in AWS region
3. Create database per customer
4. Get connection string

**Terraform with Atlas:**

```hcl
# terraform-tenants/aws-mongodb/atlas.tf

terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.14"
    }
  }
}

provider "mongodbatlas" {
  public_key  = var.atlas_public_key
  private_key = var.atlas_private_key
}

# MongoDB Atlas Project
resource "mongodbatlas_project" "usualstore" {
  name   = "usualstore-tenants"
  org_id = var.atlas_org_id
}

# Shared cluster (M10 for multiple tenants)
resource "mongodbatlas_cluster" "shared" {
  project_id = mongodbatlas_project.usualstore.id
  name       = "usualstore-shared"
  
  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"  # Shared by customers
  
  mongo_db_major_version = "7.0"
}

# Database per customer (created via Atlas UI or API)
resource "mongodbatlas_database_user" "building_shop" {
  project_id = mongodbatlas_project.usualstore.id
  auth_database_name = "admin"
  
  username = "building_admin"
  password = var.building_shop_password
  
  roles {
    role_name     = "readWrite"
    database_name = "building_shop"  # Customer database!
  }
}

# Lambda connects to Atlas
resource "aws_lambda_function" "building_shop_api" {
  function_name = "building-shop-api"
  
  environment {
    variables = {
      MONGODB_URI   = "mongodb+srv://building_admin:${var.building_shop_password}@usualstore-shared.mongodb.net"
      DATABASE_NAME = "building_shop"
      TENANT_NAME   = "building_shop"
    }
  }
}
```

**Cost (MongoDB Atlas):**
- M10 cluster (shared): ~$57/month
- Free tier (M0): $0/month (for testing)
- M10 can handle multiple tenant databases
- **Total: ~$57/month for multiple customers**

### Option 3: Self-Managed on EC2

**When to use:**
- Full control needed
- Cost optimization for large scale
- Custom MongoDB configuration

```hcl
# EC2 with MongoDB
resource "aws_instance" "mongodb" {
  ami           = "ami-0c55b159cbfafe1f0"  # Ubuntu
  instance_class = "t3.medium"
  
  user_data = <<-EOF
    #!/bin/bash
    # Install MongoDB
    wget -qO - https://www.mongodb.org/static/pgp/server-7.0.asc | apt-key add -
    echo "deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/7.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-7.0.list
    apt-get update
    apt-get install -y mongodb-org
    systemctl start mongod
    systemctl enable mongod
  EOF
  
  vpc_security_group_ids = [aws_security_group.mongodb.id]
}
```

---

## 🔄 Migration: PostgreSQL to MongoDB

### Strategy 1: Dual Write (Safe Migration)

```go
// Write to both databases during migration
func (app *application) CreateOrder(order *Order) error {
    // Write to PostgreSQL
    if err := app.PostgreSQL.CreateOrder(order); err != nil {
        return err
    }
    
    // Write to MongoDB (new)
    if err := app.MongoDB.CreateOrder(order); err != nil {
        log.Printf("MongoDB write failed: %v", err)
        // Don't fail the request
    }
    
    return nil
}

// Read from PostgreSQL (old), fallback to MongoDB
func (app *application) GetOrder(id string) (*Order, error) {
    order, err := app.PostgreSQL.GetOrder(id)
    if err == nil {
        return order, nil
    }
    
    // Fallback to MongoDB
    return app.MongoDB.GetOrder(id)
}
```

### Strategy 2: Data Export/Import

```bash
# Export from PostgreSQL
pg_dump -U building_admin -d building_shop -t customers --data-only --column-inserts > customers.sql

# Convert SQL to MongoDB (using script)
python sql_to_mongodb.py customers.sql

# Import to MongoDB
mongoimport --db building_shop --collection customers --file customers.json --jsonArray
```

### Strategy 3: New Customers Start with MongoDB

- Existing customers: Keep PostgreSQL
- New customers: Use MongoDB
- Migrate customers gradually

---

## 📊 Comparison Matrix

### Local Development

| Feature | PostgreSQL | MongoDB |
|---------|-----------|---------|
| **Setup** | `docker run postgres` | `docker run mongo` |
| **Admin UI** | pgAdmin, DBeaver | MongoDB Compass |
| **Connection** | `psql` | `mongosh` |
| **Complexity** | Medium | Low |

### AWS Deployment

| Feature | RDS PostgreSQL | DocumentDB | MongoDB Atlas |
|---------|---------------|------------|---------------|
| **Management** | AWS managed | AWS managed | MongoDB managed |
| **Cost (small)** | ~$15-25/mo | ~$70-90/mo | ~$57/mo (shared) |
| **Compatibility** | Full PostgreSQL | MongoDB 4.0 | Full MongoDB |
| **Multi-region** | AWS regions | AWS regions | Multi-cloud |
| **Backup** | Automated | Automated | Automated |
| **Scaling** | Vertical/horizontal | Horizontal | Automatic |

---

## 🎯 Recommendations

### For Development/Testing
```
✅ Use MongoDB locally in Docker
✅ Same multi-tenant pattern
✅ Fast setup
✅ Easy to test
```

### For Production (Small Scale)
```
✅ MongoDB Atlas M10 (shared cluster)
✅ ~$57/month for multiple customers
✅ Fully managed
✅ Free tier for testing
```

### For Production (Large Scale)
```
✅ Amazon DocumentDB (separate clusters)
✅ Full AWS integration
✅ VPC isolation
✅ Custom scaling
```

### For Mixed Approach
```
✅ Keep PostgreSQL for existing customers
✅ Use MongoDB for new customers
✅ Migrate gradually
✅ Both work with same architecture!
```

---

## 📁 File Structure with MongoDB

```
usual_store/
├── terraform-tenants/
│   ├── mongodb/              ← Local MongoDB Terraform
│   │   ├── main.tf
│   │   └── variables.tf
│   │
│   ├── aws-mongodb/          ← AWS MongoDB Terraform
│   │   ├── documentdb.tf     ← DocumentDB option
│   │   ├── atlas.tf          ← Atlas option
│   │   └── lambda.tf
│   │
│   ├── MONGODB-GUIDE.md      ← This guide
│   └── add-customer-mongodb.sh  ← Script for MongoDB
│
├── internal/
│   ├── driver/
│   │   ├── postgres.go       ← PostgreSQL driver
│   │   └── mongodb.go        ← MongoDB driver (new)
│   │
│   └── models/
│       ├── models_postgres.go
│       └── models_mongodb.go
│
└── cmd/api/
    ├── main.go               ← Updated for both DBs
    └── config.go
```

---

## ✅ Quick Start Checklist

### Local MongoDB Setup

- [ ] Install Docker
- [ ] Create `docker-compose-mongodb.yml`
- [ ] Start MongoDB container
- [ ] Create customer databases
- [ ] Update Go code for MongoDB
- [ ] Test locally

### AWS DocumentDB Setup

- [ ] Create AWS account
- [ ] Set up VPC and subnets
- [ ] Create DocumentDB cluster
- [ ] Create customer database
- [ ] Update Lambda for MongoDB
- [ ] Test connection

### MongoDB Atlas Setup

- [ ] Create Atlas account
- [ ] Create cluster in AWS region
- [ ] Create customer database
- [ ] Get connection string
- [ ] Update Lambda
- [ ] Test connection

---

## 🎯 Summary

### Can You Use MongoDB?

**YES!** ✅

- ✅ Same multi-tenant architecture
- ✅ Customer chooses database name
- ✅ Terraform creates databases
- ✅ Works locally (Docker)
- ✅ Works in AWS (DocumentDB or Atlas)
- ✅ Complete isolation per customer
- ✅ Minimal code changes

### Three AWS Options

1. **DocumentDB** - AWS managed, MongoDB compatible, ~$70-90/month
2. **MongoDB Atlas** - Official MongoDB SaaS, ~$57/month shared
3. **EC2** - Self-managed, full control, variable cost

### Best Choice

For most use cases: **MongoDB Atlas**
- Fully managed
- Free tier available
- Works with Lambda
- ~$57/month for shared cluster serving multiple tenants

---

**Next Steps:**
1. Try MongoDB locally with Docker
2. Update your Go application
3. Test with one customer
4. Choose AWS option (Atlas recommended)
5. Deploy to production

Your multi-tenant architecture is database-agnostic! 🎉

