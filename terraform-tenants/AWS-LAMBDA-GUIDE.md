# 🚀 AWS Lambda Deployment Guide

**How to deploy your multi-tenant architecture to AWS Lambda + RDS**

---

## 🎯 Architecture Overview

### Current (Local/Docker)

```
Customer 1 → building_shop DB (PostgreSQL) → App Container (Go)
Customer 2 → hardware_store DB (PostgreSQL) → App Container (Go)
```

### With AWS Lambda

```
Customer 1 → building_shop DB (RDS) → Lambda Function (Go)
Customer 2 → hardware_store DB (RDS) → Lambda Function (Go)
                                        ↑
                                 API Gateway
```

---

## 🏗️ AWS Components

### Per Customer Deployment

Each customer gets:

1. **RDS PostgreSQL Instance** (or database on shared RDS)
   - Database: `building_shop`
   - Isolated from other customers
   - Automatic backups
   - Multi-AZ for high availability

2. **Lambda Function**
   - Go runtime (using provided.al2 or provided.al2023)
   - Environment variables point to customer's RDS
   - Same code for ALL customers!
   - Auto-scaling

3. **API Gateway**
   - Custom domain per customer
   - HTTPS endpoints
   - Routes to Lambda

4. **S3 Buckets** (for frontend)
   - Static React/TypeScript build
   - CloudFront CDN

---

## 📋 Two Deployment Patterns

### Pattern 1: Separate RDS Instance Per Customer (Isolated)

**Best for:**
- Large customers
- Compliance requirements (HIPAA, SOC2)
- Performance isolation
- Different AWS regions per customer

```
Customer 1: Building Shop
├── RDS: building-shop.xyz.rds.amazonaws.com
├── Lambda: building-shop-api
├── API Gateway: api.buildingshop.com
└── S3 + CloudFront: buildingshop.com

Customer 2: Hardware Store
├── RDS: hardware-store.xyz.rds.amazonaws.com
├── Lambda: hardware-store-api
├── API Gateway: api.hardwarestore.com
└── S3 + CloudFront: hardwarestore.com
```

**Cost:** Higher (RDS instances ~$50-200/month each)

### Pattern 2: Shared RDS, Multiple Databases (Cost-Effective)

**Best for:**
- Small to medium customers
- Cost optimization
- Centralized management

```
Shared RDS: main.xyz.rds.amazonaws.com
├── Database: building_shop
├── Database: hardware_store
└── Database: bookstore_2025

Customer 1 Lambda → Connects to building_shop DB
Customer 2 Lambda → Connects to hardware_store DB
Customer 3 Lambda → Connects to bookstore_2025 DB
```

**Cost:** Lower (One RDS ~$50-200/month, multiple DBs)

---

## 🔧 Lambda Function Structure

### Your Go API as Lambda

Your existing Go application needs minimal changes:

```go
// cmd/lambda/main.go (new file)
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    
    // Your existing imports
    "your-app/internal/driver"
    "your-app/internal/models"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Your existing application logic
    // Environment variable DATABASE_DSN tells which customer DB to use
    
    db, err := driver.ConnectSQL(os.Getenv("DATABASE_DSN"))
    if err != nil {
        return events.APIGatewayProxyResponse{
            StatusCode: 500,
            Body:       "Database connection failed",
        }, nil
    }
    
    // Route request to appropriate handler
    // Your existing handler code here
    
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       "response",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }, nil
}

func main() {
    lambda.Start(handler)
}
```

**Key Point:** Same application logic, different entry point!

---

## 🌍 Terraform for AWS Lambda Deployment

### Directory Structure

```
terraform-tenants/
├── aws/                          ← New directory for AWS
│   ├── main.tf                   ← AWS Lambda + RDS setup
│   ├── variables.tf
│   ├── lambda.tf                 ← Lambda configuration
│   ├── rds.tf                    ← RDS configuration
│   ├── api_gateway.tf            ← API Gateway
│   ├── s3_frontend.tf            ← Frontend hosting
│   │
│   └── customers/                ← Per-customer configs
│       ├── building_shop.tfvars
│       ├── hardware_store.tfvars
│       └── bookstore.tfvars
│
└── add-customer-aws.sh           ← Script for AWS customers
```

### Terraform Configuration Example

#### Lambda Function (per customer)

```hcl
# terraform-tenants/aws/lambda.tf

# Building Shop Lambda Function
resource "aws_lambda_function" "building_shop_api" {
  filename      = "../../dist/lambda.zip"  # Your compiled Go binary
  function_name = "building-shop-api"
  role          = aws_iam_role.lambda_exec.arn
  handler       = "bootstrap"              # For Go on Lambda
  runtime       = "provided.al2023"
  
  environment {
    variables = {
      DATABASE_DSN = "host=${aws_db_instance.building_shop.endpoint} port=5432 user=building_admin password=${var.db_password} dbname=building_shop sslmode=require"
      TENANT_NAME  = "building_shop"
      API_PORT     = "4001"
    }
  }
  
  vpc_config {
    subnet_ids         = aws_subnet.private[*].id
    security_group_ids = [aws_security_group.lambda.id]
  }
  
  memory_size = 512
  timeout     = 30
  
  tags = {
    Customer = "Building Shop"
  }
}
```

#### RDS Database (per customer or shared)

**Option 1: Separate RDS per customer**

```hcl
# terraform-tenants/aws/rds.tf

resource "aws_db_instance" "building_shop" {
  identifier           = "building-shop-db"
  engine              = "postgres"
  engine_version      = "14.7"
  instance_class      = "db.t3.micro"      # Start small
  allocated_storage   = 20
  storage_encrypted   = true
  
  db_name  = "building_shop"
  username = "building_admin"
  password = var.db_password               # From secrets
  
  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name
  
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "mon:04:00-mon:05:00"
  
  skip_final_snapshot = false
  final_snapshot_identifier = "building-shop-final-snapshot"
  
  tags = {
    Customer = "Building Shop"
  }
}
```

**Option 2: Shared RDS with multiple databases**

```hcl
# One RDS instance
resource "aws_db_instance" "shared" {
  identifier        = "multi-tenant-db"
  instance_class    = "db.t3.small"     # Larger for multiple tenants
  allocated_storage = 100
  # ... other configs ...
}

# Create databases via null_resource (like local Terraform)
resource "null_resource" "create_customer_dbs" {
  for_each = var.customers
  
  provisioner "local-exec" {
    command = <<-EOT
      PGPASSWORD=${var.master_password} psql \
        -h ${aws_db_instance.shared.endpoint} \
        -U postgres \
        -c "CREATE DATABASE ${each.value.database_name};"
    EOT
  }
}
```

#### API Gateway

```hcl
# terraform-tenants/aws/api_gateway.tf

resource "aws_apigatewayv2_api" "building_shop" {
  name          = "building-shop-api"
  protocol_type = "HTTP"
  
  cors_configuration {
    allow_origins = ["https://buildingshop.com"]
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers = ["content-type", "authorization"]
  }
}

resource "aws_apigatewayv2_integration" "building_shop" {
  api_id           = aws_apigatewayv2_api.building_shop.id
  integration_type = "AWS_PROXY"
  
  integration_uri    = aws_lambda_function.building_shop_api.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "building_shop_default" {
  api_id    = aws_apigatewayv2_api.building_shop.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.building_shop.id}"
}

resource "aws_apigatewayv2_stage" "building_shop" {
  api_id      = aws_apigatewayv2_api.building_shop.id
  name        = "prod"
  auto_deploy = true
}

# Custom domain (optional)
resource "aws_apigatewayv2_domain_name" "building_shop" {
  domain_name = "api.buildingshop.com"
  
  domain_name_configuration {
    certificate_arn = aws_acm_certificate.building_shop.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}
```

#### Frontend (S3 + CloudFront)

```hcl
# terraform-tenants/aws/s3_frontend.tf

resource "aws_s3_bucket" "building_shop_frontend" {
  bucket = "buildingshop-com-frontend"
}

resource "aws_s3_bucket_website_configuration" "building_shop" {
  bucket = aws_s3_bucket.building_shop_frontend.id
  
  index_document {
    suffix = "index.html"
  }
  
  error_document {
    key = "index.html"  # For React Router
  }
}

resource "aws_cloudfront_distribution" "building_shop" {
  enabled             = true
  default_root_object = "index.html"
  
  origin {
    domain_name = aws_s3_bucket.building_shop_frontend.bucket_regional_domain_name
    origin_id   = "S3-buildingshop"
    
    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.building_shop.cloudfront_access_identity_path
    }
  }
  
  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "S3-buildingshop"
    viewer_protocol_policy = "redirect-to-https"
    
    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }
  
  aliases = ["buildingshop.com"]
  
  viewer_certificate {
    acm_certificate_arn = aws_acm_certificate.building_shop.arn
    ssl_support_method  = "sni-only"
  }
  
  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }
}
```

---

## 🔐 IAM Roles & Permissions

### Lambda Execution Role

```hcl
# terraform-tenants/aws/iam.tf

resource "aws_iam_role" "lambda_exec" {
  name = "lambda-tenant-execution-role"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_vpc" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Custom policy for RDS access
resource "aws_iam_role_policy" "lambda_rds" {
  name = "lambda-rds-access"
  role = aws_iam_role.lambda_exec.id
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "rds:DescribeDBInstances",
        "rds:Connect"
      ]
      Resource = "*"
    }]
  })
}
```

---

## 🔒 VPC & Security Groups

### Network Configuration

```hcl
# terraform-tenants/aws/vpc.tf

resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  
  tags = {
    Name = "multi-tenant-vpc"
  }
}

# Private subnets for Lambda
resource "aws_subnet" "private" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.${count.index + 1}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]
  
  tags = {
    Name = "private-subnet-${count.index + 1}"
  }
}

# Security group for Lambda
resource "aws_security_group" "lambda" {
  name        = "lambda-sg"
  vpc_id      = aws_vpc.main.id
  description = "Security group for Lambda functions"
  
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Security group for RDS
resource "aws_security_group" "rds" {
  name        = "rds-sg"
  vpc_id      = aws_vpc.main.id
  description = "Security group for RDS instances"
  
  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.lambda.id]
    description     = "PostgreSQL from Lambda"
  }
}
```

---

## 📦 Building Go Lambda Binary

### Makefile for Lambda Build

```makefile
# Add to your Makefile

.PHONY: build-lambda
build-lambda:
	@echo "Building Go binary for Lambda..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
		-tags lambda.norpc \
		-o bootstrap \
		cmd/lambda/main.go
	zip lambda.zip bootstrap
	@echo "Lambda package created: lambda.zip"

.PHONY: deploy-lambda-building-shop
deploy-lambda-building-shop: build-lambda
	@echo "Deploying to Building Shop Lambda..."
	aws lambda update-function-code \
		--function-name building-shop-api \
		--zip-file fileb://lambda.zip \
		--region us-east-1
```

---

## 🚀 Deployment Workflow

### Adding New Customer to AWS

#### Step 1: Run Interactive Script (Modified for AWS)

```bash
cd terraform-tenants
./add-customer-aws.sh
```

**Prompts:**
```
Customer Company Name: Building Shop Inc
Database Name: building_shop
AWS Region: us-east-1
RDS Option:
  1) Separate RDS instance (isolated, higher cost)
  2) Shared RDS instance (cost-effective)
Choice: 1

Admin Username: building_admin
Admin Password: ••••••••••
Admin Email: admin@buildingshop.com

Domain: buildingshop.com
API Domain: api.buildingshop.com
```

**Script generates:**
- `terraform-tenants/aws/customers/building_shop.tfvars`
- `terraform-tenants/aws/customers/building_shop-deploy.sh`

#### Step 2: Deploy Infrastructure

```bash
cd terraform-tenants/aws

# Initialize Terraform
terraform init

# Plan deployment
terraform plan -var-file="customers/building_shop.tfvars"

# Apply
terraform apply -var-file="customers/building_shop.tfvars"
```

**Terraform creates:**
- ✅ RDS PostgreSQL database
- ✅ Lambda function with environment variables
- ✅ API Gateway with custom domain
- ✅ S3 bucket for frontend
- ✅ CloudFront distribution
- ✅ Security groups, IAM roles

#### Step 3: Apply Database Schema

```bash
# Get RDS endpoint from Terraform output
RDS_ENDPOINT=$(terraform output -raw building_shop_rds_endpoint)

# Apply schema
PGPASSWORD=AdminPass123! psql \
  -h $RDS_ENDPOINT \
  -U building_admin \
  -d building_shop \
  -f ../modules/tenant/database-schema.sql
```

#### Step 4: Deploy Lambda Code

```bash
# Build Lambda binary
cd ../../
make build-lambda

# Deploy
cd terraform-tenants/aws
terraform apply -var-file="customers/building_shop.tfvars"
```

#### Step 5: Deploy Frontend to S3

```bash
# Build React frontend with customer config
cd react-frontend
REACT_APP_API_URL=https://api.buildingshop.com \
REACT_APP_TENANT_NAME="Building Shop" \
npm run build

# Upload to S3
aws s3 sync build/ s3://buildingshop-com-frontend/ --delete

# Invalidate CloudFront cache
aws cloudfront create-invalidation \
  --distribution-id E1234567890ABC \
  --paths "/*"
```

---

## 💰 Cost Estimation

### Per Customer (Pattern 1: Separate RDS)

| Service | Configuration | Monthly Cost |
|---------|---------------|--------------|
| **RDS** | db.t3.micro, 20GB | ~$15-25 |
| **Lambda** | 1M requests, 512MB, 200ms avg | ~$5-10 |
| **API Gateway** | 1M requests | ~$3.50 |
| **S3** | 10GB storage, 100GB transfer | ~$3 |
| **CloudFront** | 100GB transfer | ~$8.50 |
| **Route53** | Hosted zone | ~$0.50 |
| **Total** | | **~$35-50/month** |

### Shared RDS (Pattern 2: Multiple Databases)

| Service | Configuration | Monthly Cost |
|---------|---------------|--------------|
| **RDS** | db.t3.small, 100GB (shared) | ~$30-50 |
| **Lambda per customer** | 1M requests each | ~$5-10 each |
| **API Gateway per customer** | 1M requests each | ~$3.50 each |
| **S3 per customer** | 10GB, 100GB transfer each | ~$3 each |
| **CloudFront per customer** | 100GB transfer each | ~$8.50 each |
| **Total for 3 customers** | | **~$90-120/month** |

**Comparison:**
- 3 customers, separate RDS: ~$105-150/month
- 3 customers, shared RDS: ~$90-120/month

---

## 📊 Monitoring & Observability

### CloudWatch Integration

```hcl
# Lambda log group
resource "aws_cloudwatch_log_group" "building_shop_lambda" {
  name              = "/aws/lambda/building-shop-api"
  retention_in_days = 14
}

# Alarms
resource "aws_cloudwatch_metric_alarm" "lambda_errors" {
  alarm_name          = "building-shop-lambda-errors"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = "300"
  statistic           = "Sum"
  threshold           = "10"
  alarm_description   = "Alert on Lambda errors"
  
  dimensions = {
    FunctionName = aws_lambda_function.building_shop_api.function_name
  }
}

resource "aws_cloudwatch_metric_alarm" "rds_cpu" {
  alarm_name          = "building-shop-rds-cpu"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/RDS"
  period              = "300"
  statistic           = "Average"
  threshold           = "80"
  
  dimensions = {
    DBInstanceIdentifier = aws_db_instance.building_shop.id
  }
}
```

---

## 🔄 CI/CD Pipeline

### GitHub Actions for Lambda Deployment

```yaml
# .github/workflows/deploy-lambda.yml

name: Deploy to AWS Lambda

on:
  push:
    branches: [main]
    paths:
      - 'cmd/**'
      - 'internal/**'

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      
      - name: Build Lambda binary
        run: |
          GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
            -tags lambda.norpc \
            -o bootstrap \
            cmd/lambda/main.go
          zip lambda.zip bootstrap
      
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1
      
      - name: Deploy to all customer Lambdas
        run: |
          # Building Shop
          aws lambda update-function-code \
            --function-name building-shop-api \
            --zip-file fileb://lambda.zip
          
          # Hardware Store
          aws lambda update-function-code \
            --function-name hardware-store-api \
            --zip-file fileb://lambda.zip
          
          # Add more customers as needed
```

---

## 🎯 Migration from Local/Docker to AWS

### Phase 1: Set Up Infrastructure

1. Create AWS account
2. Set up Terraform backend (S3 + DynamoDB for state)
3. Deploy VPC, subnets, security groups

### Phase 2: Migrate First Customer

1. Create RDS instance
2. Export local database:
   ```bash
   pg_dump -U building_admin -d building_shop > building_shop.sql
   ```
3. Import to RDS:
   ```bash
   psql -h building-shop.xyz.rds.amazonaws.com -U building_admin -d building_shop < building_shop.sql
   ```
4. Deploy Lambda function
5. Test thoroughly
6. Update DNS to point to API Gateway

### Phase 3: Migrate Remaining Customers

Repeat for each customer.

---

## 🔐 Secrets Management

### Use AWS Secrets Manager

```hcl
# Store database credentials securely
resource "aws_secretsmanager_secret" "building_shop_db" {
  name = "building-shop/database"
}

resource "aws_secretsmanager_secret_version" "building_shop_db" {
  secret_id = aws_secretsmanager_secret.building_shop_db.id
  secret_string = jsonencode({
    username = "building_admin"
    password = var.db_password
    host     = aws_db_instance.building_shop.endpoint
    database = "building_shop"
  })
}

# Lambda accesses secrets
resource "aws_iam_role_policy" "lambda_secrets" {
  role = aws_iam_role.lambda_exec.id
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "secretsmanager:GetSecretValue"
      ]
      Resource = aws_secretsmanager_secret.building_shop_db.arn
    }]
  })
}
```

---

## ✅ Checklist for AWS Deployment

### Before Deployment

- [ ] AWS account created
- [ ] IAM user with appropriate permissions
- [ ] AWS CLI configured
- [ ] Terraform installed
- [ ] Go 1.25 installed
- [ ] Domain name registered (Route53 or external)
- [ ] SSL certificates created (ACM)

### Per Customer Deployment

- [ ] Run `add-customer-aws.sh` script
- [ ] Review generated Terraform configuration
- [ ] `terraform plan` - review changes
- [ ] `terraform apply` - deploy infrastructure
- [ ] Apply database schema to RDS
- [ ] Build and deploy Lambda function
- [ ] Build and deploy frontend to S3
- [ ] Configure DNS records
- [ ] Test API endpoints
- [ ] Test frontend application
- [ ] Set up monitoring alerts
- [ ] Document credentials securely

---

## 📚 Summary

### Key Points

1. **Same Application Code** - No code changes needed, just different deployment target
2. **Lambda Per Customer** - Each customer gets their own Lambda function
3. **RDS Options** - Choose separate or shared based on requirements
4. **Terraform Automation** - Infrastructure as code for AWS resources
5. **Cost-Effective** - Pay only for what you use, auto-scaling
6. **Secure** - VPC isolation, IAM roles, Secrets Manager

### Next Steps

1. Read this guide thoroughly
2. Set up AWS account and credentials
3. Create `terraform-tenants/aws/` directory structure
4. Start with one test customer
5. Migrate production customers gradually

---

**Documentation Links:**
- [AWS Lambda Go](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [API Gateway](https://docs.aws.amazon.com/apigateway/)
- [RDS PostgreSQL](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_PostgreSQL.html)

