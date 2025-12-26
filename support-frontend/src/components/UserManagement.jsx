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
  InputLabel
} from '@mui/material';
import {
  Delete as DeleteIcon,
  Refresh as RefreshIcon,
  People as PeopleIcon,
  ArrowBack as ArrowBackIcon,
  PersonAdd as PersonAddIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const API_URL = process.env.REACT_APP_SUPPORT_API_URL || 'http://localhost:4001';

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

  const fetchUsers = async () => {
    setLoading(true);
    setError('');
    
    try {
      const response = await axios.get(`${API_URL}/api/users`, {
        params: {
          page: page + 1, // Backend uses 1-based indexing
          page_size: rowsPerPage,
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

  useEffect(() => {
    fetchUsers();
  }, [page, rowsPerPage]);

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
      await axios.post(`${API_URL}/api/users`, {
        first_name: newUser.firstName,
        last_name: newUser.lastName,
        email: newUser.email,
        password: newUser.password,
        role: newUser.role
      });
      
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

  return (
    <Box>
      {/* AppBar */}
      <AppBar position="static">
        <Toolbar>
          <IconButton
            color="inherit"
            onClick={() => navigate('/support/dashboard')}
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
        {/* Total Count */}
        <Paper sx={{ p: 2, mb: 3 }}>
          <Typography variant="h6">
            Total Users: <strong>{totalCount}</strong>
          </Typography>
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
                <TableCell><strong>ID</strong></TableCell>
                <TableCell><strong>Name</strong></TableCell>
                <TableCell><strong>Email</strong></TableCell>
                <TableCell><strong>Role</strong></TableCell>
                <TableCell><strong>Created At</strong></TableCell>
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
                      No users found
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
                      {new Date(user.created_at).toLocaleDateString()}
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

