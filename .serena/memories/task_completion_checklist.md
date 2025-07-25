# PicoShare - Task Completion Checklist  

## After Code Changes

### 1. Docker Build and Test
```bash
# Build Docker image - CRITICAL step
docker build -t picoshare .

# If build takes >3 minutes, cancel and retry
# Run the new image
docker run -p 4001:4001 picoshare

# Test at http://10.10.10.202:4001 (not localhost!)
```

### 2. Functionality Verification
- [ ] **Basic functionality** works as expected
- [ ] **Dark/light theme toggle** functions properly
- [ ] **Authentication flow** works (both shared secret and OAuth2)
- [ ] **File uploads** and downloads work
- [ ] **Responsive design** works on mobile/desktop

### 3. Code Quality Checks
```bash
# Format Go code
./dev-scripts/check-go-formatting

# Format frontend code  
./dev-scripts/format-frontend

# Run Go tests
./dev-scripts/run-go-tests
```

### 4. Automated Testing (When Available)
```bash
# Run E2E tests with Playwright MCP
./dev-scripts/run-e2e-tests

# Test specific user flows
# - File upload/download
# - Authentication
# - Theme switching
# - Thumbnail generation
```

### 5. Security Checks
- [ ] **No sensitive data** in logs or responses
- [ ] **Authentication** properly enforced
- [ ] **Input validation** for all user inputs
- [ ] **File upload restrictions** working
- [ ] **CSRF protection** in place

### 6. Performance Verification
- [ ] **Thumbnail generation** works efficiently
- [ ] **File serving** has proper caching headers
- [ ] **Database queries** are optimized
- [ ] **Static assets** load correctly

## Before Completing Task
- [ ] All modified files compile without errors
- [ ] Docker image builds successfully
- [ ] Application runs at http://10.10.10.202:4001
- [ ] All required features work as specified
- [ ] No breaking changes to existing functionality