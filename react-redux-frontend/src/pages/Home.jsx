import { Container, Typography, Button, Box, Grid, Card, CardContent } from '@mui/material';
import { ShoppingCart, Security, LocalShipping } from '@mui/icons-material';
import { Link } from 'react-router-dom';

function Home() {
  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Box sx={{ textAlign: 'center', mb: 6 }}>
        <Typography variant="h2" component="h1" gutterBottom>
          Welcome to Usual Store
        </Typography>
        <Typography variant="h5" color="text.secondary" paragraph>
          Powered by React + Redux Toolkit
        </Typography>
        <Button
          variant="contained"
          size="large"
          component={Link}
          to="/products"
          sx={{ mt: 2 }}
        >
          Shop Now
        </Button>
      </Box>

      <Grid container spacing={4} sx={{ mt: 4 }}>
        <Grid item xs={12} md={4}>
          <Card sx={{ height: '100%', textAlign: 'center', p: 2 }}>
            <ShoppingCart sx={{ fontSize: 60, color: 'primary.main', mb: 2 }} />
            <CardContent>
              <Typography variant="h5" component="h2" gutterBottom>
                Wide Selection
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Browse through our extensive collection of widgets and subscriptions
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={4}>
          <Card sx={{ height: '100%', textAlign: 'center', p: 2 }}>
            <Security sx={{ fontSize: 60, color: 'primary.main', mb: 2 }} />
            <CardContent>
              <Typography variant="h5" component="h2" gutterBottom>
                Secure Payments
              </Typography>
              <Typography variant="body2" color="text.secondary">
                All transactions are encrypted and secure with Stripe integration
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={4}>
          <Card sx={{ height: '100%', textAlign: 'center', p: 2 }}>
            <LocalShipping sx={{ fontSize: 60, color: 'primary.main', mb: 2 }} />
            <CardContent>
              <Typography variant="h5" component="h2" gutterBottom>
                Fast Delivery
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Quick and reliable shipping to your doorstep
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Container>
  );
}

export default Home;

