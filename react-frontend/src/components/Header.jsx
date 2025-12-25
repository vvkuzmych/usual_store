import React, { useState, useEffect } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  Container,
  IconButton,
  Avatar,
  Menu,
  MenuItem,
  Badge,
} from '@mui/material';
import {
  ShoppingCart as ShoppingCartIcon,
  Person as PersonIcon,
  Logout as LogoutIcon,
} from '@mui/icons-material';

const Header = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const [anchorEl, setAnchorEl] = useState(null);
  const [cartItemsCount, setCartItemsCount] = useState(0);

  // Load cart count from localStorage
  const updateCartCount = () => {
    const existingCart = localStorage.getItem('cart');
    if (existingCart) {
      const cart = JSON.parse(existingCart);
      const totalItems = cart.reduce((sum, item) => sum + item.quantity, 0);
      setCartItemsCount(totalItems);
    } else {
      setCartItemsCount(0);
    }
  };

  useEffect(() => {
    // Initial load
    updateCartCount();

    // Listen for storage events (cart updates from other tabs/windows)
    window.addEventListener('storage', updateCartCount);

    // Listen for custom cart update event
    window.addEventListener('cartUpdated', updateCartCount);

    return () => {
      window.removeEventListener('storage', updateCartCount);
      window.removeEventListener('cartUpdated', updateCartCount);
    };
  }, []);

  const handleMenu = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    handleClose();
    logout();
  };

  return (
    <AppBar 
      position="static"
      sx={{
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        boxShadow: '0 2px 10px rgba(0, 0, 0, 0.1)',
      }}
    >
      <Container maxWidth="lg">
        <Toolbar disableGutters sx={{ justifyContent: 'space-between' }}>
          {/* Logo */}
          <Typography
            variant="h6"
            component={RouterLink}
            to="/"
            sx={{
              textDecoration: 'none',
              color: 'white',
              fontWeight: 700,
              display: 'flex',
              alignItems: 'center',
              gap: 1,
              '&:hover': {
                opacity: 0.9,
              },
            }}
          >
            🛍️ Usual Store
          </Typography>

          {/* Navigation */}
          <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
            <Button
              component={RouterLink}
              to="/"
              sx={{
                color: 'white',
                '&:hover': {
                  backgroundColor: 'rgba(255, 255, 255, 0.1)',
                },
              }}
            >
              Home
            </Button>

            <Button
              component={RouterLink}
              to="/products"
              sx={{
                color: 'white',
                '&:hover': {
                  backgroundColor: 'rgba(255, 255, 255, 0.1)',
                },
              }}
            >
              Products
            </Button>

            <IconButton
              component={RouterLink}
              to="/cart"
              sx={{
                color: 'white',
                '&:hover': {
                  backgroundColor: 'rgba(255, 255, 255, 0.1)',
                },
              }}
            >
              <Badge 
                badgeContent={cartItemsCount} 
                color="error"
                sx={{
                  '& .MuiBadge-badge': {
                    fontWeight: 'bold',
                    fontSize: '0.75rem',
                  },
                }}
              >
                <ShoppingCartIcon />
              </Badge>
            </IconButton>

            {isAuthenticated() ? (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <IconButton
                  onClick={handleMenu}
                  sx={{
                    padding: 0.5,
                    '&:hover': {
                      backgroundColor: 'rgba(255, 255, 255, 0.1)',
                    },
                  }}
                >
                  <Avatar
                    sx={{
                      bgcolor: 'rgba(255, 255, 255, 0.2)',
                      width: 36,
                      height: 36,
                    }}
                  >
                    {user?.firstName?.[0] || user?.email?.[0] || 'U'}
                  </Avatar>
                </IconButton>

                <Menu
                  anchorEl={anchorEl}
                  open={Boolean(anchorEl)}
                  onClose={handleClose}
                  anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'right',
                  }}
                  transformOrigin={{
                    vertical: 'top',
                    horizontal: 'right',
                  }}
                >
                  <MenuItem disabled>
                    <Typography variant="body2" color="text.secondary">
                      {user?.firstName || user?.email}
                    </Typography>
                  </MenuItem>
                  <MenuItem onClick={handleLogout}>
                    <LogoutIcon fontSize="small" sx={{ mr: 1 }} />
                    Logout
                  </MenuItem>
                </Menu>
              </Box>
            ) : (
              <Button
                component={RouterLink}
                to="/login"
                variant="contained"
                sx={{
                  bgcolor: 'rgba(40, 167, 69, 0.9)',
                  '&:hover': {
                    bgcolor: 'rgba(40, 167, 69, 1)',
                  },
                }}
                startIcon={<PersonIcon />}
              >
                Login
              </Button>
            )}
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default Header;
