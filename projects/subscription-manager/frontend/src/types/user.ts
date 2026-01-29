export interface User {
  id: string;
  email: string;
  emailVerified: boolean;
  preferences: UserPreferences;
  createdAt: string;
}

export interface UserPreferences {
  currency: string;
  theme: 'light' | 'dark' | 'system';
  reminderDays: number;
  emailNotifications: boolean;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput {
  email: string;
  password: string;
}
