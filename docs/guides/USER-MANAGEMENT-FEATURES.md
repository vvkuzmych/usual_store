# User Management Features

## Overview
The User Management page (`http://localhost:3005/support/users`) provides comprehensive tools for managing user accounts with advanced sorting, searching, and filtering capabilities.

## Features

### 1. Sortable Columns
All table columns are sortable by clicking on the column headers:

- **ID**: Sort users by their unique identifier
- **Name**: Sort alphabetically by full name (first + last)
- **Email**: Sort alphabetically by email address
- **Role**: Sort by role hierarchy (super_admin → admin → supporter → user)
- **Created At**: Sort by creation date (oldest → newest or newest → oldest)

**Usage:**
- First click: Sort ascending (↑)
- Second click: Sort descending (↓)
- Visual indicator shows active sort column and direction

### 2. Universal Search
Real-time search functionality across all user fields:

**Searchable Fields:**
- User ID (numeric search)
- First Name
- Last Name
- Full Name (combined)
- Email address
- Role (including display names like "Super Admin")
- Creation date (formatted date string)

**Features:**
- Real-time filtering as you type
- Search icon indicator
- Clear button (×) to reset search
- Shows filtered count: "Showing X filtered"
- Works seamlessly with pagination

**Examples:**
```
"admin"     → Shows all users with admin roles
"john"      → Shows all users named John (first or last name)
"@gmail"    → Shows all Gmail users
"Dec 26"    → Shows users created on December 26
"super"     → Shows super admin users
```

### 3. Fixed Date Display
Proper date formatting with error handling:

**Format:** `MMM DD, YYYY, HH:MM AM/PM`
**Example:** `Dec 26, 2025, 03:45 PM`

**Error Handling:**
- Invalid dates show "Invalid Date"
- Missing dates show "N/A"
- Proper timezone handling

### 4. Pagination
- Adjustable rows per page: 5, 10, 25, 50, 100
- Works with filtered results
- Automatic page reset when searching
- Shows total count and filtered count

## Access Control

### Super Admin Only
This page is restricted to users with the `super_admin` role.

**Access via:**
1. Login at: `http://localhost:3005/support/login`
2. Credentials: `admin@example.com` / `qwerty12`
3. Navigate to: Manage Users button in dashboard

## Technical Implementation

### Client-Side Operations
- **Sorting**: Client-side for instant response
- **Filtering**: Real-time search without server requests
- **Pagination**: Applied after sorting and filtering

### State Management
```javascript
const [orderBy, setOrderBy] = useState('id');
const [order, setOrder] = useState('asc');
const [searchTerm, setSearchTerm] = useState('');
const [filteredUsers, setFilteredUsers] = useState([]);
```

### Performance
- Fetches all users once (up to 1000)
- All operations performed client-side
- Instant sorting and filtering
- Efficient re-rendering with React hooks

## User Actions

### Available Actions
1. **Create User**: Add new users with any role
2. **Delete User**: Remove users (with super admin protection)
3. **Sort**: Order users by any column
4. **Search**: Filter users by any field
5. **Refresh**: Reload user list from database

### Protected Actions
- Cannot delete the last super admin
- Deletion requires confirmation dialog
- Warning shown when deleting super admins

## UI Components

### Material UI Components Used
- `TableSortLabel`: Sortable column headers
- `TextField`: Search input with adornments
- `InputAdornment`: Search icon and clear button
- `Chip`: Role badges with color coding
- `TablePagination`: Pagination controls
- `CircularProgress`: Loading indicator
- `Alert`: Error messages

### Role Colors
- **Super Admin**: Red (error)
- **Admin**: Orange (warning)
- **Supporter**: Blue (info)
- **User**: Gray (default)

## Browser Cache
After updates, users should:
1. Open in incognito window, OR
2. Hard refresh (Cmd/Ctrl + Shift + R), OR
3. Clear browser cache

## Related Files
- `support-frontend/src/components/UserManagement.jsx`
- `cmd/api/handlers-api.go` (GetAllUsers endpoint)
- `internal/models/models.go` (GetAllUsersPaginated)

## See Also
- [User Creation Email Notifications](USER-CREATION-EMAIL-NOTIFICATIONS.md)
- [Support Dashboard](SUPPORT-DASHBOARD-GUIDE.md)
- [Role-Based Access Control](RBAC-IMPLEMENTATION.md)

