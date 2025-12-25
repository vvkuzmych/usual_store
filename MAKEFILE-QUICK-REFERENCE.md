# Makefile Quick Reference 🚀

## Frontend Commands

### Start Frontends
```bash
make react-start          # React (port 3000) + backend
make typescript-start     # TypeScript (port 3001) + backend
make go-start            # Go (port 4000) + backend
make all-frontends-start # ALL frontends + backend
```

### Stop Frontends
```bash
make react-stop          # Stop React
make typescript-stop     # Stop TypeScript
make go-stop            # Stop Go
make docker-down        # Stop everything
```

### View Logs
```bash
make react-logs         # React logs
make typescript-logs    # TypeScript logs
make go-logs           # Go logs
make docker-logs       # All logs
```

### Build Images
```bash
make react-build            # Build React image
make typescript-build       # Build TypeScript image
make go-build-docker        # Build Go image
make build-all-frontends    # Build all images
```

### Status & Restart
```bash
make frontend-status    # Show frontend status
make react-restart      # Restart React
make typescript-restart # Restart TypeScript
make go-restart        # Restart Go
```

---

## Access URLs

| Frontend | Command | URL |
|----------|---------|-----|
| React | `make react-start` | http://localhost:3000 |
| TypeScript | `make typescript-start` | http://localhost:3001 |
| Go | `make go-start` | http://localhost:4000 |
| Backend | (auto-started) | http://localhost:4001 |

---

## Common Workflows

### Quick Start TypeScript
```bash
make typescript-start    # Start
make typescript-logs     # View logs
make typescript-stop     # Stop
```

### Compare All Frontends
```bash
make all-frontends-start
make frontend-status
# Visit: 3000, 3001, 4000
make docker-down
```

### Development
```bash
make typescript-start    # Start
# ... make changes ...
make typescript-build    # Rebuild
make typescript-restart  # Restart
make typescript-logs     # Check logs
```

---

## Database Commands
```bash
make db-shell-ipv4      # Connect via IPv4
make migrate            # Run migrations
make rollback           # Rollback migrations
```

---

## Help
```bash
make help              # Show all commands
```

---

**Full documentation**: `docs/guides/MAKEFILE-COMMANDS.md`

