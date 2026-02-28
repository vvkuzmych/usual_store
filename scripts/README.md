# Scripts for usual_store

This directory contains helper scripts for testing and verification.

---

## 🔍 Worker Pool Health Check Scripts

### 1. `check_worker_pool.sh` - **Comprehensive Check**

Full health check of the worker pool implementation with detailed reporting.

**Usage:**
```bash
./scripts/check_worker_pool.sh
```

**What it checks:**
- ✅ Worker pool package exists
- ✅ All required files present
- ✅ All tests passing (runs full test suite)
- ✅ Test coverage (reports %)
- ✅ Binary builds successfully
- ✅ Configuration flags available
- ✅ Code quality (imports, integration)
- 📊 Implementation summary
- 📖 Configuration options
- 💡 Usage examples

**Duration:** ~45 seconds

**When to use:**
- Before committing code
- After making changes to worker pool
- Weekly health checks
- Before deployments

---

### 2. `quick_check.sh` - **Fast Verification**

Quick smoke test to verify worker pool is functional.

**Usage:**
```bash
./scripts/quick_check.sh
```

**What it checks:**
- ✅ Basic test passing
- ✅ Build succeeds
- 📊 Coverage percentage

**Duration:** ~5 seconds

**When to use:**
- During development (quick feedback)
- After small changes
- Before pushing to git
- In CI/CD pipelines (fast gate)

---

## 📊 Example Output

### Comprehensive Check:
```
========================================
  Worker Pool Health Check
========================================

▶ 1. Checking Worker Pool Package
-------------------------------------------
✅ Worker pool package exists

▶ 2. Verifying Required Files
-------------------------------------------
✅ internal/workerpool/workerpool.go exists
✅ internal/workerpool/job.go exists
...

✅ Worker Pool Health Check Complete!
All checks passed successfully!
Worker pool is ready for production use.
```

### Quick Check:
```
🔍 Quick Worker Pool Check...

▶ Running tests...
✅ Tests passing
▶ Building service...
✅ Build successful
▶ Checking coverage...
   Coverage: 97.7%

✅ Worker pool is healthy!
```

### 3. `worker_pool_info.sh` - **Information & Status**

Shows worker pool configuration and stats without running tests.

**Usage:**
```bash
./scripts/worker_pool_info.sh
```

**What it shows:**
- 📦 Implementation stats (lines of code, files)
- 🧪 Test statistics
- ⚙️ Configuration defaults
- 🔗 Integration points
- 💡 Usage examples
- 📊 Performance estimates

**Duration:** < 1 second

**When to use:**
- Quick info lookup
- Documentation reference
- Onboarding new developers
- Before making configuration changes

---

## 🛠️ Other Scripts

### `start-typescript.sh`
Starts TypeScript services.

### `test-ipv6.sh`
Tests IPv6 connectivity.

---

## 📝 Adding New Scripts

When adding new scripts:
1. Make them executable: `chmod +x scripts/your_script.sh`
2. Add shebang: `#!/bin/bash`
3. Add description comment at top
4. Update this README
5. Test before committing

---

## 🔧 Troubleshooting

**Scripts not executable?**
```bash
chmod +x scripts/*.sh
```

**Wrong working directory?**
Scripts automatically `cd` to project root.

**Tests failing?**
Run comprehensive check for detailed error messages:
```bash
./scripts/check_worker_pool.sh
```
