import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import Container from "@mui/material/Container";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import CardMedia from "@mui/material/CardMedia";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import Chip from "@mui/material/Chip";
import CircularProgress from "@mui/material/CircularProgress";
import Alert from "@mui/material/Alert";
import { apiService } from "../services/api";
import { useAuth } from "../contexts/AuthContext";
import type { Product } from "../types";

const ProductDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProduct = async () => {
      if (!id) return;
      try {
        const data = await apiService.getProduct(parseInt(id));
        setProduct(data);
      } catch (err) {
        setError("Failed to load product. Please try again later.");
        console.error("Error fetching product:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchProduct();
  }, [id]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="50vh">
        <CircularProgress size={60} />
      </Box>
    );
  }

  if (error || !product) {
    return (
      <Container maxWidth="lg">
        <Alert severity="error">{error || "Product not found"}</Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg">
      <Button
        onClick={() => navigate("/products")}
        sx={{ mb: 3 }}
      >
        ← Back to Products
      </Button>

      <Typography variant="h1" gutterBottom>
        {product.is_recurring ? "Subscribe to " : "Buy "}{product.name}
      </Typography>

      <Grid container spacing={4}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardMedia
              component="img"
              image={product.image || "/placeholder.jpg"}
              alt={product.name}
              sx={{ height: 400, objectFit: "cover" }}
            />
          </Card>
          <Box sx={{ mt: 2 }}>
            <Typography variant="h4" gutterBottom>
              {product.name}
            </Typography>
            <Typography variant="h5" color="primary" gutterBottom>
              ${product.price.toFixed(2)}
              {product.is_recurring && " per month"}
            </Typography>
            <Typography variant="body1" paragraph>
              {product.description}
            </Typography>
            <Chip
              label={`${product.inventory_level} in stock`}
              color={product.inventory_level > 0 ? "success" : "error"}
              size="small"
            />
          </Box>
        </Grid>

        <Grid item xs={12} md={6}>
          {!isAuthenticated ? (
            <Card sx={{ p: 3 }}>
              <Typography variant="h5" gutterBottom>
                Login Required
              </Typography>
              <Typography variant="body1" paragraph>
                You must be logged in to purchase this product.
              </Typography>
              <Button
                variant="contained"
                fullWidth
                onClick={() => navigate("/login")}
              >
                Login to Continue
              </Button>
            </Card>
          ) : (
            <Card sx={{ p: 3 }}>
              <Typography variant="h5" gutterBottom>
                Payment Information
              </Typography>
              <Box component="form" sx={{ mt: 2 }}>
                <TextField
                  label="First Name"
                  required
                  fullWidth
                  margin="normal"
                />
                <TextField
                  label="Last Name"
                  required
                  fullWidth
                  margin="normal"
                />
                <TextField
                  label="Email"
                  type="email"
                  required
                  fullWidth
                  margin="normal"
                />
                <TextField
                  label="Name on Card"
                  required
                  fullWidth
                  margin="normal"
                />
                <Box sx={{ my: 2 }}>
                  <Typography variant="body2" gutterBottom>
                    Credit Card (Stripe Integration)
                  </Typography>
                  <Alert severity="info">
                    Stripe card element would be integrated here
                  </Alert>
                </Box>
                <Button
                  variant="contained"
                  fullWidth
                  size="large"
                  sx={{ mt: 2 }}
                >
                  Charge ${product.price.toFixed(2)}
                </Button>
              </Box>
            </Card>
          )}
        </Grid>
      </Grid>
    </Container>
  );
};

export default ProductDetail;

