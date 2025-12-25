// Product types
export interface Product {
  id: number;
  name: string;
  description: string;
  inventory_level: number;
  price: number;
  image: string;
  is_recurring: boolean;
  plan_id?: string;
  created_at?: string;
  updated_at?: string;
}

// User types
export interface User {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
}

// Auth types
export interface LoginCredentials {
  email: string;
  password: string;
}

export interface AuthResponse {
  authentication_token: {
    token: string;
    expiry: string;
  };
  user: User;
}

// Payment types
export interface PaymentRequest {
  product_id: number;
  first_name: string;
  last_name: string;
  email: string;
  payment_method: string;
  payment_intent_id: string;
  payment_amount: number;
  payment_currency: string;
}

// AI Assistant types
export interface ChatMessage {
  role: "user" | "assistant";
  content: string;
  timestamp?: string;
}

export interface ChatRequest {
  message: string;
  conversation_id?: string;
  user_id?: number;
}

export interface ChatResponse {
  response: string;
  conversation_id: string;
  suggestions?: string[];
  products?: Product[];
}

export interface FeedbackRequest {
  conversation_id: string;
  message_index: number;
  rating: number;
  comment?: string;
}

export interface AIStats {
  total_conversations: number;
  total_cost: number;
  purchases: number;
  conversion_rate: number;
}

