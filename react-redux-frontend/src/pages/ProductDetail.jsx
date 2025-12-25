import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Container,
  Typography,
  Button,
  Box,
  CircularProgress,
  Grid,
  Paper,
  TextField,
  Alert,
  Card,
  CardContent,
  CardMedia,
  Chip,
  Divider,
} from '@mui/material';
import {
  ArrowBack,
  Lock as LockIcon,
  Login as LoginIcon,
  PersonAdd as PersonAddIcon,
  CheckCircle as CheckCircleIcon,
} from '@mui/icons-material';
import { useAppDispatch, useAppSelector } from '../hooks/useRedux';
import { fetchProductById } from '../store/slices/productsSlice';
import { addNotification } from '../store/slices/uiSlice';
import StripeCardElement from '../components/StripeCardElement';

function ProductDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  
  const { selectedProduct: product, loading, error } = useAppSelector((state) => state.products);
  const { isAuthenticated, user } = useAppSelector((state) => state.auth);
  
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
    dispatch(fetchProductById(id));
  }, [dispatch, id]);

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

  const formatCurrency = (cents) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(cents / 100);
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
        
        dispatch(addNotification({
          message: 'Payment successful! Redirecting...',
          severity: 'success',
        }));

        setTimeout(() => {
          navigate('/products');
        }, 1500);
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

  if (loading) {
    return (
      <Container sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '60vh' }}>
        <CircularProgress size={60} />
      </Container>
    );
  }

  if (error || !product) {
    return (
      <Container sx={{ mt: 4 }}>
        <Alert severity="error" action={
          <Button color="inherit" size="small" onClick={() => navigate('/products')}>
            Back to Products
          </Button>
        }>
          {error || 'Product not found'}
        </Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Button
        startIcon={<ArrowBack />}
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
      {!isAuthenticated ? (
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
          </Box>

          <Paper sx={{ p: 2, bgcolor: 'background.paper', maxWidth: 400, mx: 'auto' }}>
            <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
              Why login?
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
              icon={cardMessages.type === 'success' ? <CheckCircleIcon /> : undefined}
            >
              {cardMessages.message}
            </Alert>
          )}

          <Paper elevation={3} sx={{ p: 3 }}>
            <Typography variant="h5" gutterBottom>
              Payment Information
            </Typography>

            <Box component="form" onSubmit={handleSubmit} noValidate>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <TextField
                    required
                    fullWidth
                    label="First Name"
                    name="firstName"
                    value={formData.firstName}
                    onChange={handleInputChange}
                    error={Boolean(formErrors.firstName)}
                    helperText={formErrors.firstName}
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <TextField
                    required
                    fullWidth
                    label="Last Name"
                    name="lastName"
                    value={formData.lastName}
                    onChange={handleInputChange}
                    error={Boolean(formErrors.lastName)}
                    helperText={formErrors.lastName}
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    required
                    fullWidth
                    label="Email"
                    name="email"
                    type="email"
                    value={formData.email}
                    onChange={handleInputChange}
                    error={Boolean(formErrors.email)}
                    helperText={formErrors.email}
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    required
                    fullWidth
                    label="Cardholder Name"
                    name="cardholderName"
                    value={formData.cardholderName}
                    onChange={handleInputChange}
                    error={Boolean(formErrors.cardholderName)}
                    helperText={formErrors.cardholderName}
                  />
                </Grid>
                <Grid item xs={12}>
                  <Typography variant="subtitle2" gutterBottom>
                    Card Details
                  </Typography>
                  <StripeCardElement onCardReady={handleCardReady} />
                </Grid>
              </Grid>

              <Box sx={{ mt: 3 }}>
                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  fullWidth
                  disabled={processing || !stripe || !cardElement}
                >
                  {processing ? (
                    <>
                      <CircularProgress size={24} sx={{ mr: 1 }} />
                      Processing...
                    </>
                  ) : (
                    `Pay ${formatCurrency(product.price)}`
                  )}
                </Button>
              </Box>

              <Typography variant="caption" display="block" textAlign="center" sx={{ mt: 2 }}>
                🔒 Your payment information is secure and encrypted
              </Typography>
            </Box>
          </Paper>
        </>
      )}
    </Container>
  );
}

export default ProductDetail;
