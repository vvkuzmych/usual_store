# Stripe Payment Integration Setup

## 🔑 Required Configuration

To enable Stripe payments in your React frontend, you need to add your Stripe publishable key to the `.env` file.

### Step 1: Get Your Stripe Keys

1. Go to https://dashboard.stripe.com/test/apikeys
2. Copy your **Publishable key** (starts with `pk_test_`)
3. Copy your **Secret key** (starts with `sk_test_`)

### Step 2: Add Keys to .env File

Add these lines to your `.env` file in the project root:

```bash
# Backend Stripe keys
STRIPE_SECRET=sk_test_your_secret_key_here
STRIPE_KEY=pk_test_your_publishable_key_here

# React Frontend Stripe key (same as STRIPE_KEY)
REACT_APP_STRIPE_PUBLISHABLE_KEY=pk_test_your_publishable_key_here
```

### Step 3: Rebuild Docker Containers

```bash
# Rebuild React frontend with new env var
docker compose build react-frontend

# Restart services
docker compose up -d
```

## 🧪 Test Cards

Use these test card numbers in development:

| Card Number         | Description              |
|---------------------|--------------------------|
| 4242 4242 4242 4242 | Visa - Success           |
| 4000 0000 0000 0002 | Visa - Card declined     |
| 4000 0000 0000 9995 | Visa - Insufficient funds|
| 5555 5555 5555 4444 | Mastercard - Success     |
| 3782 822463 10005   | American Express         |

**Expiry Date:** Any future date (e.g., 12/25)  
**CVC:** Any 3 digits (e.g., 123)  
**ZIP:** Any 5 digits (e.g., 12345)

## 📋 Features Implemented

✅ **Stripe.js Integration**
- Secure card input element
- Real-time card validation
- PCI-compliant (Stripe handles card data)

✅ **Payment Processing**
- Create payment method
- Handle payment intent
- Support for one-time payments
- Support for subscriptions

✅ **Form Validation**
- Client-side validation
- Real-time error messages
- Field-level feedback

✅ **User Experience**
- Loading spinner during payment
- Success/error messages
- Test card hint in footer

## 🔒 Security Notes

- **Never** commit `.env` file to git
- Card data never touches your server
- Stripe handles all PCI compliance
- Use test keys for development
- Use live keys only in production

## 🚀 Production Checklist

Before going live:

- [ ] Replace test keys with live keys
- [ ] Enable production mode in Stripe Dashboard
- [ ] Set up webhook endpoints
- [ ] Test with real payment methods
- [ ] Configure proper error handling
- [ ] Set up payment confirmation emails
- [ ] Test subscription renewals (if using subscriptions)

## 📝 Backend Integration

The React frontend sends payment data to these endpoints:

**One-time payment:**
```
POST /api/payment-intent
```

**Subscription:**
```
POST /api/create-customer-and-subscribe-to-plan
```

Make sure your Go backend has these endpoints configured!

## 🆘 Troubleshooting

### "Stripe is not loaded yet"
- Check that REACT_APP_STRIPE_PUBLISHABLE_KEY is set in .env
- Rebuild the React frontend container
- Check browser console for Stripe.js loading errors

### "Payment failed"
- Verify Stripe keys are correct
- Check backend logs for errors
- Ensure backend endpoints are working
- Try a different test card

### Card element not appearing
- Check browser console for errors
- Verify Stripe.js script loaded (check Network tab)
- Clear cache and hard refresh (Cmd+Shift+R)

## 📚 Resources

- [Stripe React Documentation](https://stripe.com/docs/stripe-js/react)
- [Stripe Testing Guide](https://stripe.com/docs/testing)
- [Stripe API Reference](https://stripe.com/docs/api)

