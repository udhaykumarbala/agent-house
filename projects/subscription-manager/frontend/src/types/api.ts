export interface APIResponse<T> {
  data: T;
  meta?: {
    total?: number;
    page?: number;
    perPage?: number;
  };
}

export interface APIError {
  error: {
    code: string;
    message: string;
    details?: Array<{
      field: string;
      message: string;
    }>;
  };
}

export interface Summary {
  totalMonthly: number;
  totalYearly: number;
  activeCount: number;
  byCategory: Record<string, number>;
  dueSoon: number;
  dueSoonAmount: number;
}

export interface UpcomingGroup {
  label: string;
  subscriptions: import('./subscription').Subscription[];
  total: number;
}

export interface UpcomingResponse {
  groups: UpcomingGroup[];
  total: number;
  daysSpan: number;
}
