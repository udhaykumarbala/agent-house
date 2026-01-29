import type { APIResponse, APIError, Summary, UpcomingResponse } from '@/types/api';
import type { Subscription, SubscriptionInput } from '@/types/subscription';
import type { User, LoginInput, RegisterInput, UserPreferences } from '@/types/user';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class APIClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseUrl}/api/v1${endpoint}`;

    const res = await fetch(url, {
      ...options,
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    const json = await res.json();

    if (!res.ok) {
      const error = json as APIError;
      throw new Error(error.error?.message || 'An error occurred');
    }

    return (json as APIResponse<T>).data;
  }

  // Auth
  async register(input: RegisterInput): Promise<User> {
    return this.request<User>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(input),
    });
  }

  async login(input: LoginInput): Promise<User> {
    return this.request<User>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(input),
    });
  }

  async logout(): Promise<void> {
    await this.request<{ message: string }>('/auth/logout', {
      method: 'POST',
    });
  }

  async me(): Promise<User> {
    return this.request<User>('/auth/me');
  }

  // Subscriptions
  async getSubscriptions(): Promise<Subscription[]> {
    return this.request<Subscription[]>('/subscriptions');
  }

  async getSubscription(id: string): Promise<Subscription> {
    return this.request<Subscription>(`/subscriptions/${id}`);
  }

  async createSubscription(input: SubscriptionInput): Promise<Subscription> {
    return this.request<Subscription>('/subscriptions', {
      method: 'POST',
      body: JSON.stringify(input),
    });
  }

  async updateSubscription(id: string, input: SubscriptionInput): Promise<Subscription> {
    return this.request<Subscription>(`/subscriptions/${id}`, {
      method: 'PUT',
      body: JSON.stringify(input),
    });
  }

  async deleteSubscription(id: string): Promise<void> {
    await this.request<{ message: string }>(`/subscriptions/${id}`, {
      method: 'DELETE',
    });
  }

  // Analytics
  async getSummary(): Promise<Summary> {
    return this.request<Summary>('/analytics/summary');
  }

  async getUpcoming(days: number = 30): Promise<UpcomingResponse> {
    return this.request<UpcomingResponse>(`/analytics/upcoming?days=${days}`);
  }

  // User
  async updatePreferences(prefs: UserPreferences): Promise<UserPreferences> {
    return this.request<UserPreferences>('/user/preferences', {
      method: 'PUT',
      body: JSON.stringify(prefs),
    });
  }

  async deleteAccount(): Promise<void> {
    await this.request<{ message: string }>('/user/account', {
      method: 'DELETE',
    });
  }
}

export const api = new APIClient(API_URL);
