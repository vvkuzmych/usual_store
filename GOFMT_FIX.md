# Fix: gofmt "Failed to exec spawn helper" Error

## 🚨 Problem

```
Cannot run program "/opt/homebrew/opt/go/libexec/bin/gofmt"
error=0, Failed to exec spawn helper: pid: 18430, exit value: 1
```

---

## ✅ Solution 1: Remove Extended Attributes (Quick Fix)

macOS sometimes blocks execution due to extended attributes:

```bash
# Remove attributes from gofmt
sudo xattr -cr /opt/homebrew/opt/go/libexec/bin/gofmt

# Verify it works
/opt/homebrew/opt/go/libexec/bin/gofmt -help
```

---

## ✅ Solution 2: Fix GoLand/IDE Configuration

### Update GOROOT in IDE:

1. **Open GoLand Preferences** (⌘ + ,)
2. Go to **Go** → **GOROOT**
3. Click **+** → **Local**
4. Select: `/opt/homebrew/opt/go/libexec`
5. Click **OK**

### Alternative: Set to use system gofmt:

1. **Preferences** → **Tools** → **File Watchers**
2. Find **gofmt** watcher
3. Change program path to: `/opt/homebrew/bin/gofmt`

---

## ✅ Solution 3: Invalidate IDE Caches

Sometimes IDE caches cause issues:

1. **File** → **Invalidate Caches**
2. Check **Clear file system cache and Local History**
3. Click **Invalidate and Restart**

---

## ✅ Solution 4: Reinstall Go Tools

If nothing works, reinstall:

```bash
# Reinstall Go
brew reinstall go

# Verify installation
go version
which gofmt
gofmt -help

# Restart IDE
```

---

## 🔍 Diagnosis Commands

```bash
# Check Go installation
go version
go env GOROOT

# Check gofmt location
which gofmt
ls -la /opt/homebrew/opt/go/libexec/bin/gofmt

# Check extended attributes
xattr /opt/homebrew/opt/go/libexec/bin/gofmt
xattr -l /opt/homebrew/opt/go/libexec/bin/gofmt

# Test gofmt directly
/opt/homebrew/opt/go/libexec/bin/gofmt -help
```

---

## 🎯 Current Status

Your system:
- ✅ Go installed: `go1.25.3 darwin/arm64`
- ✅ gofmt works in terminal: `/opt/homebrew/bin/gofmt`
- ⚠️  IDE path has extended attributes: `@` symbol

The issue is IDE-specific, not a Go installation problem.

---

## 💡 Recommended Fix Order

1. **Try Solution 1** (remove attributes) - 30 seconds
2. **Try Solution 3** (invalidate caches) - 1 minute
3. **Try Solution 2** (update IDE config) - 2 minutes
4. **Try Solution 4** (reinstall) - 5 minutes

---

## 📝 Quick Fix Script

Run this to fix the most common issues:

```bash
#!/bin/bash
echo "🔧 Fixing gofmt..."

# Remove extended attributes
sudo xattr -cr /opt/homebrew/opt/go/libexec/bin/

# Verify
if /opt/homebrew/opt/go/libexec/bin/gofmt -help > /dev/null 2>&1; then
    echo "✅ gofmt fixed!"
    echo "Now restart your IDE (GoLand/VSCode)"
else
    echo "❌ Still broken. Try reinstalling Go:"
    echo "brew reinstall go"
fi
```

Save as `fix_gofmt.sh` and run:
```bash
chmod +x fix_gofmt.sh
./fix_gofmt.sh
```

---

## 🚀 After Fix

1. Restart your IDE completely (not just reload)
2. Open your Go project
3. Try formatting a Go file (⌥ + ⇧ + F or ⌘ + ⌥ + L)
4. Should work now! ✅

---

**Most likely fix: Restart IDE after running Solution 1** 🔄

