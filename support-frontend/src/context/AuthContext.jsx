import React, { createContext, useState, useContext, useEffect } from 'react';
import axios from 'axios';

const AuthContext = createContext(null);

// Use the main backend authentication endpoint
const AUTH_API_URL = process.env.REACT_APP_AUTH_API_URL || 'http://localhost:4001';

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Check if user is already logged in on mount
  useEffect(() => {
    const storedUser = localStorage.getItem('support_user');
    if (storedUser) {
      try {
        setUser(JSON.parse(storedUser));
      } catch (e) {
        localStorage.removeItem('support_user');
      }
    }
    setLoading(false);
  }, []);

  const login = async (email, password) => {
    setError(null);
    setLoading(true);
    
    try {
      // Authenticate with main backend
      const response = await axios.post(`${AUTH_API_URL}/api/authenticate`, {
        email,
        password,
      });

      const userData = response.data;
      
      // Check if user has super_admin, admin, or supporter role
      const allowedRoles = ['super_admin', 'admin', 'supporter'];
      if (!allowedRoles.includes(userData.role)) {
        throw new Error('Access denied. Only administrators and support staff can access this dashboard.');
      }

      const userInfo = {
        id: userData.id,
        email: userData.email || email,
        firstName: userData.first_name || 'User',
        lastName: userData.last_name || '',
        role: userData.role,
      };

      setUser(userInfo);
      localStorage.setItem('support_user', JSON.stringify(userInfo));
      setLoading(false);
      return { success: true };
    } catch (err) {
      const errorMessage = err.response?.data?.error || err.message || 'Authentication failed';
      setError(errorMessage);
      setLoading(false);
      return { success: false, error: errorMessage };
    }
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem('support_user');
  };

  const isAdmin = () => {
    return user?.role === 'admin' || user?.role === 'super_admin';
  };

  const isSuperAdmin = () => {
    return user?.role === 'super_admin';
  };

  const value = {
    user,
    loading,
    error,
    login,
    logout,
    isAdmin,
    isSuperAdmin,
    isAuthenticated: !!user,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export default AuthContext;

