import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { getProducts } from '../services/api';
import {
  Box,
  Container,
  Typography,
  Grid,
  Card,
  CardContent,
  CardMedia,
  CardActions,
  Button,
  Chip,
  CircularProgress,
  Alert,
} from '@mui/material';
import {
  ShoppingCart as ShoppingCartIcon,
  Autorenew as AutorenewIcon,
  Warning as WarningIcon,
} from '@mui/icons-material';

const Products = () => {
  const navigate = useNavigate();
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadProducts();
  }, []);

  const loadProducts = async () => {
    try {
      const data = await getProducts();
      setProducts(data || []);
    } catch (err) {
      console.error('Error loading products:', err);
      setError('Failed to load products');
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (cents) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(cents / 100);
  };

  const handleProductClick = (productId) => {
    navigate(`/product/${productId}`);
  };

  const handleBuyNow = (e, product) => {
    e.stopPropagation();

    // Get existing cart from localStorage
    const existingCart = localStorage.getItem('cart');
    let cart = existingCart ? JSON.parse(existingCart) : [];

    // Check if product already in cart
    const existingItemIndex = cart.findIndex(item => item.id === product.id);

    if (existingItemIndex >= 0) {
      // Update quantity
      cart[existingItemIndex].quantity += 1;
      cart[existingItemIndex].totalPrice = cart[existingItemIndex].price * cart[existingItemIndex].quantity;
    } else {
      // Add new item
      cart.push({
        id: product.id,
        name: product.name,
        price: product.price,
        quantity: 1,
        totalPrice: product.price,
        image: product.image,
        is_recurring: product.is_recurring,
        plan_id: product.plan_id || '',
      });
    }

    // Save to localStorage
    localStorage.setItem('cart', JSON.stringify(cart));
    
    // Dispatch custom event to update cart badge
    window.dispatchEvent(new Event('cartUpdated'));
    
    // Navigate to cart
    navigate('/cart');
  };

  if (loading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="70vh"
      >
        <CircularProgress size={60} />
      </Box>
    );
  }

  if (error) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error" action={
          <Button color="inherit" size="small" onClick={loadProducts}>
            Try Again
          </Button>
        }>
          {error}
        </Alert>
      </Container>
    );
  }

  if (products.length === 0) {
    return (
      <Container maxWidth="lg" sx={{ py: 8, textAlign: 'center' }}>
        <Typography variant="h4" gutterBottom>
          No Products Available
        </Typography>
        <Typography color="text.secondary">
          Check back soon for new products!
        </Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      {/* Header */}
      <Box sx={{ textAlign: 'center', mb: 6 }}>
        <Typography variant="h2" component="h1" gutterBottom>
          Our Products
        </Typography>
        <Typography variant="h6" color="text.secondary">
          Browse our collection of amazing widgets and plans
        </Typography>
      </Box>

      {/* Products Grid */}
      <Grid container spacing={4}>
        {products.map((product) => (
          <Grid item xs={12} sm={6} md={4} key={product.id}>
            <Card
              sx={{
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                cursor: 'pointer',
                position: 'relative',
              }}
              onClick={() => handleProductClick(product.id)}
            >
              {/* Product Image */}
              <CardMedia
                component="img"
                height="240"
                image={product.image || '/static/mac-mini.png'}
                alt={product.name}
                sx={{
                  objectFit: 'cover',
                }}
                onError={(e) => {
                  e.target.src = '/static/mac-mini.png';
                }}
              />

              {/* Badges */}
              <Box 
                sx={{ 
                  position: 'absolute', 
                  top: 12, 
                  left: 12, 
                  right: 12, 
                  display: 'flex', 
                  gap: 1, 
                  flexWrap: 'wrap',
                  zIndex: 1,
                }}
              >
                {product.is_recurring && (
                  <Chip
                    icon={<AutorenewIcon />}
                    label="Subscription"
                    color="warning"
                    size="small"
                    sx={{ 
                      fontWeight: 600,
                      backgroundColor: 'rgba(255, 152, 0, 0.95)',
                      color: 'white',
                      boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
                    }}
                  />
                )}
                {product.inventory_level !== undefined && product.inventory_level < 10 && (
                  <Chip
                    icon={<WarningIcon />}
                    label={`Only ${product.inventory_level} left!`}
                    color="error"
                    size="small"
                    sx={{ 
                      fontWeight: 600,
                      backgroundColor: 'rgba(211, 47, 47, 0.95)',
                      color: 'white',
                      boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
                    }}
                  />
                )}
              </Box>

              {/* Product Info */}
              <CardContent sx={{ flexGrow: 1, pt: 2.5 }}>
                <Typography 
                  gutterBottom 
                  variant="h5" 
                  component="h2" 
                  fontWeight="bold"
                  sx={{ 
                    minHeight: '64px',
                    display: 'flex',
                    alignItems: 'center',
                  }}
                >
                  {product.name}
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  {product.description}
                </Typography>
              </CardContent>

              {/* Footer */}
              <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
                <Box>
                  <Typography variant="h6" color="primary.main" fontWeight="bold">
                    {formatCurrency(product.price)}
                  </Typography>
                  {product.is_recurring && (
                    <Typography variant="caption" color="text.secondary">
                      per month
                    </Typography>
                  )}
                </Box>
                <Button
                  variant="contained"
                  size="small"
                  startIcon={<ShoppingCartIcon />}
                  onClick={(e) => handleBuyNow(e, product)}
                >
                  {product.is_recurring ? 'Subscribe' : 'Buy Now'}
                </Button>
              </CardActions>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Container>
  );
};

export default Products;
