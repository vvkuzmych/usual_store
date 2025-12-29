import React, { useState, useEffect } from 'react';
import {
  Box,
  Container,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  TableSortLabel,
  Chip,
  IconButton,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
  Alert,
  CircularProgress,
  AppBar,
  Toolbar,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  InputAdornment
} from '@mui/material';
import {
  Delete as DeleteIcon,
  Refresh as RefreshIcon,
  People as PeopleIcon,
  ArrowBack as ArrowBackIcon,
  PersonAdd as PersonAddIcon,
  Search as SearchIcon,
  Clear as ClearIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const API_URL = process.env.REACT_APP_SUPPORT_API_URL || 'http://localhost:4001';
const MESSAGING_API_URL = process.env.REACT_APP_MESSAGING_API_URL || 'http://localhost:4001'; // Backend API publishes to Kafka

function UserManagement() {
  const navigate = useNavigate();
  const [users, setUsers] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [deleteDialog, setDeleteDialog] = useState({ open: false, user: null });
  const [deleting, setDeleting] = useState(false);
  const [createDialog, setCreateDialog] = useState(false);
  const [creating, setCreating] = useState(false);
  const [newUser, setNewUser] = useState({
    firstName: '',
    lastName: '',
    email: '',
    password: '',
    role: 'user'
  });
  
  // Sorting state
  const [orderBy, setOrderBy] = useState('id');
  const [order, setOrder] = useState('asc');
  
  // Search state
  const [searchTerm, setSearchTerm] = useState('');
  const [searchDebounce, setSearchDebounce] = useState('');

  const fetchUsers = async () => {
    setLoading(true);
    setError('');
    
    try {
      // Fetch users with server-side filtering and sorting
      const response = await axios.get(`${API_URL}/api/users`, {
        params: {
          page: page + 1, // Backend uses 1-based indexing
          page_size: rowsPerPage,
          search: searchDebounce,
          sort_by: orderBy,
          sort_order: order,
        },
      });

      setUsers(response.data.users || []);
      setTotalCount(response.data.total_count || 0);
    } catch (err) {
      console.error('Error fetching users:', err);
      setError(err.response?.data?.message || 'Failed to load users');
    } finally {
      setLoading(false);
    }
  };

  // Fetch users when page, rowsPerPage, orderBy, order, or searchDebounce changes
  useEffect(() => {
    fetchUsers();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, rowsPerPage, orderBy, order, searchDebounce]);

  // Debounce search input (wait 500ms after user stops typing)
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      setSearchDebounce(searchTerm);
      setPage(0); // Reset to first page when search changes
    }, 500);

    return () => clearTimeout(timeoutId);
  }, [searchTerm]);

  // Handle sorting
  const handleRequestSort = (property) => {
    const isAsc = orderBy === property && order === 'asc';
    setOrder(isAsc ? 'desc' : 'asc');
    setOrderBy(property);
    setPage(0); // Reset to first page when sorting changes
  };

  // Handle search
  const handleSearch = (event) => {
    setSearchTerm(event.target.value);
  };

  const handleClearSearch = () => {
    setSearchTerm('');
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleDeleteClick = (user) => {
    setDeleteDialog({ open: true, user });
  };

  const handleDeleteConfirm = async () => {
    const userToDelete = deleteDialog.user;
    setDeleting(true);
    setError('');

    try {
      await axios.delete(`${API_URL}/api/users/${userToDelete.id}`);
      
      // Refresh the list
      await fetchUsers();
      
      setDeleteDialog({ open: false, user: null });
    } catch (err) {
      console.error('Error deleting user:', err);
      setError(err.response?.data?.message || 'Failed to delete user');
    } finally {
      setDeleting(false);
    }
  };

  const handleDeleteCancel = () => {
    setDeleteDialog({ open: false, user: null });
  };

  const handleCreateClick = () => {
    setNewUser({
      firstName: '',
      lastName: '',
      email: '',
      password: '',
      role: 'user'
    });
    setCreateDialog(true);
  };

  const handleCreateCancel = () => {
    setCreateDialog(false);
    setNewUser({
      firstName: '',
      lastName: '',
      email: '',
      password: '',
      role: 'user'
    });
  };

  const generatePassword = () => {
    const length = 12;
    const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*';
    let password = '';
    for (let i = 0; i < length; i++) {
      password += charset.charAt(Math.floor(Math.random() * charset.length));
    }
    setNewUser({ ...newUser, password });
  };

  const handleCreateSubmit = async () => {
    setCreating(true);
    setError('');

    try {
      const response = await axios.post(`${API_URL}/api/users`, {
        first_name: newUser.firstName,
        last_name: newUser.lastName,
        email: newUser.email,
        password: newUser.password,
        role: newUser.role
      });
      
      // Send welcome email via Kafka messaging service
      try {
        const roleDisplayName = {
          'super_admin': 'Super Administrator',
          'admin': 'Administrator',
          'supporter': 'Support Agent',
          'user': 'User'
        }[newUser.role] || 'User';

        // Determine if user needs dashboard access
        const isDashboardUser = ['super_admin', 'admin', 'supporter'].includes(newUser.role);
        const loginUrl = isDashboardUser 
          ? 'http://localhost:3005/support/dashboard' 
          : 'http://localhost:3000';

        // Build email message
        let emailMessage = `
══════════════════════════════════════════════════════════════
                USUAL STORE - ACCOUNT REGISTRATION
══════════════════════════════════════════════════════════════

Dear ${newUser.firstName} ${newUser.lastName},

Your account has been successfully created! Below are your credentials
and account information:

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                    ACCOUNT INFORMATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Full Name:  ${newUser.firstName} ${newUser.lastName}
Email:      ${newUser.email}
Password:   ${newUser.password}
Role:       ${roleDisplayName}

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                      ACCESS DETAILS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

You have been registered with the role of: ${roleDisplayName}
`;

        if (isDashboardUser) {
          emailMessage += `
Your role grants you access to the Support Dashboard where you can:
${newUser.role === 'super_admin' ? '  • Manage all users and permissions\n  • Create administrators and supporters\n  • View and manage all support tickets\n  • Access system analytics and reports' : ''}${newUser.role === 'admin' ? '  • Create support staff accounts\n  • View and manage support tickets\n  • Access system reports' : ''}${newUser.role === 'supporter' ? '  • View and respond to support tickets\n  • Assist customers with their inquiries' : ''}

Dashboard URL: ${loginUrl}

Please bookmark this URL for easy access.
`;
        } else {
          emailMessage += `
You can now log in to shop for products and services.

Store URL: ${loginUrl}
`;
        }

        emailMessage += `
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                    SECURITY REMINDER
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

⚠️  IMPORTANT: For your security, please change your password after
    your first login.

⚠️  Keep your credentials confidential and do not share them with
    anyone.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

If you have any questions or need assistance, please don't hesitate
to contact our support team.

Best regards,
The Usual Store Team

══════════════════════════════════════════════════════════════
        `;

        await axios.post(`${MESSAGING_API_URL}/api/messaging/send`, {
          to: newUser.email,
          subject: `Welcome to Usual Store - Your ${roleDisplayName} Account Has Been Created`,
          message: emailMessage.trim(),
        });
        console.log('✅ Welcome email sent successfully via Kafka');
      } catch (emailError) {
        console.error('⚠️ Failed to send welcome email:', emailError);
        // Don't fail the whole process if email fails
      }
      
      // Refresh the list
      await fetchUsers();
      
      setCreateDialog(false);
      setNewUser({
        firstName: '',
        lastName: '',
        email: '',
        password: '',
        role: 'user'
      });
    } catch (err) {
      console.error('Error creating user:', err);
      setError(err.response?.data?.message || 'Failed to create user');
    } finally {
      setCreating(false);
    }
  };

  const getRoleColor = (role) => {
    switch (role) {
      case 'super_admin':
        return 'error';
      case 'admin':
        return 'warning';
      case 'supporter':
        return 'info';
      default:
        return 'default';
    }
  };

  const getRoleLabel = (role) => {
    switch (role) {
      case 'super_admin':
        return 'Super Admin';
      case 'admin':
        return 'Admin';
      case 'supporter':
        return 'Supporter';
      default:
        return 'User';
    }
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    
    try {
      const date = new Date(dateString);
      
      // Check if date is valid
      if (isNaN(date.getTime())) {
        return 'Invalid Date';
      }
      
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch (error) {
      return 'Invalid Date';
    }
  };

  return (
    <Box>
      {/* AppBar */}
      <AppBar position="static">
        <Toolbar>
          <IconButton
            color="inherit"
            onClick={() => window.location.href = 'http://localhost:8000/support/dashboard'}
            sx={{ mr: 2 }}
          >
            <ArrowBackIcon />
          </IconButton>
          <PeopleIcon sx={{ mr: 2 }} />
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            User Management
          </Typography>
          <Button
            color="inherit"
            startIcon={<PersonAddIcon />}
            onClick={handleCreateClick}
            sx={{ mr: 2 }}
          >
            Create User
          </Button>
          <Button
            color="inherit"
            startIcon={<RefreshIcon />}
            onClick={fetchUsers}
            disabled={loading}
          >
            Refresh
          </Button>
        </Toolbar>
      </AppBar>

      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        {/* Total Count and Search */}
        <Paper sx={{ p: 2, mb: 3 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 2 }}>
            <Typography variant="h6">
              {searchTerm ? (
                <>
                  Found: <strong>{totalCount}</strong> user{totalCount !== 1 ? 's' : ''}
                  <Typography component="span" variant="body2" color="text.secondary" sx={{ ml: 2 }}>
                    (filtered)
                  </Typography>
                </>
              ) : (
                <>
                  Total Users: <strong>{totalCount}</strong>
                </>
              )}
            </Typography>
            <TextField
              size="small"
              placeholder="Search by ID, name, email, role, or date..."
              value={searchTerm}
              onChange={handleSearch}
              sx={{ minWidth: 300 }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
                endAdornment: searchTerm && (
                  <InputAdornment position="end">
                    <IconButton size="small" onClick={handleClearSearch}>
                      <ClearIcon />
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
          </Box>
        </Paper>

        {/* Error Alert */}
        {error && (
          <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError('')}>
            {error}
          </Alert>
        )}

        {/* Users Table */}
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>
                  <TableSortLabel
                    active={orderBy === 'id'}
                    direction={orderBy === 'id' ? order : 'asc'}
                    onClick={() => handleRequestSort('id')}
                  >
                    <strong>ID</strong>
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={orderBy === 'name'}
                    direction={orderBy === 'name' ? order : 'asc'}
                    onClick={() => handleRequestSort('name')}
                  >
                    <strong>Name</strong>
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={orderBy === 'email'}
                    direction={orderBy === 'email' ? order : 'asc'}
                    onClick={() => handleRequestSort('email')}
                  >
                    <strong>Email</strong>
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={orderBy === 'role'}
                    direction={orderBy === 'role' ? order : 'asc'}
                    onClick={() => handleRequestSort('role')}
                  >
                    <strong>Role</strong>
                  </TableSortLabel>
                </TableCell>
                <TableCell>
                  <TableSortLabel
                    active={orderBy === 'created_at'}
                    direction={orderBy === 'created_at' ? order : 'asc'}
                    onClick={() => handleRequestSort('created_at')}
                  >
                    <strong>Created At</strong>
                  </TableSortLabel>
                </TableCell>
                <TableCell align="center"><strong>Actions</strong></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : users.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                    <Typography variant="body1" color="text.secondary">
                      {searchTerm ? 'No users match your search' : 'No users found'}
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                users.map((user) => (
                  <TableRow key={user.id} hover>
                    <TableCell>{user.id}</TableCell>
                    <TableCell>
                      {user.first_name} {user.last_name}
                    </TableCell>
                    <TableCell>{user.email}</TableCell>
                    <TableCell>
                      <Chip
                        label={getRoleLabel(user.role)}
                        color={getRoleColor(user.role)}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      {formatDate(user.created_at)}
                    </TableCell>
                    <TableCell align="center">
                      <IconButton
                        color="error"
                        onClick={() => handleDeleteClick(user)}
                        size="small"
                      >
                        <DeleteIcon />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>

          {/* Pagination */}
          <TablePagination
            component="div"
            count={totalCount}
            page={page}
            onPageChange={handleChangePage}
            rowsPerPage={rowsPerPage}
            onRowsPerPageChange={handleChangeRowsPerPage}
            rowsPerPageOptions={[5, 10, 25, 50, 100]}
          />
        </TableContainer>
      </Container>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialog.open} onClose={handleDeleteCancel}>
        <DialogTitle>Confirm Delete</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete user <strong>{deleteDialog.user?.email}</strong>?
            {deleteDialog.user?.role === 'super_admin' && (
              <Alert severity="warning" sx={{ mt: 2 }}>
                Warning: You are about to delete a Super Admin account!
              </Alert>
            )}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDeleteCancel} disabled={deleting}>
            Cancel
          </Button>
          <Button
            onClick={handleDeleteConfirm}
            color="error"
            variant="contained"
            disabled={deleting}
          >
            {deleting ? (
              <>
                <CircularProgress size={20} sx={{ mr: 1 }} />
                Deleting...
              </>
            ) : (
              'Delete'
            )}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Create User Dialog */}
      <Dialog open={createDialog} onClose={handleCreateCancel} maxWidth="sm" fullWidth>
        <DialogTitle>Create New User</DialogTitle>
        <DialogContent>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          <TextField
            fullWidth
            label="First Name"
            value={newUser.firstName}
            onChange={(e) => setNewUser({ ...newUser, firstName: e.target.value })}
            margin="normal"
            required
            disabled={creating}
          />

          <TextField
            fullWidth
            label="Last Name"
            value={newUser.lastName}
            onChange={(e) => setNewUser({ ...newUser, lastName: e.target.value })}
            margin="normal"
            required
            disabled={creating}
          />

          <TextField
            fullWidth
            label="Email Address"
            type="email"
            value={newUser.email}
            onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
            margin="normal"
            required
            disabled={creating}
          />

          <Box sx={{ mt: 2, mb: 1 }}>
            <TextField
              fullWidth
              label="Password"
              type="text"
              value={newUser.password}
              onChange={(e) => setNewUser({ ...newUser, password: e.target.value })}
              required
              disabled={creating}
              helperText="User will receive this password via email"
            />
            <Button
              size="small"
              onClick={generatePassword}
              disabled={creating}
              sx={{ mt: 1 }}
            >
              Generate Secure Password
            </Button>
          </Box>

          <FormControl fullWidth margin="normal" required>
            <InputLabel>Role</InputLabel>
            <Select
              value={newUser.role}
              label="Role"
              onChange={(e) => setNewUser({ ...newUser, role: e.target.value })}
              disabled={creating}
            >
              <MenuItem value="user">User (Regular Customer)</MenuItem>
              <MenuItem value="supporter">Supporter (Support Agent)</MenuItem>
              <MenuItem value="admin">Admin (Administrator)</MenuItem>
              <MenuItem value="super_admin">Super Admin (Full Access)</MenuItem>
            </Select>
          </FormControl>

          {newUser.role === 'super_admin' && (
            <Alert severity="warning" sx={{ mt: 2 }}>
              Warning: Super Admins have full access to all system features, including user management!
            </Alert>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCreateCancel} disabled={creating}>
            Cancel
          </Button>
          <Button
            onClick={handleCreateSubmit}
            variant="contained"
            disabled={creating || !newUser.email || !newUser.password || !newUser.firstName || !newUser.lastName}
          >
            {creating ? (
              <>
                <CircularProgress size={20} sx={{ mr: 1 }} />
                Creating...
              </>
            ) : (
              'Create User'
            )}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default UserManagement;

