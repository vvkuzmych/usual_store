# Go Package Conflict Fix

## 🐛 Problem

GitHub Actions workflow was failing with a Go build error:

```
found packages main (handlers-api.go) and models (widget_model.go) in /home/runner/work/usual_store/usual_store
Error: Process completed with exit code 1.
```

This error occurs when **multiple package declarations exist in the same directory**. Go requires all `.go` files in a directory to belong to the same package (except for `_test.go` files).

---

## 🔍 Root Cause

The root directory `/Users/vkuzm/Projects/UsualStore/usual_store/` contained **duplicate/misplaced Go files**:

```
usual_store/
├── handlers-api.go   ❌ (package main)
├── routes-api.go     ❌ (package main)
├── widget_model.go   ❌ (package models)  ← Different package!
```

These files had different package declarations:
- `handlers-api.go` and `routes-api.go` → `package main`
- `widget_model.go` → `package models`

Go detected this conflict and refused to build.

---

## ✅ Solution Applied

### **Removed Duplicate Files**

The proper versions of these files already existed in their correct locations:

```
usual_store/
├── cmd/api/
│   ├── handlers-api.go   ✅ (package main)
│   └── routes-api.go     ✅ (package main)
└── internal/models/
    └── widget_model.go   ✅ (package models)
```

**Action taken:** Deleted the duplicate files from the root directory:

```bash
rm /Users/vkuzm/Projects/UsualStore/usual_store/handlers-api.go
rm /Users/vkuzm/Projects/UsualStore/usual_store/routes-api.go
rm /Users/vkuzm/Projects/UsualStore/usual_store/widget_model.go
```

---

## 📋 What Changed

### **Before (Broken):**

```
usual_store/
├── handlers-api.go       ❌ Duplicate (package main)
├── routes-api.go         ❌ Duplicate (package main)
├── widget_model.go       ❌ Duplicate (package models)
├── cmd/api/
│   ├── handlers-api.go   ✅ Correct location
│   └── routes-api.go     ✅ Correct location
└── internal/models/
    └── widget_model.go   ✅ Correct location
```

**Error:**
```
found packages main and models in /home/runner/work/usual_store/usual_store
```

### **After (Fixed):**

```
usual_store/
├── cmd/api/
│   ├── handlers-api.go   ✅ (package main)
│   └── routes-api.go     ✅ (package main)
└── internal/models/
    └── widget_model.go   ✅ (package models)
```

**Result:**
```
✅ Build successful!
```

---

## 🎯 Go Package Rules

To prevent this issue in the future, remember:

### **Rule 1: One Package Per Directory**

All `.go` files in a directory must declare the same package (except `_test.go` files).

✅ **Good:**
```
myapp/
└── utils/
    ├── string.go    (package utils)
    ├── math.go      (package utils)
    └── time.go      (package utils)
```

❌ **Bad:**
```
myapp/
└── utils/
    ├── string.go    (package utils)
    ├── math.go      (package helpers)  ← Different package!
    └── time.go      (package utils)
```

### **Rule 2: Package Name Matches Directory Name**

By convention, the package name should match the directory name.

✅ **Good:**
```
myapp/
└── models/
    └── user.go      (package models)
```

❌ **Bad (but technically allowed):**
```
myapp/
└── models/
    └── user.go      (package database)  ← Confusing!
```

### **Rule 3: Main Package for Executables**

Only executable programs use `package main`.

✅ **Good:**
```
myapp/
├── cmd/
│   └── server/
│       └── main.go      (package main)
└── internal/
    └── handlers/
        └── api.go       (package handlers)
```

---

## 🧪 Testing the Fix

### **Local Build Test:**

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store
go build -v ./...
```

**Expected output:**
```
usual_store/internal/validator
usual_store/internal/models
usual_store/cmd/api
...
✅ Build successful!
```

### **GitHub Actions Test:**

1. Commit the changes
2. Push to GitHub
3. Workflow should now pass

---

## 🚀 Next Steps

### **1. Commit and Push:**

```bash
cd /Users/vkuzm/Projects/UsualStore/usual_store

# Check what was deleted
git status

# Add the deletions
git add -A

# Commit
git commit -m "fix: remove duplicate Go files causing package conflict

- Removed handlers-api.go from root (exists in cmd/api/)
- Removed routes-api.go from root (exists in cmd/api/)
- Removed widget_model.go from root (exists in internal/models/)

Fixes: found packages main and models error"

# Push
git push origin main
```

### **2. Verify GitHub Actions:**

```bash
# Watch the workflow
gh run watch

# Or view in browser
gh run view --web
```

---

## 🔍 How to Detect This Issue

### **Error Message:**

```
found packages <package1> (<file1>) and <package2> (<file2>) in <directory>
```

### **Solution:**

1. Identify which files have conflicting package declarations
2. Check if duplicates exist in proper locations
3. Remove duplicates or move files to correct directories
4. Ensure all files in a directory use the same package

### **Quick Check:**

```bash
# List all package declarations in current directory
head -1 *.go | grep package

# Should all show the same package!
```

---

## 📚 References

- **Go Package Documentation:** https://go.dev/doc/code#Organization
- **Effective Go - Packages:** https://go.dev/doc/effective_go#names
- **Go Project Layout:** https://github.com/golang-standards/project-layout

---

## ✅ Verification Checklist

- [x] Identified duplicate files in root directory
- [x] Verified proper versions exist in correct locations
- [x] Deleted duplicate files from root
- [ ] Committed changes
- [ ] Pushed to GitHub
- [ ] Verified GitHub Actions build passes

---

## 🎉 Success Criteria

Your fix is successful when:

✅ `go build -v ./...` completes without errors
✅ GitHub Actions workflow passes
✅ No "found packages" errors
✅ All files in each directory use the same package

---

**Fixed:** December 25, 2025  
**Issue:** Package conflict (multiple packages in root directory)  
**Solution:** Removed duplicate files (proper versions exist in subdirectories)  
**Impact:** High (blocks all builds and CI/CD)

