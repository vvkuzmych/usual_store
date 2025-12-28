import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material';
import { AuthProvider, useAuth } from './context/AuthContext';
import SupportChatWidget from './components/SupportChatWidget';
import SupporterDashboard from './components/SupporterDashboard';
import UserManagement from './components/UserManagement';
import Login from './components/Login';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
});

// Protected Route Component
function ProtectedRoute({ children }) {
  const { isAuthenticated, loading } = useAuth();
  
  if (loading) {
    return <div>Loading...</div>;
  }
  
  return isAuthenticated ? children : <Navigate to="/support/login" replace />;
}

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <AuthProvider>
          <Routes>
            {/* User-facing chat widget (public) */}
            <Route path="/support" element={<SupportChatWidget />} />
            
            {/* Login page */}
            <Route path="/support/login" element={<Login />} />
            
            {/* Supporter dashboard (public - accessible from admin panel back button) */}
            <Route 
              path="/support/dashboard" 
              element={<SupporterDashboard />} 
            />
            
            {/* User management (protected - super_admin only) */}
            <Route 
              path="/support/users" 
              element={
                <ProtectedRoute>
                  <UserManagement />
                </ProtectedRoute>
              } 
            />
            
            {/* User management via admin path (for API gateway routing) - PUBLIC */}
            <Route 
              path="/admin/users" 
              element={<UserManagement />} 
            />
            
            {/* Default redirect */}
            <Route path="/" element={<Navigate to="/support/login" replace />} />
          </Routes>
        </AuthProvider>
      </Router>
    </ThemeProvider>
  );
}

export default App;

