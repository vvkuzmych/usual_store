import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getProduct } from '../services/api';
import {
  Box,
  Container,
  Typography,
  Button,
  Card,
  CardContent,
  CardMedia,
  Alert,
  CircularProgress,
  Chip,
  Paper,
  Divider,
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  ShoppingCart as ShoppingCartIcon,
} from '@mui/icons-material';

const ProductDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [addToCartSuccess, setAddToCartSuccess] = useState(false);

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        setLoading(true);
        const data = await getProduct(id);
        setProduct(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchProduct();
  }, [id]);

  const handleAddToCart = () => {
    if (!product) return;

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
      });
    }

    // Save to localStorage
    localStorage.setItem('cart', JSON.stringify(cart));
    
    // Dispatch custom event to update cart badge
    window.dispatchEvent(new Event('cartUpdated'));
    
    // Show success message
    setAddToCartSuccess(true);
    setTimeout(() => setAddToCartSuccess(false), 3000);
  };

  const formatCurrency = (cents) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(cents / 100);
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="70vh">
        <CircularProgress size={60} />
      </Box>
    );
  }

  if (error || !product) {
    return (
      <Container maxWidth="md" sx={{ py: 4 }}>
        <Alert 
          severity="error"
          action={
            <Button color="inherit" size="small" onClick={() => navigate('/products')}>
              Back to Products
            </Button>
          }
        >
          {error || 'Product not found'}
        </Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Button
        startIcon={<ArrowBackIcon />}
        onClick={() => navigate('/products')}
        sx={{ mb: 3 }}
      >
        Back to Products
      </Button>

      <Typography variant="h3" component="h1" gutterBottom textAlign="center">
        {product.is_recurring ? 'Subscribe to ' : 'Buy '}
        {product.name}
      </Typography>

      <Divider sx={{ my: 3 }} />

      {/* Product Display */}
      <Card sx={{ mb: 4 }}>
        {product.image && (
          <CardMedia
            component="img"
            height="400"
            image={product.image || '/static/mac-mini.png'}
            alt={product.name}
            sx={{ objectFit: 'contain', bgcolor: 'grey.100' }}
            onError={(e) => {
              e.target.src = '/static/mac-mini.png';
            }}
          />
        )}
        <CardContent>
          <Box sx={{ textAlign: 'center' }}>
            <Typography variant="h4" gutterBottom fontWeight="bold">
              {product.name}
            </Typography>
            <Typography variant="h5" color="primary.main" gutterBottom>
              {formatCurrency(product.price)}
              {product.is_recurring && (
                <Chip
                  label="/month"
                  color="warning"
                  size="small"
                  sx={{ ml: 1, fontWeight: 600 }}
                />
              )}
            </Typography>
            <Typography variant="body1" color="text.secondary" paragraph>
              {product.description}
            </Typography>
            {product.inventory_level !== undefined && (
              <Chip
                label={`${product.inventory_level} in stock`}
                color="success"
                size="small"
              />
            )}
          </Box>

          {/* Add to Cart Button */}
          <Box sx={{ mt: 3, display: 'flex', gap: 2, justifyContent: 'center' }}>
            <Button
              variant="contained"
              size="large"
              startIcon={<ShoppingCartIcon />}
              onClick={handleAddToCart}
              sx={{ minWidth: 200 }}
            >
              Add to Cart
            </Button>
          </Box>
          
          {addToCartSuccess && (
            <Alert severity="success" sx={{ mt: 2 }}>
              Product added to cart! <Button size="small" onClick={() => navigate('/cart')}>View Cart</Button>
            </Alert>
          )}
        </CardContent>
      </Card>

      {/* Optional: Show message to go to cart for checkout */}
      <Paper elevation={3} sx={{ p: 4, mt: 3, textAlign: 'center' }}>
        <Typography variant="h5" gutterBottom>
          Ready to checkout?
        </Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          Add items to your cart and proceed to checkout when you're ready.
        </Typography>
        <Button
          variant="contained"
          size="large"
          onClick={() => navigate('/cart')}
          sx={{ mt: 2 }}
        >
          Go to Cart
        </Button>
      </Paper>
    </Container>
  );
};

export default ProductDetail;
