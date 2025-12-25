import Box from "@mui/material/Box";
import Container from "@mui/material/Container";
import Grid from "@mui/material/Grid";
import Typography from "@mui/material/Typography";
import Link from "@mui/material/Link";
import { Link as RouterLink } from "react-router-dom";

const Footer = () => {
  return (
    <Box
      component="footer"
      sx={{
        py: 4,
        px: 2,
        mt: "auto",
        backgroundColor: (theme) =>
          theme.palette.mode === "light"
            ? theme.palette.grey[200]
            : theme.palette.grey[800],
      }}
    >
      <Container maxWidth="lg">
        <Grid container spacing={4}>
          <Grid item xs={12} sm={4}>
            <Typography variant="h6" color="text.primary" gutterBottom>
              🛍️ Usual Store
            </Typography>
            <Typography variant="body2" color="text.secondary">
              TypeScript Edition - Your modern shop for quality products.
            </Typography>
          </Grid>
          <Grid item xs={12} sm={4}>
            <Typography variant="h6" color="text.primary" gutterBottom>
              Quick Links
            </Typography>
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
              <Link component={RouterLink} to="/" color="text.secondary">
                Home
              </Link>
              <Link component={RouterLink} to="/products" color="text.secondary">
                Products
              </Link>
              <Link component={RouterLink} to="/cart" color="text.secondary">
                Cart
              </Link>
              <Link component={RouterLink} to="/login" color="text.secondary">
                Login
              </Link>
            </Box>
          </Grid>
          <Grid item xs={12} sm={4}>
            <Typography variant="h6" color="text.primary" gutterBottom>
              Support
            </Typography>
            <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
              <Link href="#" color="text.secondary">
                Help Center
              </Link>
              <Link href="#" color="text.secondary">
                Contact Us
              </Link>
              <Link href="#" color="text.secondary">
                Privacy Policy
              </Link>
              <Link href="#" color="text.secondary">
                Terms of Service
              </Link>
            </Box>
          </Grid>
        </Grid>
        <Box sx={{ mt: 3, textAlign: "center" }}>
          <Typography variant="body2" color="text.secondary">
            © 2025 Usual Store (TypeScript Edition). All rights reserved.
          </Typography>
        </Box>
      </Container>
    </Box>
  );
};

export default Footer;

