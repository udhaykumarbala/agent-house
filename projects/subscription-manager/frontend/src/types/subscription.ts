export type BillingCycle = 'weekly' | 'monthly' | 'yearly' | 'custom';

export type Category =
  | 'entertainment'
  | 'productivity'
  | 'lifestyle'
  | 'utilities'
  | 'news'
  | 'education'
  | 'finance'
  | 'health'
  | 'other';

export type Currency = 'USD' | 'EUR' | 'GBP' | 'CAD' | 'AUD' | 'INR';

export interface Subscription {
  id: string;
  name: string;
  description?: string;
  price: number;
  currency: Currency;
  billingCycle: BillingCycle;
  customDays?: number;
  nextBillingDate: string;
  category: Category;
  icon?: string;
  color?: string;
  notes?: string;
  isActive: boolean;
  reminderDays: number;
  createdAt: string;
  updatedAt: string;
}

export interface SubscriptionInput {
  name: string;
  description?: string;
  price: number;
  currency?: Currency;
  billingCycle: BillingCycle;
  customDays?: number;
  nextBillingDate: string;
  category?: Category;
  icon?: string;
  color?: string;
  notes?: string;
  reminderDays?: number;
}

export interface SubscriptionTemplate {
  name: string;
  icon: string;
  color: string;
  category: Category;
  defaultPrice: number;
  billingCycle: BillingCycle;
}

export const CATEGORIES: { value: Category; label: string; icon: string }[] = [
  { value: 'entertainment', label: 'Entertainment', icon: '🎬' },
  { value: 'productivity', label: 'Productivity', icon: '⚡' },
  { value: 'lifestyle', label: 'Lifestyle', icon: '🌟' },
  { value: 'utilities', label: 'Utilities', icon: '🔧' },
  { value: 'news', label: 'News', icon: '📰' },
  { value: 'education', label: 'Education', icon: '📚' },
  { value: 'finance', label: 'Finance', icon: '💰' },
  { value: 'health', label: 'Health', icon: '❤️' },
  { value: 'other', label: 'Other', icon: '📦' },
];

export const BILLING_CYCLES: { value: BillingCycle; label: string }[] = [
  { value: 'weekly', label: 'Weekly' },
  { value: 'monthly', label: 'Monthly' },
  { value: 'yearly', label: 'Yearly' },
  { value: 'custom', label: 'Custom' },
];

export function getMonthlyEquivalent(sub: Subscription): number {
  switch (sub.billingCycle) {
    case 'weekly':
      return (sub.price * 52) / 12;
    case 'monthly':
      return sub.price;
    case 'yearly':
      return sub.price / 12;
    case 'custom':
      if (sub.customDays && sub.customDays > 0) {
        return (sub.price * 30) / sub.customDays;
      }
      return sub.price;
    default:
      return sub.price;
  }
}

export function getYearlyEquivalent(sub: Subscription): number {
  switch (sub.billingCycle) {
    case 'weekly':
      return sub.price * 52;
    case 'monthly':
      return sub.price * 12;
    case 'yearly':
      return sub.price;
    case 'custom':
      if (sub.customDays && sub.customDays > 0) {
        return (sub.price * 365) / sub.customDays;
      }
      return sub.price * 12;
    default:
      return sub.price * 12;
  }
}

export function getDaysUntilBilling(nextBillingDate: string): number {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const billing = new Date(nextBillingDate);
  billing.setHours(0, 0, 0, 0);
  return Math.ceil((billing.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
}

export function isDueSoon(nextBillingDate: string, days: number = 7): boolean {
  const daysUntil = getDaysUntilBilling(nextBillingDate);
  return daysUntil >= 0 && daysUntil <= days;
}

export function formatBillingCycle(cycle: BillingCycle): string {
  switch (cycle) {
    case 'weekly':
      return '/wk';
    case 'monthly':
      return '/mo';
    case 'yearly':
      return '/yr';
    default:
      return '';
  }
}

export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
  });
}

export function formatFullDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  });
}
