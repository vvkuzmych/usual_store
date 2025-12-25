import axios from 'axios';

// Use empty string to make API calls relative to the current origin
// This allows Nginx to proxy requests to the backend
const API_BASE_URL = process.env.REACT_APP_API_URL || '';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

// Products API
export const getProducts = async () => {
  const response = await api.get('/api/products');
  return response.data;
};

export const getProduct = async (id) => {
  const response = await api.get(`/api/product/${id}`);
  return response.data;
};

// Cart API
export const addToCart = async (productId, quantity = 1) => {
  const response = await api.post('/api/cart/add', { productId, quantity });
  return response.data;
};

export const getCart = async () => {
  const response = await api.get('/api/cart');
  return response.data;
};

// Auth API
export const login = async (email, password) => {
  const response = await api.post('/api/authenticate', { email, password });
  return response.data;
};

export const signup = async (userData) => {
  const response = await api.post('/api/signup', userData);
  return response.data;
};

export const logout = async () => {
  const response = await api.post('/api/logout');
  return response.data;
};

// Checkout API
export const createCheckoutSession = async (cartData) => {
  const response = await api.post('/api/checkout', cartData);
  return response.data;
};

export default api;

