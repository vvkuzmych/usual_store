import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import {
  Box,
  Container,
  Typography,
  Link,
  Grid,
  Divider,
} from '@mui/material';
import {
  Favorite as FavoriteIcon,
} from '@mui/icons-material';

const Footer = () => {
  return (
    <Box
      component="footer"
      sx={{
        bgcolor: 'grey.900',
        color: 'white',
        py: 6,
        mt: 'auto',
      }}
    >
      <Container maxWidth="lg">
        <Grid container spacing={4}>
          {/* Company Info */}
          <Grid item xs={12} sm={4}>
            <Typography variant="h6" gutterBottom fontWeight="bold">
              🛍️ Usual Store
            </Typography>
            <Typography variant="body2" color="grey.400">
              Your one-stop shop for quality products and subscriptions.
            </Typography>
          </Grid>

          {/* Quick Links */}
          <Grid item xs={12} sm={4}>
            <Typography variant="h6" gutterBottom fontWeight="bold">
              Quick Links
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
              <Link
                component={RouterLink}
                to="/"
                color="grey.400"
                underline="hover"
              >
                Home
              </Link>
              <Link
                component={RouterLink}
                to="/products"
                color="grey.400"
                underline="hover"
              >
                Products
              </Link>
              <Link
                component={RouterLink}
                to="/cart"
                color="grey.400"
                underline="hover"
              >
                Cart
              </Link>
              <Link
                component={RouterLink}
                to="/login"
                color="grey.400"
                underline="hover"
              >
                Login
              </Link>
            </Box>
          </Grid>

          {/* Support */}
          <Grid item xs={12} sm={4}>
            <Typography variant="h6" gutterBottom fontWeight="bold">
              Support
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
              <Link href="#" color="grey.400" underline="hover">
                Help Center
              </Link>
              <Link href="#" color="grey.400" underline="hover">
                Contact Us
              </Link>
              <Link href="#" color="grey.400" underline="hover">
                Privacy Policy
              </Link>
              <Link href="#" color="grey.400" underline="hover">
                Terms of Service
              </Link>
            </Box>
          </Grid>
        </Grid>

        <Divider sx={{ my: 4, bgcolor: 'grey.800' }} />

        <Box sx={{ textAlign: 'center' }}>
          <Typography variant="body2" color="grey.400">
            Made with <FavoriteIcon sx={{ fontSize: 14, color: 'error.main', verticalAlign: 'middle' }} /> by Usual Store Team
          </Typography>
          <Typography variant="body2" color="grey.400" sx={{ mt: 1 }}>
            © {new Date().getFullYear()} Usual Store. All rights reserved.
          </Typography>
        </Box>
      </Container>
    </Box>
  );
};

export default Footer;
