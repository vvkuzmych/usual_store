# UsualStore Database Information

## Database Type

**PostgreSQL 15**

## Architecture

The database runs entirely in Docker and is managed by Terraform.

```
┌─────────────────────────────────────────────────┐
│              Docker Host (Your Mac)              │
│                                                  │
│  ┌────────────────────────────────────────┐    │
│  │   usualstore-database Container        │    │
│  │                                        │    │
│  │   PostgreSQL 15                        │    │
│  │   Database: usualstore                 │    │
│  │   User: postgres                       │    │
│  │   Internal Port: 5432                  │    │
│  │                                        │    │
│  │   Data: /var/lib/postgresql/data      │    │
│  │   (Mapped to Docker volume)            │    │
│  └────────────────────────────────────────┘    │
│              ↕                                   │
│    Port Mapping: 5433 → 5432                    │
│              ↕                                   │
│  ┌────────────────────────────────────────┐    │
│  │  Application Containers                 │    │
│  │  - usualstore-back-end                  │    │
│  │  - usualstore-ai-assistant              │    │
│  │  - usualstore-support-service           │    │
│  └────────────────────────────────────────┘    │
│                                                  │
│  External Access: localhost:5433                │
└─────────────────────────────────────────────────┘
```

## Installation

### ✅ What You DON'T Need to Install

- ❌ PostgreSQL server
- ❌ Database initialization scripts
- ❌ Homebrew PostgreSQL service
- ❌ Manual database creation
- ❌ Port configuration

**The database runs completely in Docker!**

### ⚠️ What You DO Need

#### 1. Docker Desktop
```bash
brew install --cask docker
```

#### 2. Terraform
```bash
brew install terraform
```

#### 3. Migration Tool (soda)
```bash
brew install gobuffalo/tap/pop
```

#### 4. PostgreSQL Client (Optional)
For direct database access:
```bash
# Option A: Just the client (recommended)
brew install libpq
echo 'export PATH="/opt/homebrew/opt/libpq/bin:$PATH"' >> ~/.zshrc

# Option B: Full PostgreSQL (includes client)
brew install postgresql@15
```

## Database Configuration

### Connection Details

| Setting | Value |
|---------|-------|
| **Host** | `localhost` (from host) or `database` (from containers) |
| **Port** | `5433` (external) / `5432` (internal) |
| **Database** | `usualstore` |
| **User** | `postgres` |
| **Password** | Set in `terraform.tfvars` |

### Connection Strings

**From Host Machine**:
```
postgres://postgres:your_password@localhost:5433/usualstore
```

**From Docker Containers** (internal network):
```
postgres://postgres:your_password@database:5432/usualstore
```

**Environment Variable**:
```bash
export DATABASE_URL="postgres://postgres:your_password@localhost:5433/usualstore"
```

## Schema Management

### Migrations

**Tool**: Soda (Pop migrations)

**Location**: `/Users/vkuzm/Projects/UsualStore/usual_store/migrations/`

**Total**: 24 migration pairs (48 files)

### Tables Created

After running migrations (`make migrate`), the following tables are created:

#### Core Tables
- `users` - User accounts (admin, super_admin, supporter, user)
- `widgets` - Products/items for sale
- `orders` - Customer orders
- `transactions` - Payment transactions
- `transaction_status` - Transaction states
- `statuses` - Order status tracking
- `customers` - Customer information
- `tokens` - Authentication tokens
- `sessions` - User sessions

#### AI Assistant Tables
- `ai_conversations` - Chat sessions
- `ai_messages` - Individual messages
- `ai_user_preferences` - User AI preferences
- `ai_feedback` - User feedback on responses
- `ai_product_cache` - Product recommendation cache

#### Support System Tables
- `support_tickets` - Support tickets
- `support_messages` - Ticket messages
- `support_sessions` - WebSocket sessions

#### System Tables
- `schema_migration` - Migration tracking (managed by soda)

## Access Methods

### 1. Using psql (PostgreSQL Client)

```bash
# Connect from host
psql "postgres://postgres:password@localhost:5433/usualstore"

# Or using environment variable
psql $DATABASE_URL

# Common commands
\dt              # List all tables
\d users         # Describe users table
\l               # List databases
\q               # Quit
```

### 2. Using Docker exec

```bash
# Connect via Docker
docker exec -it usualstore-database psql -U postgres -d usualstore

# Run a query
docker exec -it usualstore-database psql -U postgres -d usualstore -c "SELECT * FROM users;"

# Execute SQL file
docker exec -i usualstore-database psql -U postgres -d usualstore < script.sql
```

### 3. Using Makefile Commands

```bash
# IPv6 connection
make db-shell-ipv6

# IPv4 connection
make db-shell-ipv4

# Docker internal connection
make db-shell-docker
```

### 4. Using GUI Tools

Popular PostgreSQL GUI clients:
- **pgAdmin** (https://www.pgadmin.org/)
- **TablePlus** (https://tableplus.com/)
- **DBeaver** (https://dbeaver.io/)
- **DataGrip** (https://www.jetbrains.com/datagrip/)

**Connection Settings**:
- Host: `localhost`
- Port: `5433`
- Database: `usualstore`
- User: `postgres`
- Password: (from terraform.tfvars)

## Data Persistence

### Docker Volume

Data is persisted in a Docker volume: `usualstore_db_data`

```bash
# Check volume
docker volume inspect usualstore_db_data

# Backup volume
docker run --rm -v usualstore_db_data:/data -v $(pwd):/backup alpine tar czf /backup/db-backup.tar.gz /data

# Restore volume
docker run --rm -v usualstore_db_data:/data -v $(pwd):/backup alpine sh -c "cd / && tar xzf /backup/db-backup.tar.gz"
```

### Data Survives

✅ Container restart
✅ Container recreation
✅ Terraform apply/destroy (unless volume deleted)

### Data is Lost

❌ `docker volume rm usualstore_db_data`
❌ `terraform destroy` with volume removal
❌ `docker system prune --volumes`

## Backup & Restore

### Backup Database

```bash
# Full database backup
docker exec usualstore-database pg_dump -U postgres usualstore > backup.sql

# Compressed backup
docker exec usualstore-database pg_dump -U postgres usualstore | gzip > backup.sql.gz

# Backup specific table
docker exec usualstore-database pg_dump -U postgres -t users usualstore > users-backup.sql
```

### Restore Database

```bash
# Restore from backup
docker exec -i usualstore-database psql -U postgres -d usualstore < backup.sql

# Restore compressed backup
gunzip < backup.sql.gz | docker exec -i usualstore-database psql -U postgres -d usualstore
```

## Performance

### Configuration

The database is configured with:
- `max_connections=200`
- `listen_addresses=*`
- Health checks every 10 seconds
- 5-second timeout
- 5 retries before unhealthy

### Monitoring

```bash
# Check database status
docker exec usualstore-database pg_isready -U postgres

# View active connections
docker exec -it usualstore-database psql -U postgres -d usualstore -c "SELECT * FROM pg_stat_activity;"

# Database size
docker exec -it usualstore-database psql -U postgres -d usualstore -c "SELECT pg_size_pretty(pg_database_size('usualstore'));"

# Table sizes
docker exec -it usualstore-database psql -U postgres -d usualstore -c "SELECT schemaname,tablename,pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size FROM pg_tables WHERE schemaname='public' ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"
```

## Troubleshooting

### Database Won't Start

```bash
# Check logs
docker logs usualstore-database

# Check health
docker inspect usualstore-database | grep -A 10 Health

# Restart
docker restart usualstore-database
```

### Cannot Connect

```bash
# Verify port is listening
lsof -i :5433

# Check from Docker network
docker exec usualstore-back-end nc -zv database 5432

# Test connection
docker exec usualstore-database pg_isready -U postgres -d usualstore
```

### Migration Errors

```bash
# Check if tables exist
docker exec -it usualstore-database psql -U postgres -d usualstore -c "\dt"

# Run migrations
make migrate

# Check migration status
soda migrate status
```

### Reset Database (⚠️ DESTRUCTIVE)

```bash
# Option 1: Recreate container (keeps volume)
terraform destroy -target=module.database.docker_container.database
terraform apply -target=module.database.docker_container.database
make migrate

# Option 2: Delete volume (loses all data)
docker stop usualstore-database
docker rm usualstore-database
docker volume rm usualstore_db_data
terraform apply
make migrate
```

## Security

### Recommendations

1. **Change Default Password**
   - Edit `terraform.tfvars`
   - Set strong `postgres_password`

2. **Restrict Access**
   - Database only accessible from localhost by default
   - Internal Docker network for container-to-container

3. **Backup Regularly**
   - Use `pg_dump` for backups
   - Store backups securely
   - Test restore procedures

4. **Production**
   - Use managed PostgreSQL (AWS RDS, Azure Database, etc.)
   - Enable SSL/TLS
   - Set up replication
   - Configure automated backups

## Production Considerations

### For Production Deployment

**Don't use Docker PostgreSQL in production!**

Use managed database services:
- **AWS**: Amazon RDS for PostgreSQL
- **Azure**: Azure Database for PostgreSQL
- **GCP**: Cloud SQL for PostgreSQL
- **Heroku**: Heroku Postgres
- **DigitalOcean**: Managed PostgreSQL

**Why?**
- Automated backups
- High availability
- Automatic failover
- Monitoring & alerts
- Scaling
- Security patches
- Professional support

### Migration to Production

1. Export schema:
   ```bash
   pg_dump -U postgres -s usualstore > schema.sql
   ```

2. Export data:
   ```bash
   pg_dump -U postgres -a usualstore > data.sql
   ```

3. Import to production database

4. Update connection strings in application

---

**Related Documentation**:
- `SETUP-GUIDE.md` - Complete setup instructions
- `QUICK-START.md` - Quick reference
- `DEVELOPMENT-WORKFLOW.md` - Development practices

**Last Updated**: December 31, 2025

