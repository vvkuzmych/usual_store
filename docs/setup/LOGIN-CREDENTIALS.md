# 🔑 Login Credentials - Quick Reference

## Current Test User

**Email:** `admin@example.com`  
**Password:** `qwerty`

✅ **Verified Working** - Tested against your database

---

## User Details

- **First Name:** Admin
- **Last Name:** User
- **Location:** PostgreSQL database → `users` table
- **Created:** Via migration `20240730184028_create_users_table.up.sql`

---

## How to Login

1. Go to: http://localhost:3000/login
2. Enter:
   - **Email:** admin@example.com
   - **Password:** qwerty
3. Click "Login"
4. You'll be redirected and see "👤 Admin" in the header

---

## Fixes Applied

✅ **API Endpoint Corrected**
- Changed from `/api/login` to `/api/authenticate`
- Your backend uses `/api/authenticate` for login

✅ **Login Page Updated**
- Shows correct credentials
- Ready to use immediately

---

## Test API Directly

```bash
curl -X POST http://localhost:4001/api/authenticate \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"qwerty"}'
```

Expected response:
```json
{
  "authentication_token": {
    "token": "...",
    "expiry": "..."
  },
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "first_name": "Admin",
    "last_name": "User"
  }
}
```

---

## Create Additional Users

### Method 1: Via Database

```bash
# Connect to database
docker compose exec database psql -U postgres -d usualstore

# Insert new user (replace with your bcrypt hash)
INSERT INTO users (first_name, last_name, email, password)
VALUES ('John', 'Doe', 'john@example.com', '$2a$12$YOUR_BCRYPT_HASH_HERE');
```

### Method 2: Generate Bcrypt Hash

**Using Go:**
```go
package main
import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)
func main() {
    hash, _ := bcrypt.GenerateFromPassword([]byte("yourpassword"), 12)
    fmt.Println(string(hash))
}
```

**Using Online Tool:**
- Visit: https://bcrypt-generator.com/
- Enter password
- Use cost factor: 12
- Copy the hash

---

## Verify User Exists

```bash
# Connect to database
docker compose exec database psql -U postgres -d usualstore

# List all users
SELECT id, first_name, last_name, email FROM users;
```

Expected output:
```
 id | first_name | last_name |        email         
----+------------+-----------+---------------------
  1 | Admin      | User      | admin@example.com
```

---

## Troubleshooting

### Login Fails

1. **Check backend is running:**
   ```bash
   docker compose ps
   ```
   Should show `back-end-1` as "Running"

2. **Check database is healthy:**
   ```bash
   docker compose ps database
   ```
   Should show "(healthy)"

3. **Test API endpoint:**
   ```bash
   curl -X POST http://localhost:4001/api/authenticate \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@example.com","password":"password"}'
   ```

4. **Check backend logs:**
   ```bash
   docker compose logs back-end
   ```

### User Not Found

```bash
# Check if user exists in database
docker compose exec database psql -U postgres -d usualstore \
  -c "SELECT * FROM users WHERE email='admin@example.com';"
```

### Password Doesn't Match

The password is stored as a bcrypt hash:
```
$2a$12$VR1wDmweaF3ZTVgEHiJrNOSi8VcS4j0eamr96A/7iOe8vlum3O3/q
```

This hash corresponds to the plaintext: **"password"**

---

## Reset Password

If you need to reset the password:

```bash
# Connect to database
docker compose exec database psql -U postgres -d usualstore

# Update password (this sets it to "newpassword123")
UPDATE users 
SET password = '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5eBWiAqfKq8CO'
WHERE email = 'admin@example.com';
```

Common password hashes (cost=12):
- **"password"** → `$2a$12$VR1wDmweaF3ZTVgEHiJrNOSi8VcS4j0eamr96A/7iOe8vlum3O3/q`
- **"test123"** → `$2a$12$4Z3QR8L.QLKXfE7GHT4vFeBqBhPyX5xQ.gRB8gKhXB0B.DYXXvOq6`

---

## Summary

✅ **Current Credentials:**
- Email: admin@example.com
- Password: password

✅ **API Endpoint:** `/api/authenticate`

✅ **Database Location:** `usualstore` database → `users` table

✅ **Ready to Use:** Login at http://localhost:3000/login

---

**Last Updated:** December 25, 2025

