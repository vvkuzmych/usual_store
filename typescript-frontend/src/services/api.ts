import axios, { AxiosInstance } from "axios";
import type {
  Product,
  LoginCredentials,
  AuthResponse,
  PaymentRequest,
  ChatRequest,
  ChatResponse,
  FeedbackRequest,
  AIStats,
} from "../types";

const API_BASE_URL = "";

class ApiService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        "Content-Type": "application/json",
      },
    });

    // Add auth token to requests
    this.api.interceptors.request.use((config) => {
      const token = localStorage.getItem("authToken");
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });
  }

  // Product APIs
  async getProducts(): Promise<Product[]> {
    const response = await this.api.get<Product[]>("/api/products");
    return response.data;
  }

  async getProduct(id: number): Promise<Product> {
    const response = await this.api.get<Product>(`/api/product/${id}`);
    return response.data;
  }

  // Auth APIs
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await this.api.post<AuthResponse>(
      "/api/authenticate",
      credentials
    );
    return response.data;
  }

  // Payment APIs
  async createPaymentIntent(data: {
    amount: number;
    currency: string;
  }): Promise<{ client_secret: string }> {
    const response = await this.api.post("/api/payment-intent", data);
    return response.data;
  }

  async processPayment(data: PaymentRequest): Promise<{ success: boolean }> {
    const response = await this.api.post("/api/charge", data);
    return response.data;
  }

  // AI Assistant APIs
  async sendChatMessage(data: ChatRequest): Promise<ChatResponse> {
    const response = await this.api.post<ChatResponse>("/api/ai/chat", data);
    return response.data;
  }

  async submitFeedback(data: FeedbackRequest): Promise<{ success: boolean }> {
    const response = await this.api.post("/api/ai/feedback", data);
    return response.data;
  }

  async getAIStats(): Promise<AIStats> {
    const response = await this.api.get<AIStats>("/api/ai/stats");
    return response.data;
  }
}

export const apiService = new ApiService();

