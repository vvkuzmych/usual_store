import { useNavigate } from "react-router-dom";
import Container from "@mui/material/Container";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import ShoppingBagIcon from "@mui/icons-material/ShoppingBag";
import AutorenewIcon from "@mui/icons-material/Autorenew";
import SecurityIcon from "@mui/icons-material/Security";
import LocalShippingIcon from "@mui/icons-material/LocalShipping";

const Home = () => {
  const navigate = useNavigate();

  const features = [
    {
      icon: <ShoppingBagIcon sx={{ fontSize: 48, color: "primary.main" }} />,
      title: "Wide Selection",
      description: "Browse through our curated collection of quality products.",
    },
    {
      icon: <AutorenewIcon sx={{ fontSize: 48, color: "secondary.main" }} />,
      title: "Flexible Subscriptions",
      description: "Subscribe to your favorite products with flexible plans.",
    },
    {
      icon: <SecurityIcon sx={{ fontSize: 48, color: "success.main" }} />,
      title: "Secure Payments",
      description: "Shop with confidence using our secure Stripe integration.",
    },
    {
      icon: <LocalShippingIcon sx={{ fontSize: 48, color: "warning.main" }} />,
      title: "Fast Shipping",
      description: "Get your orders delivered quickly and reliably.",
    },
  ];

  return (
    <Container maxWidth="lg">
      {/* Hero Section */}
      <Box
        sx={{
          textAlign: "center",
          py: 8,
          background: "linear-gradient(45deg, #2196f3 30%, #21cbf3 90%)",
          borderRadius: 3,
          color: "white",
          mb: 6,
        }}
      >
        <Typography variant="h1" gutterBottom>
          Welcome to Usual Store
        </Typography>
        <Typography variant="h6" paragraph>
          TypeScript Edition - Your one-stop shop for quality products and subscriptions
        </Typography>
        <Button
          variant="contained"
          size="large"
          onClick={() => navigate("/products")}
          sx={{
            bgcolor: "white",
            color: "primary.main",
            "&:hover": { bgcolor: "grey.100" },
            mt: 2,
          }}
        >
          Shop Now
        </Button>
      </Box>

      {/* Features Section */}
      <Box sx={{ mb: 6 }}>
        <Typography variant="h2" textAlign="center" gutterBottom>
          Why Choose Us
        </Typography>
        <Typography
          variant="body1"
          textAlign="center"
          color="text.secondary"
          paragraph
        >
          We provide the best shopping experience for our customers
        </Typography>
        <Grid container spacing={3} sx={{ mt: 2 }}>
          {features.map((feature, index) => (
            <Grid item xs={12} sm={6} md={3} key={index}>
              <Card
                sx={{
                  height: "100%",
                  textAlign: "center",
                  p: 2,
                  "&:hover": {
                    boxShadow: 6,
                  },
                }}
              >
                <CardContent>
                  {feature.icon}
                  <Typography variant="h6" gutterBottom sx={{ mt: 2 }}>
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
      </Box>

      {/* CTA Section */}
      <Box
        sx={{
          textAlign: "center",
          py: 6,
          bgcolor: "background.paper",
          borderRadius: 2,
          mb: 4,
        }}
      >
        <Typography variant="h4" gutterBottom>
          Ready to Start Shopping?
        </Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          Explore our products and find exactly what you need
        </Typography>
        <Button
          variant="contained"
          size="large"
          onClick={() => navigate("/products")}
        >
          Browse Products
        </Button>
      </Box>
    </Container>
  );
};

export default Home;

