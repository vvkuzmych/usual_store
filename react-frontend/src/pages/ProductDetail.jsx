import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getProduct } from '../services/api';
import { useAuth } from '../context/AuthContext';
import StripeCardElement from '../components/StripeCardElement';
import {
  Box,
  Container,
  Typography,
  Button,
  Card,
  CardContent,
  CardMedia,
  TextField,
  Alert,
  CircularProgress,
  Chip,
  Paper,
  Divider,
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  Lock as LockIcon,
  ShoppingCart as ShoppingCartIcon,
  CheckCircle as CheckCircleIcon,
  PersonAdd as PersonAddIcon,
  Login as LoginIcon,
} from '@mui/icons-material';

const ProductDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user, isAuthenticated } = useAuth();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [processing, setProcessing] = useState(false);
  const [cardMessages, setCardMessages] = useState({ type: '', message: '' });
  const [stripe, setStripe] = useState(null);
  const [cardElement, setCardElement] = useState(null);
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    cardholderName: ''
  });
  const [formErrors, setFormErrors] = useState({});

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

  const handleInputChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
    if (formErrors[e.target.name]) {
      setFormErrors({
        ...formErrors,
        [e.target.name]: ''
      });
    }
  };

  const handleCardReady = (card, stripeInstance) => {
    setCardElement(card);
    setStripe(stripeInstance);
  };

  const validateForm = () => {
    const errors = {};
    if (!formData.firstName.trim()) errors.firstName = 'First name is required';
    if (!formData.lastName.trim()) errors.lastName = 'Last name is required';
    if (!formData.email.trim()) errors.email = 'Email is required';
    else if (!/\S+@\S+\.\S+/.test(formData.email)) errors.email = 'Email is invalid';
    if (!formData.cardholderName.trim()) errors.cardholderName = 'Cardholder name is required';
    
    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const showCardError = (message) => {
    setCardMessages({ type: 'error', message });
  };

  const showCardSuccess = () => {
    setCardMessages({ type: 'success', message: 'Transaction successful!' });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    if (!stripe || !cardElement) {
      showCardError('Stripe is not loaded yet. Please wait a moment.');
      return;
    }

    setProcessing(true);
    setCardMessages({ type: '', message: '' });

    try {
      const { error: stripeError, paymentMethod } = await stripe.createPaymentMethod({
        type: 'card',
        card: cardElement,
        billing_details: {
          name: formData.cardholderName,
          email: formData.email,
        },
      });

      if (stripeError) {
        showCardError(stripeError.message);
        setProcessing(false);
        return;
      }

      const payload = {
        product_id: product.id,
        amount: product.price,
        payment_method: paymentMethod.id,
        card_brand: paymentMethod.card.brand,
        exp_month: paymentMethod.card.exp_month,
        exp_year: paymentMethod.card.exp_year,
        last_four: paymentMethod.card.last4,
        email: formData.email,
        first_name: formData.firstName,
        last_name: formData.lastName,
        plan: product.plan_id || '',
      };

      const endpoint = product.is_recurring 
        ? '/api/create-customer-and-subscribe-to-plan'
        : '/api/payment-intent';

      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      const data = await response.json();

      if (data.error === false || response.ok) {
        showCardSuccess();
        sessionStorage.setItem('first_name', formData.firstName);
        sessionStorage.setItem('last_name', formData.lastName);
        sessionStorage.setItem('amount', formatCurrency(product.price));
        sessionStorage.setItem('last_four', paymentMethod.card.last4);
        
        setTimeout(() => {
          if (product.is_recurring) {
            navigate('/receipt/golden');
          } else {
            navigate('/receipt');
          }
        }, 1000);
      } else {
        showCardError(data.message || 'Payment failed. Please try again.');
        setProcessing(false);
      }
    } catch (err) {
      console.error('Payment error:', err);
      showCardError('An error occurred processing your payment. Please try again.');
      setProcessing(false);
    }
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
        </CardContent>
      </Card>

      {/* Login Prompt for Non-Authenticated Users */}
      {!isAuthenticated() ? (
        <Paper
          elevation={3}
          sx={{
            p: 4,
            textAlign: 'center',
            background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)',
          }}
        >
          <LockIcon sx={{ fontSize: 64, color: 'primary.main', mb: 2 }} />
          <Typography variant="h4" gutterBottom>
            Login Required
          </Typography>
          <Typography variant="body1" color="text.secondary" paragraph>
            You must be logged in to purchase this product.
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            Secure checkout with saved payment methods and order history.
          </Typography>

          <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center', my: 3 }}>
            <Button
              variant="contained"
              size="large"
              startIcon={<LoginIcon />}
              onClick={() => navigate('/login', { state: { from: `/product/${id}` } })}
            >
              Login to Continue
            </Button>
            <Button
              variant="outlined"
              size="large"
              startIcon={<PersonAddIcon />}
              onClick={() => navigate('/signup', { state: { from: `/product/${id}` } })}
            >
              Create Account
            </Button>
          </Box>

          <Paper sx={{ p: 2, bgcolor: 'background.paper', maxWidth: 400, mx: 'auto' }}>
            <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
              Why create an account?
            </Typography>
            <Typography variant="body2" component="div" textAlign="left">
              • Faster checkout<br />
              • Order history & tracking<br />
              • Save payment methods securely<br />
              • Exclusive member deals
            </Typography>
          </Paper>
        </Paper>
      ) : (
        <>
          {cardMessages.message && (
            <Alert 
              severity={cardMessages.type === 'error' ? 'error' : 'success'}
              sx={{ mb: 3 }}
            >
              {cardMessages.message}
            </Alert>
          )}

          <Paper elevation={3} sx={{ p: 3 }}>
            <Typography variant="h5" gutterBottom>
              Payment Information
            </Typography>

            <Box component="form" onSubmit={handleSubmit} noValidate>
              <TextField
                margin="normal"
                required
                fullWidth
                label="First Name"
                name="firstName"
                value={formData.firstName}
                onChange={handleInputChange}
                error={!!formErrors.firstName}
                helperText={formErrors.firstName}
                disabled={processing}
              />

              <TextField
                margin="normal"
                required
                fullWidth
                label="Last Name"
                name="lastName"
                value={formData.lastName}
                onChange={handleInputChange}
                error={!!formErrors.lastName}
                helperText={formErrors.lastName}
                disabled={processing}
              />

              <TextField
                margin="normal"
                required
                fullWidth
                label="Email"
                name="email"
                type="email"
                value={formData.email}
                onChange={handleInputChange}
                error={!!formErrors.email}
                helperText={formErrors.email}
                disabled={processing}
              />

              <TextField
                margin="normal"
                required
                fullWidth
                label="Name on Card"
                name="cardholderName"
                value={formData.cardholderName}
                onChange={handleInputChange}
                error={!!formErrors.cardholderName}
                helperText={formErrors.cardholderName}
                disabled={processing}
              />

              <Box sx={{ mt: 2, mb: 3 }}>
                <Typography variant="body2" gutterBottom fontWeight="bold">
                  Credit Card
                </Typography>
                <StripeCardElement 
                  onCardReady={handleCardReady}
                  onCardError={() => {}}
                />
              </Box>

              <Divider sx={{ my: 3 }} />

              {!processing ? (
                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  fullWidth
                  startIcon={<ShoppingCartIcon />}
                >
                  {product.is_recurring 
                    ? `Pay ${formatCurrency(product.price)}/month` 
                    : `Charge ${formatCurrency(product.price)}`
                  }
                </Button>
              ) : (
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', py: 2 }}>
                  <CircularProgress size={24} sx={{ mr: 2 }} />
                  <Typography>Processing payment...</Typography>
                </Box>
              )}

              <Typography variant="caption" color="text.secondary" sx={{ display: 'block', textAlign: 'center', mt: 2 }}>
                🔒 Secure payment powered by Stripe • Test card: 4242 4242 4242 4242
              </Typography>
            </Box>
          </Paper>
        </>
      )}
    </Container>
  );
};

export default ProductDetail;
