import React, { useState } from 'react';
import { useNavigate, useLocation, Link as RouterLink } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import {
  Box,
  Container,
  Paper,
  Typography,
  TextField,
  Button,
  Link,
  Alert,
  Checkbox,
  FormControlLabel,
  CircularProgress,
  Divider,
} from '@mui/material';
import {
  Login as LoginIcon,
  PersonAdd as PersonAddIcon,
} from '@mui/icons-material';

const Login = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const from = location.state?.from || '/';

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
    if (error) setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.email || !formData.password) {
      setError('Please enter both email and password');
      return;
    }

    setLoading(true);
    setError('');

    const result = await login(formData.email, formData.password);
    
    if (result.success) {
      navigate(from, { replace: true });
    } else {
      setError(result.error);
      setLoading(false);
    }
  };

  return (
    <Box
      sx={{
        minHeight: '80vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        py: 4,
      }}
    >
      <Container maxWidth="sm">
        <Paper
          elevation={10}
          sx={{
            p: 4,
            borderRadius: 3,
          }}
        >
          {/* Header */}
          <Box sx={{ textAlign: 'center', mb: 3 }}>
            <LoginIcon sx={{ fontSize: 48, color: 'primary.main', mb: 1 }} />
            <Typography variant="h4" component="h1" gutterBottom fontWeight="bold">
              Welcome Back
            </Typography>
            <Typography variant="body1" color="text.secondary">
              Login to your account to continue shopping
            </Typography>
          </Box>

          {/* Error Alert */}
          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          {/* Login Form */}
          <Box component="form" onSubmit={handleSubmit} noValidate>
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="Email Address"
              name="email"
              autoComplete="email"
              autoFocus
              value={formData.email}
              onChange={handleChange}
              disabled={loading}
            />

            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type="password"
              id="password"
              autoComplete="current-password"
              value={formData.password}
              onChange={handleChange}
              disabled={loading}
            />

            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', my: 2 }}>
              <FormControlLabel
                control={<Checkbox value="remember" color="primary" />}
                label="Remember me"
              />
              <Link
                component={RouterLink}
                to="/forgot-password"
                variant="body2"
                sx={{ textDecoration: 'none' }}
              >
                Forgot password?
              </Link>
            </Box>

            <Button
              type="submit"
              fullWidth
              variant="contained"
              size="large"
              disabled={loading}
              sx={{ mt: 2, mb: 2, py: 1.5 }}
            >
              {loading ? <CircularProgress size={24} /> : 'Login'}
            </Button>

            <Divider sx={{ my: 3 }}>
              <Typography variant="body2" color="text.secondary">
                Don't have an account?
              </Typography>
            </Divider>

            <Button
              fullWidth
              variant="outlined"
              size="large"
              component={RouterLink}
              to="/signup"
              state={{ from }}
              startIcon={<PersonAddIcon />}
              sx={{ py: 1.5 }}
            >
              Create Account
            </Button>
          </Box>

          {/* Demo Credentials */}
          <Paper
            elevation={0}
            sx={{
              mt: 3,
              p: 2,
              bgcolor: 'success.50',
              border: '2px dashed',
              borderColor: 'success.200',
            }}
          >
            <Typography variant="subtitle2" gutterBottom fontWeight="bold">
              🧪 Demo Credentials
            </Typography>
            <Typography variant="body2" component="div">
              <strong>Email:</strong> admin@example.com<br />
              <strong>Password:</strong> qwerty
            </Typography>
            <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
              ✅ Verified working credentials
            </Typography>
          </Paper>
        </Paper>
      </Container>
    </Box>
  );
};

export default Login;
