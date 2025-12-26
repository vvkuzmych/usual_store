import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Alert,
  Typography,
  CircularProgress,
  Divider,
  Paper
} from '@mui/material';
import {
  PersonAdd as PersonAddIcon,
  ContentCopy as ContentCopyIcon
} from '@mui/icons-material';
import axios from 'axios';

const API_URL = process.env.REACT_APP_SUPPORT_API_URL || 'http://localhost:4001';
const MESSAGING_API_URL = process.env.REACT_APP_MESSAGING_API_URL || 'http://localhost:4001'; // Backend API publishes to Kafka

function CreateSupporterAccount({ open, onClose, onSuccess }) {
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    password: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [createdAccount, setCreatedAccount] = useState(null);

  const generatePassword = () => {
    const length = 12;
    const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*';
    let password = '';
    for (let i = 0; i < length; i++) {
      password += charset.charAt(Math.floor(Math.random() * charset.length));
    }
    setFormData({ ...formData, password });
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      // Create user account via main backend
      const response = await axios.post('http://localhost:4001/api/users', {
        first_name: formData.firstName,
        last_name: formData.lastName,
        email: formData.email,
        password: formData.password,
      });

      const newUser = response.data;
      setCreatedAccount({
        ...formData,
        id: newUser.id,
      });

      // Send credentials via email using messaging service
      try {
        const roleDisplayName = 'Support Agent';
        const loginUrl = 'http://localhost:3005/support/dashboard';

        const emailMessage = `
══════════════════════════════════════════════════════════════
                USUAL STORE - ACCOUNT REGISTRATION
══════════════════════════════════════════════════════════════

Dear ${formData.firstName} ${formData.lastName},

Your account has been successfully created! Below are your credentials
and account information:

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                    ACCOUNT INFORMATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Full Name:  ${formData.firstName} ${formData.lastName}
Email:      ${formData.email}
Password:   ${formData.password}
Role:       ${roleDisplayName}

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                      ACCESS DETAILS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

You have been registered with the role of: ${roleDisplayName}

Your role grants you access to the Support Dashboard where you can:
  • View and respond to support tickets
  • Assist customers with their inquiries
  • Provide technical support

Dashboard URL: ${loginUrl}

Please bookmark this URL for easy access.

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
          to: formData.email,
          subject: 'Welcome to Usual Store - Your Support Agent Account Has Been Created',
          message: emailMessage.trim(),
        });
      } catch (emailError) {
        console.error('Failed to send email:', emailError);
        // Don't fail the whole process if email fails
      }

      setSuccess(true);
      if (onSuccess) {
        onSuccess(newUser);
      }
    } catch (err) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to create account';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setFormData({ firstName: '', lastName: '', email: '', password: '' });
    setError('');
    setSuccess(false);
    setCreatedAccount(null);
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <PersonAddIcon />
          <span>Create Supporter Account</span>
        </Box>
      </DialogTitle>

      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {success && createdAccount ? (
          <Box>
            <Alert severity="success" sx={{ mb: 3 }}>
              Account created successfully! Credentials have been sent to {createdAccount.email}
            </Alert>

            <Paper variant="outlined" sx={{ p: 2, bgcolor: 'grey.50' }}>
              <Typography variant="subtitle2" gutterBottom>
                Account Credentials (Save these!)
              </Typography>
              <Divider sx={{ my: 1 }} />
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="caption" color="text.secondary">
                  Name
                </Typography>
                <Typography variant="body1">
                  {createdAccount.firstName} {createdAccount.lastName}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="caption" color="text.secondary">
                  Email (Login)
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Typography variant="body1">{createdAccount.email}</Typography>
                  <Button
                    size="small"
                    startIcon={<ContentCopyIcon />}
                    onClick={() => copyToClipboard(createdAccount.email)}
                  >
                    Copy
                  </Button>
                </Box>
              </Box>

              <Box>
                <Typography variant="caption" color="text.secondary">
                  Password
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Typography variant="body1" sx={{ fontFamily: 'monospace' }}>
                    {createdAccount.password}
                  </Typography>
                  <Button
                    size="small"
                    startIcon={<ContentCopyIcon />}
                    onClick={() => copyToClipboard(createdAccount.password)}
                  >
                    Copy
                  </Button>
                </Box>
              </Box>
            </Paper>

            <Alert severity="info" sx={{ mt: 2 }}>
              The credentials have also been sent to the user's email address.
            </Alert>
          </Box>
        ) : (
          <form onSubmit={handleSubmit}>
            <TextField
              fullWidth
              label="First Name"
              name="firstName"
              value={formData.firstName}
              onChange={handleChange}
              margin="normal"
              required
              disabled={loading}
            />

            <TextField
              fullWidth
              label="Last Name"
              name="lastName"
              value={formData.lastName}
              onChange={handleChange}
              margin="normal"
              required
              disabled={loading}
            />

            <TextField
              fullWidth
              label="Email Address"
              name="email"
              type="email"
              value={formData.email}
              onChange={handleChange}
              margin="normal"
              required
              disabled={loading}
            />

            <Box sx={{ mt: 2, mb: 1 }}>
              <TextField
                fullWidth
                label="Password"
                name="password"
                type="text"
                value={formData.password}
                onChange={handleChange}
                required
                disabled={loading}
                helperText="Password will be sent to the user via email"
              />
              <Button
                size="small"
                onClick={generatePassword}
                disabled={loading}
                sx={{ mt: 1 }}
              >
                Generate Secure Password
              </Button>
            </Box>
          </form>
        )}
      </DialogContent>

      <DialogActions>
        {success ? (
          <Button onClick={handleClose} variant="contained">
            Done
          </Button>
        ) : (
          <>
            <Button onClick={handleClose} disabled={loading}>
              Cancel
            </Button>
            <Button
              onClick={handleSubmit}
              variant="contained"
              disabled={loading || !formData.email || !formData.password}
            >
              {loading ? (
                <>
                  <CircularProgress size={20} sx={{ mr: 1 }} />
                  Creating...
                </>
              ) : (
                'Create Account'
              )}
            </Button>
          </>
        )}
      </DialogActions>
    </Dialog>
  );
}

export default CreateSupporterAccount;

