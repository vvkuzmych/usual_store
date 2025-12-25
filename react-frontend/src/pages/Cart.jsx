import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import StripeCardElement from '../components/StripeCardElement';
import {
  Container,
  Typography,
  Paper,
  Box,
  Button,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Alert,
  Snackbar,
  TextField,
  Divider,
  CircularProgress
} from '@mui/material';
import {
  Delete,
  Add,
  Remove,
  ShoppingCart as ShoppingCartIcon,
  Lock as LockIcon,
  Login as LoginIcon
} from '@mui/icons-material';

const Cart = () => {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [cartItems, setCartItems] = useState([]);
  const [notification, setNotification] = useState({ open: false, message: '', severity: 'info' });
  const [processing, setProcessing] = useState(false);
  const [stripe, setStripe] = useState(null);
  const [cardElement, setCardElement] = useState(null);
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    cardholderName: ''
  });
  const [formErrors, setFormErrors] = useState({});
  const [paymentMessage, setPaymentMessage] = useState({ type: '', message: '' });

  useEffect(() => {
    // Load cart from localStorage
    loadCart();
  }, []);

  const loadCart = () => {
    const savedCart = localStorage.getItem('cart');
    if (savedCart) {
      try {
        const items = JSON.parse(savedCart);
        setCartItems(items);
      } catch (error) {
        console.error('Error loading cart:', error);
        setCartItems([]);
      }
    }
  };

  const saveCart = (items) => {
    localStorage.setItem('cart', JSON.stringify(items));
    setCartItems(items);
    // Dispatch custom event to update cart badge
    window.dispatchEvent(new Event('cartUpdated'));
  };

  const handleRemoveItem = (id, name) => {
    const updatedCart = cartItems.filter(item => item.id !== id);
    saveCart(updatedCart);
    showNotification(`${name} removed from cart`, 'info');
  };

  const handleQuantityChange = (id, currentQuantity, change) => {
    const newQuantity = currentQuantity + change;
    if (newQuantity > 0) {
      const updatedCart = cartItems.map(item =>
        item.id === id
          ? { ...item, quantity: newQuantity, totalPrice: item.price * newQuantity }
          : item
      );
      saveCart(updatedCart);
    }
  };

  const handleClearCart = () => {
    saveCart([]);
    showNotification('Cart cleared', 'info');
  };

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

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!validateForm()) {
      setPaymentMessage({ type: 'error', message: 'Please fill in all required fields' });
      return;
    }

    if (!stripe || !cardElement) {
      setPaymentMessage({ type: 'error', message: 'Stripe is not loaded. Please refresh the page.' });
      return;
    }

    if (cartItems.length === 0) {
      setPaymentMessage({ type: 'error', message: 'Your cart is empty' });
      return;
    }

    setProcessing(true);
    setPaymentMessage({ type: '', message: '' });

    try {
      // Create payment method with Stripe
      const { error, paymentMethod } = await stripe.createPaymentMethod({
        type: 'card',
        card: cardElement,
        billing_details: {
          name: formData.cardholderName,
          email: formData.email,
        },
      });

      if (error) {
        setPaymentMessage({ type: 'error', message: error.message });
        setProcessing(false);
        return;
      }

      // Process first item in cart (you can modify this to handle multiple items)
      const firstItem = cartItems[0];
      
      // Determine if it's a subscription or one-time payment
      if (firstItem.is_recurring) {
        // Subscription payment
        const response = await fetch('/api/create-customer-and-subscribe-to-plan', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            first_name: formData.firstName,
            last_name: formData.lastName,
            email: formData.email,
            payment_method: paymentMethod.id,
            plan_id: firstItem.plan_id,
          }),
        });

        const data = await response.json();

        if (!response.ok) {
          throw new Error(data.message || 'Payment failed');
        }

        setPaymentMessage({ type: 'success', message: 'Subscription created successfully!' });
        
        // Clear cart after successful payment
        setTimeout(() => {
          saveCart([]);
          showNotification('Payment successful! Redirecting...', 'success');
          navigate('/');
        }, 2000);
      } else {
        // One-time payment
        const response = await fetch('/api/payment-intent', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            amount: firstItem.price,
            currency: 'usd',
            payment_method_id: paymentMethod.id,
            first_name: formData.firstName,
            last_name: formData.lastName,
            email: formData.email,
          }),
        });

        const data = await response.json();

        if (!response.ok) {
          throw new Error(data.message || 'Payment failed');
        }

        setPaymentMessage({ type: 'success', message: 'Payment successful!' });
        
        // Clear cart after successful payment
        setTimeout(() => {
          saveCart([]);
          showNotification('Payment successful! Redirecting...', 'success');
          navigate('/');
        }, 2000);
      }
    } catch (error) {
      console.error('Payment error:', error);
      setPaymentMessage({ 
        type: 'error', 
        message: error.message || 'An error occurred during payment. Please try again.' 
      });
    } finally {
      setProcessing(false);
    }
  };

  const showNotification = (message, severity = 'info') => {
    setNotification({ open: true, message, severity });
  };

  const handleCloseNotification = () => {
    setNotification({ ...notification, open: false });
  };

  const formatCurrency = (cents) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(cents / 100);
  };

  const calculateTotal = () => {
    return cartItems.reduce((total, item) => total + item.totalPrice, 0);
  };

  if (cartItems.length === 0) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4, textAlign: 'center' }}>
        <Typography variant="h4" gutterBottom>
          Your Cart is Empty
        </Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          Start shopping to add items to your cart
        </Typography>
        <Button variant="contained" onClick={() => navigate('/products')}>
          Browse Products
        </Button>
      </Container>
    );
  }

  return (
    <>
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1">
            Shopping Cart
          </Typography>
          <Button variant="outlined" color="error" onClick={handleClearCart}>
            Clear Cart
          </Button>
        </Box>

        {/* Cart Items Table */}
        <TableContainer component={Paper} sx={{ mb: 3 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Product</TableCell>
                <TableCell align="right">Price</TableCell>
                <TableCell align="center">Quantity</TableCell>
                <TableCell align="right">Total</TableCell>
                <TableCell align="center">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {cartItems.map((item) => (
                <TableRow key={item.id}>
                  <TableCell>
                    <Typography variant="body1">{item.name}</Typography>
                  </TableCell>
                  <TableCell align="right">
                    {formatCurrency(item.price)}
                  </TableCell>
                  <TableCell align="center">
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                      <IconButton
                        size="small"
                        onClick={() => handleQuantityChange(item.id, item.quantity, -1)}
                      >
                        <Remove />
                      </IconButton>
                      <Typography sx={{ mx: 2 }}>{item.quantity}</Typography>
                      <IconButton
                        size="small"
                        onClick={() => handleQuantityChange(item.id, item.quantity, 1)}
                      >
                        <Add />
                      </IconButton>
                    </Box>
                  </TableCell>
                  <TableCell align="right">
                    {formatCurrency(item.totalPrice)}
                  </TableCell>
                  <TableCell align="center">
                    <IconButton
                      color="error"
                      onClick={() => handleRemoveItem(item.id, item.name)}
                    >
                      <Delete />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        {/* Cart Total */}
        <Paper sx={{ p: 3, mb: 3 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h5">Total:</Typography>
            <Typography variant="h4" color="primary">
              {formatCurrency(calculateTotal())}
            </Typography>
          </Box>
        </Paper>

        <Divider sx={{ my: 4 }} />

        {/* Checkout Section */}
        {!isAuthenticated() ? (
          <Paper elevation={3} sx={{ p: 4, textAlign: 'center', background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)' }}>
            <LockIcon sx={{ fontSize: 64, color: 'primary.main', mb: 2 }} />
            <Typography variant="h4" gutterBottom>
              Login Required to Checkout
            </Typography>
            <Typography variant="body1" color="text.secondary" paragraph>
              Please login to complete your purchase
            </Typography>
            <Button
              variant="contained"
              size="large"
              startIcon={<LoginIcon />}
              onClick={() => navigate('/login', { state: { from: '/cart' } })}
            >
              Login to Continue
            </Button>
          </Paper>
        ) : (
          <Paper elevation={3} sx={{ p: 3 }}>
            <Typography variant="h5" gutterBottom>
              Checkout - Payment Information
            </Typography>
            <Typography variant="body2" color="text.secondary" paragraph>
              🔒 Your payment information is secure and encrypted
            </Typography>

            {paymentMessage.message && (
              <Alert severity={paymentMessage.type === 'error' ? 'error' : 'success'} sx={{ mb: 3 }}>
                {paymentMessage.message}
              </Alert>
            )}

            <Box component="form" onSubmit={handleSubmit} noValidate>
              <TextField
                fullWidth
                label="First Name"
                name="firstName"
                value={formData.firstName}
                onChange={handleInputChange}
                error={!!formErrors.firstName}
                helperText={formErrors.firstName}
                required
                sx={{ mb: 2 }}
              />

              <TextField
                fullWidth
                label="Last Name"
                name="lastName"
                value={formData.lastName}
                onChange={handleInputChange}
                error={!!formErrors.lastName}
                helperText={formErrors.lastName}
                required
                sx={{ mb: 2 }}
              />

              <TextField
                fullWidth
                label="Email"
                name="email"
                type="email"
                value={formData.email}
                onChange={handleInputChange}
                error={!!formErrors.email}
                helperText={formErrors.email}
                required
                sx={{ mb: 2 }}
              />

              <TextField
                fullWidth
                label="Cardholder Name"
                name="cardholderName"
                value={formData.cardholderName}
                onChange={handleInputChange}
                error={!!formErrors.cardholderName}
                helperText={formErrors.cardholderName}
                required
                sx={{ mb: 2 }}
              />

              <Typography variant="body2" sx={{ mb: 1, fontWeight: 'bold' }}>
                Credit Card Information *
              </Typography>
              <StripeCardElement onCardReady={handleCardReady} />

              <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
                <Button
                  variant="outlined"
                  fullWidth
                  onClick={() => navigate('/products')}
                  disabled={processing}
                >
                  Continue Shopping
                </Button>
                <Button
                  type="submit"
                  variant="contained"
                  fullWidth
                  disabled={processing || cartItems.length === 0}
                  startIcon={processing ? <CircularProgress size={20} /> : <ShoppingCartIcon />}
                >
                  {processing ? 'Processing...' : `Pay ${formatCurrency(calculateTotal())}`}
                </Button>
              </Box>

              <Typography variant="caption" display="block" textAlign="center" sx={{ mt: 2 }} color="text.secondary">
                Test Card: 4242 4242 4242 4242 | Exp: 12/25 | CVC: 123
              </Typography>
            </Box>
          </Paper>
        )}
      </Container>

      <Snackbar
        open={notification.open}
        autoHideDuration={6000}
        onClose={handleCloseNotification}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={handleCloseNotification} severity={notification.severity} sx={{ width: '100%' }}>
          {notification.message}
        </Alert>
      </Snackbar>
    </>
  );
};

export default Cart;
