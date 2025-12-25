import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import {
  Box,
  Container,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  Paper,
} from '@mui/material';
import {
  ShoppingCart as ShoppingCartIcon,
  Autorenew as AutorenewIcon,
  Security as SecurityIcon,
  LocalShipping as LocalShippingIcon,
} from '@mui/icons-material';

const Home = () => {
  const features = [
    {
      icon: <ShoppingCartIcon sx={{ fontSize: 48 }} />,
      title: 'Wide Selection',
      description: 'Browse through our curated collection of quality products.',
    },
    {
      icon: <AutorenewIcon sx={{ fontSize: 48 }} />,
      title: 'Flexible Subscriptions',
      description: 'Subscribe to your favorite products with flexible plans.',
    },
    {
      icon: <SecurityIcon sx={{ fontSize: 48 }} />,
      title: 'Secure Payments',
      description: 'Shop with confidence using our secure Stripe integration.',
    },
    {
      icon: <LocalShippingIcon sx={{ fontSize: 48 }} />,
      title: 'Fast Shipping',
      description: 'Get your orders delivered quickly and reliably.',
    },
  ];

  return (
    <Box>
      {/* Hero Section */}
      <Box
        sx={{
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          color: 'white',
          py: 12,
          textAlign: 'center',
        }}
      >
        <Container maxWidth="md">
          <Typography variant="h2" component="h1" gutterBottom fontWeight="bold">
            Welcome to Usual Store
          </Typography>
          <Typography variant="h5" paragraph sx={{ mb: 4 }}>
            Your one-stop shop for quality products and subscriptions
          </Typography>
          <Button
            component={RouterLink}
            to="/products"
            variant="contained"
            size="large"
            startIcon={<ShoppingCartIcon />}
            sx={{
              bgcolor: 'white',
              color: 'primary.main',
              px: 4,
              py: 1.5,
              fontSize: '1.1rem',
              '&:hover': {
                bgcolor: 'grey.100',
              },
            }}
          >
            Shop Now
          </Button>
        </Container>
      </Box>

      {/* Features Section */}
      <Container maxWidth="lg" sx={{ py: 8 }}>
        <Typography variant="h3" component="h2" textAlign="center" gutterBottom fontWeight="bold">
          Why Choose Us
        </Typography>
        <Typography variant="h6" textAlign="center" color="text.secondary" paragraph sx={{ mb: 6 }}>
          We provide the best shopping experience for our customers
        </Typography>

        <Grid container spacing={4}>
          {features.map((feature, index) => (
            <Grid item xs={12} sm={6} md={3} key={index}>
              <Card
                sx={{
                  height: '100%',
                  textAlign: 'center',
                  transition: 'all 0.3s ease',
                  '&:hover': {
                    transform: 'translateY(-8px)',
                    boxShadow: 6,
                  },
                }}
              >
                <CardContent>
                  <Box sx={{ color: 'primary.main', mb: 2 }}>
                    {feature.icon}
                  </Box>
                  <Typography variant="h6" gutterBottom fontWeight="bold">
                    {feature.title}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {feature.description}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Container>

      {/* CTA Section */}
      <Box sx={{ bgcolor: 'grey.100', py: 8 }}>
        <Container maxWidth="md">
          <Paper
            elevation={0}
            sx={{
              p: 6,
              textAlign: 'center',
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              color: 'white',
            }}
          >
            <Typography variant="h4" gutterBottom fontWeight="bold">
              Ready to Start Shopping?
            </Typography>
            <Typography variant="h6" paragraph sx={{ mb: 4 }}>
              Explore our products and find exactly what you need
            </Typography>
            <Button
              component={RouterLink}
              to="/products"
              variant="contained"
              size="large"
              sx={{
                bgcolor: 'white',
                color: 'primary.main',
                px: 4,
                py: 1.5,
                '&:hover': {
                  bgcolor: 'grey.100',
                },
              }}
            >
              Browse Products
            </Button>
          </Paper>
        </Container>
      </Box>
    </Box>
  );
};

export default Home;
