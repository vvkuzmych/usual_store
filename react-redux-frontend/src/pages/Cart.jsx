import { Container, Typography, Paper, Box, Button, IconButton, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@mui/material';
import { Delete, Add, Remove, ShoppingCart } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../hooks/useRedux';
import { removeFromCart, updateQuantity, clearCart, selectCartItems, selectCartTotal } from '../store/slices/cartSlice';
import { addNotification } from '../store/slices/uiSlice';

function Cart() {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  
  const cartItems = useAppSelector(selectCartItems);
  const totalAmount = useAppSelector(selectCartTotal);

  const handleRemoveItem = (id, name) => {
    dispatch(removeFromCart(id));
    dispatch(addNotification({
      message: `${name} removed from cart`,
      severity: 'info',
    }));
  };

  const handleQuantityChange = (id, currentQuantity, change) => {
    const newQuantity = currentQuantity + change;
    if (newQuantity > 0) {
      dispatch(updateQuantity({ id, quantity: newQuantity }));
    }
  };

  const handleClearCart = () => {
    dispatch(clearCart());
    dispatch(addNotification({
      message: 'Cart cleared',
      severity: 'info',
    }));
  };

  const handleCheckout = () => {
    if (cartItems.length === 0) {
      return;
    }
    // Navigate to the first product for checkout
    // User can purchase items individually
    const firstItem = cartItems[0];
    navigate(`/product/${firstItem.id}`);
    dispatch(addNotification({
      message: 'Proceeding to checkout for ' + firstItem.name,
      severity: 'info',
    }));
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
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1">
          Shopping Cart
        </Typography>
        <Button variant="outlined" color="error" onClick={handleClearCart}>
          Clear Cart
        </Button>
      </Box>

      <TableContainer component={Paper}>
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
                  ${item.price.toFixed(2)}
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
                  ${item.totalPrice.toFixed(2)}
                </TableCell>
                <TableCell align="center">
                  <Box sx={{ display: 'flex', gap: 1, justifyContent: 'center' }}>
                    <Button
                      variant="contained"
                      size="small"
                      startIcon={<ShoppingCart />}
                      onClick={() => navigate(`/product/${item.id}`)}
                    >
                      Buy Now
                    </Button>
                    <IconButton
                      color="error"
                      onClick={() => handleRemoveItem(item.id, item.name)}
                    >
                      <Delete />
                    </IconButton>
                  </Box>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Paper sx={{ mt: 3, p: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h5">Total:</Typography>
          <Typography variant="h4" color="primary">
            ${totalAmount.toFixed(2)}
          </Typography>
        </Box>
        <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
          <Button
            variant="outlined"
            fullWidth
            onClick={() => navigate('/products')}
          >
            Continue Shopping
          </Button>
          <Button
            variant="contained"
            fullWidth
            onClick={handleCheckout}
            disabled={cartItems.length === 0}
          >
            Checkout First Item
          </Button>
        </Box>
        <Typography variant="caption" display="block" textAlign="center" sx={{ mt: 2 }} color="text.secondary">
          💡 Tip: Use "Buy Now" buttons above to checkout individual items
        </Typography>
      </Paper>
    </Container>
  );
}

export default Cart;

