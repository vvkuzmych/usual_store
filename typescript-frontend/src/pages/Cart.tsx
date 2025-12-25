import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import { useNavigate } from "react-router-dom";

const Cart = () => {
  const navigate = useNavigate();

  return (
    <Container maxWidth="lg">
      <Box sx={{ textAlign: "center", py: 8 }}>
        <Typography variant="h1" gutterBottom>
          Shopping Cart
        </Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          Your cart is currently empty.
        </Typography>
        <Button
          variant="contained"
          size="large"
          onClick={() => navigate("/products")}
        >
          Continue Shopping
        </Button>
      </Box>
    </Container>
  );
};

export default Cart;

