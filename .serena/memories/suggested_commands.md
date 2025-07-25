# PicoShare - Development Commands

## Essential Commands

### Building and Running
```bash
# Build Docker image (CRITICAL: May take multiple attempts)
docker build -t picoshare .

# Run the container
docker run -p 4001:4001 picoshare

# Stop current container
docker stop $(docker ps -q --filter ancestor=picoshare)

# View logs
docker logs $(docker ps -q --filter ancestor=picoshare)
```

### Development Server
```bash
# Development with hot reload (requires modd)
./dev-scripts/serve

# Run from source
PS_SHARED_SECRET=somesecretpass PORT=4001 go run cmd/picoshare/main.go
```

### Testing and Quality
```bash
# Run Go tests
./dev-scripts/run-go-tests

# Format frontend code
./dev-scripts/format-frontend

# Check Go formatting
./dev-scripts/check-go-formatting

# Run E2E tests
./dev-scripts/run-e2e-tests
```

### Database Management
```bash
# Reset database
./dev-scripts/reset-db

# Vacuum to reclaim space
sqlite3 data/store.db 'VACUUM'
```

## Important Notes
- **Always test changes in Docker** - Use IP 10.10.10.202:4001 for testing
- **Docker builds may timeout** - Simply restart the build command if it takes >3 minutes
- **Use Playwright MCP** for automated UI testing
- **Check authentication** with both shared secret and OAuth2 modes