import { useEffect, useState } from 'react';
import { loadStripe } from '@stripe/stripe-js';
import { Elements, CardElement, useStripe, useElements } from '@stripe/react-stripe-js';
import { Box, CircularProgress } from '@mui/material';

// Load Stripe publishable key
const stripePromise = loadStripe('pk_test_51OFVZcGnbFbfXcZLp5A9sLhfv7UYB3kjGXXfqQROl8pxEHwm3FXEn8kPCKxRvDnNX0b8bBg2XZgwOILIxZQOUUDv00OVvF3Qo1');

const CARD_ELEMENT_OPTIONS = {
  style: {
    base: {
      fontSize: '16px',
      color: '#424770',
      '::placeholder': {
        color: '#aab7c4',
      },
    },
    invalid: {
      color: '#9e2146',
    },
  },
  hidePostalCode: true,
};

const CardElementWrapper = ({ onCardReady }) => {
  const stripe = useStripe();
  const elements = useElements();
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    if (stripe && elements && !isReady) {
      const cardElement = elements.getElement(CardElement);
      if (cardElement && onCardReady) {
        setIsReady(true);
        onCardReady(cardElement, stripe);
      }
    }
  }, [stripe, elements, onCardReady, isReady]);

  if (!stripe || !elements) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 2 }}>
        <CircularProgress size={24} />
      </Box>
    );
  }

  return (
    <Box
      sx={{
        border: '1px solid #ccc',
        borderRadius: '8px',
        padding: '12px',
        '&:hover': {
          borderColor: '#6a11cb',
        },
        '&:focus-within': {
          borderColor: '#6a11cb',
          borderWidth: '2px',
        },
      }}
    >
      <CardElement options={CARD_ELEMENT_OPTIONS} />
    </Box>
  );
};

const StripeCardElement = ({ onCardReady }) => {
  return (
    <Elements stripe={stripePromise}>
      <CardElementWrapper onCardReady={onCardReady} />
    </Elements>
  );
};

export default StripeCardElement;

